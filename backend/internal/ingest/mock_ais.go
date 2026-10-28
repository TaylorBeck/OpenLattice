package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type mockVessel struct {
	mmsi     int64
	shipName string
	lat      float64
	lon      float64
	sog      float64
	heading  float64
}

// StartMockAISServer starts a local HTTP/WebSocket server emulating stream.aisstream.io.
// It returns the WS endpoint URL (e.g. "ws://127.0.0.1:54321/v0/stream") or an error.
func StartMockAISServer(ctx context.Context, listenAddr string) (string, error) {
	if listenAddr == "" {
		listenAddr = "127.0.0.1:0"
	}

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return "", fmt.Errorf("failed to bind mock AIS server: %w", err)
	}

	StartMockAISWithListener(ctx, listener)
	endpoint := fmt.Sprintf("ws://%s/v0/stream", listener.Addr().String())
	log.Printf("[MockAIS] Server listening at %s", endpoint)
	return endpoint, nil
}

// StartMockAISWithListener serves the mock AIS WebSocket feed over any net.Listener.
func StartMockAISWithListener(ctx context.Context, listener net.Listener) {
	upgrader := websocket.Upgrader{
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v0/stream", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("[MockAIS] Upgrade error: %v", err)
			return
		}
		handleMockAISClient(ctx, conn)
	})

	server := &http.Server{
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		server.Close()
		listener.Close()
	}()

	go func() {
		_ = server.Serve(listener)
	}()
}

func handleMockAISClient(ctx context.Context, conn *websocket.Conn) {
	defer conn.Close()
	log.Printf("[MockAIS] Client connected: %s", conn.RemoteAddr())

	// Handshake verification: read subscription within 3 seconds
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Printf("[MockAIS] Failed to receive handshake within deadline: %v", err)
		return
	}

	var sub AISStreamSubscription
	if err := json.Unmarshal(msg, &sub); err != nil {
		log.Printf("[MockAIS] Invalid subscription JSON: %v", err)
		return
	}
	_ = conn.SetReadDeadline(time.Time{})
	log.Printf("[MockAIS] Subscription validated from %s", conn.RemoteAddr())

	// Initialize synthetic vessel fleet
	fleet := []*mockVessel{
		{
			mmsi:     244670000,
			shipName: "ROTTERDAM CARRIER",
			lat:      51.9244,
			lon:      4.4777,
			sog:      14.5,
			heading:  245.0,
		},
		{
			mmsi:     367412340,
			shipName: "PACIFIC STAR",
			lat:      33.7405,
			lon:      -118.2710,
			sog:      9.8,
			heading:  130.0,
		},
		{
			mmsi:     211281610,
			shipName: "HAMBURG EXPRESS",
			lat:      53.5320,
			lon:      9.9800,
			sog:      12.0,
			heading:  75.0,
		},
		{
			mmsi:     563012900,
			shipName: "SINGAPORE PIONEER",
			lat:      1.2902,
			lon:      103.8519,
			sog:      16.2,
			heading:  110.0,
		},
	}

	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	var writeMu sync.Mutex

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now().UTC()
			timeStr := now.Format("2006-01-02 15:04:05.000000 +0000 UTC")

			for _, v := range fleet {
				advanceVessel(v, 0.3)

				envelope := AISStreamEnvelope{
					MessageType: "PositionReport",
					MetaData: AISStreamMetaData{
						MMSI:      v.mmsi,
						ShipName:  v.shipName,
						Latitude:  v.lat,
						Longitude: v.lon,
						TimeUTC:   timeStr,
					},
					Message: AISStreamInnerMessage{
						PositionReport: &AISPositionReport{
							Latitude:    v.lat,
							Longitude:   v.lon,
							Sog:         v.sog,
							TrueHeading: v.heading,
							UserID:      v.mmsi,
							Valid:       true,
						},
					},
				}

				payload, err := json.Marshal(envelope)
				if err != nil {
					continue
				}

				writeMu.Lock()
				// Send as binary WebSocket message per AISStream specification
				err = conn.WriteMessage(websocket.BinaryMessage, payload)
				writeMu.Unlock()

				if err != nil {
					log.Printf("[MockAIS] Client write error: %v", err)
					return
				}
			}
		}
	}
}

func advanceVessel(v *mockVessel, dtSeconds float64) {
	degPerSec := (v.sog / 3600.0) / 60.0 * dtSeconds
	rad := v.heading * math.Pi / 180.0

	v.lat += degPerSec * math.Cos(rad)
	v.lon += degPerSec * math.Sin(rad) / math.Cos(v.lat*math.Pi/180.0)

	if v.lat > 80.0 {
		v.lat = 80.0
		v.heading = 180.0
	} else if v.lat < -80.0 {
		v.lat = -80.0
		v.heading = 0.0
	}
}
