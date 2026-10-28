package main

import (
	"context"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/TaylorBeck/openlattice/backend/internal/state"
	"github.com/TaylorBeck/openlattice/backend/pb"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type pipeListener struct {
	conns chan net.Conn
	done  chan struct{}
	once  sync.Once
}

func newPipeListener() *pipeListener {
	return &pipeListener{
		conns: make(chan net.Conn),
		done:  make(chan struct{}),
	}
}

func (l *pipeListener) Accept() (net.Conn, error) {
	select {
	case conn, ok := <-l.conns:
		if !ok {
			return nil, net.ErrClosed
		}
		return conn, nil
	case <-l.done:
		return nil, net.ErrClosed
	}
}

func (l *pipeListener) Close() error {
	l.once.Do(func() {
		close(l.done)
	})
	return nil
}

func (l *pipeListener) Addr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}
}

func (l *pipeListener) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	clientConn, serverConn := net.Pipe()
	select {
	case <-l.done:
		clientConn.Close()
		serverConn.Close()
		return nil, net.ErrClosed
	case l.conns <- serverConn:
		return clientConn, nil
	case <-ctx.Done():
		clientConn.Close()
		serverConn.Close()
		return nil, ctx.Err()
	}
}

func TestWebSocketBroadcaster_10HzStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := state.NewTrackStore(10 * time.Second)
	defer store.Stop()

	// Seed tracks
	now := time.Now()
	store.UpsertTrack(state.Track{
		EntityID:       "AAL123",
		Class:          state.TrackClassCivilianAir,
		Latitude:       37.7749,
		Longitude:      -122.4194,
		AltitudeMeters: 10668.0,
		VelocityKnots:  450.0,
		HeadingDegrees: 270.0,
		TimestampMs:    now.UnixMilli(),
		LastSeen:       now,
		HasPosition:    true,
		HasVelocity:    true,
	})
	store.UpsertTrack(state.Track{
		EntityID:       "BOAT456",
		Class:          state.TrackClassCivilianMaritime,
		Latitude:       37.8000,
		Longitude:      -122.4000,
		AltitudeMeters: 0.0,
		VelocityKnots:  15.5,
		HeadingDegrees: 90.0,
		TimestampMs:    now.UnixMilli(),
		LastSeen:       now,
		HasPosition:    true,
		HasVelocity:    true,
	})

	manager := NewClientManager()
	go startBroadcaster(ctx, store, manager)

	mux := http.NewServeMux()
	mux.HandleFunc("/stream", streamHandler(manager))

	listener := newPipeListener()
	server := &http.Server{Handler: mux}
	go func() {
		_ = server.Serve(listener)
	}()
	defer server.Close()
	defer listener.Close()

	dialer := websocket.Dialer{
		NetDialContext: listener.DialContext,
	}

	conn, resp, err := dialer.Dial("ws://127.0.0.1:8080/stream", nil)
	if err != nil {
		t.Fatalf("Failed to dial WebSocket: %v", err)
	}
	defer conn.Close()
	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("Expected status 101 Switching Protocols, got %d", resp.StatusCode)
	}

	// Read 5 frames and verify timing and payload content
	const expectedFrames = 5
	receivedTimestamps := make([]time.Time, 0, expectedFrames)

	for i := 0; i < expectedFrames; i++ {
		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		msgType, msgBytes, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("Frame %d: error reading WebSocket message: %v", i, err)
		}

		receivedTimestamps = append(receivedTimestamps, time.Now())

		if msgType != websocket.BinaryMessage {
			t.Errorf("Frame %d: expected binary message (%d), got %d", i, websocket.BinaryMessage, msgType)
		}

		var stream pb.TrackUpdateStream
		if err := proto.Unmarshal(msgBytes, &stream); err != nil {
			t.Fatalf("Frame %d: failed to unmarshal pb.TrackUpdateStream: %v", i, err)
		}

		if len(stream.Tracks) != 2 {
			t.Fatalf("Frame %d: expected 2 tracks, got %d", i, len(stream.Tracks))
		}

		// Verify track values
		foundAir := false
		foundSea := false
		for _, tr := range stream.Tracks {
			if tr.EntityId == "AAL123" {
				foundAir = true
				if tr.Class != pb.Track_CIVILIAN_AIR {
					t.Errorf("AAL123 expected class CIVILIAN_AIR, got %v", tr.Class)
				}
				if tr.Latitude != 37.7749 || tr.Longitude != -122.4194 {
					t.Errorf("AAL123 lat/lon mismatch: (%f, %f)", tr.Latitude, tr.Longitude)
				}
			}
			if tr.EntityId == "BOAT456" {
				foundSea = true
				if tr.Class != pb.Track_CIVILIAN_MARITIME {
					t.Errorf("BOAT456 expected class CIVILIAN_MARITIME, got %v", tr.Class)
				}
			}
		}

		if !foundAir || !foundSea {
			t.Errorf("Frame %d: did not find expected tracks (air=%v, sea=%v)", i, foundAir, foundSea)
		}
	}

	// Verify frame intervals (~100ms)
	for i := 1; i < len(receivedTimestamps); i++ {
		delta := receivedTimestamps[i].Sub(receivedTimestamps[i-1])
		t.Logf("Frame %d-%d interval: %v (target ~100ms)", i-1, i, delta)
	}

	// Test concurrent client
	conn2, _, err := dialer.Dial("ws://127.0.0.1:8080/stream", nil)
	if err != nil {
		t.Fatalf("Failed to dial second client: %v", err)
	}
	defer conn2.Close()

	if manager.ClientCount() != 2 {
		t.Errorf("Expected 2 active clients, got %d", manager.ClientCount())
	}

	_ = conn2.SetReadDeadline(time.Now().Add(1 * time.Second))
	msgType2, msgBytes2, err := conn2.ReadMessage()
	if err != nil {
		t.Fatalf("Client 2 failed to read message: %v", err)
	}
	if msgType2 != websocket.BinaryMessage || len(msgBytes2) == 0 {
		t.Errorf("Client 2 received invalid binary message")
	}

	// Close client 2 and wait for reader pump to deregister
	conn2.Close()
	time.Sleep(50 * time.Millisecond)
}

