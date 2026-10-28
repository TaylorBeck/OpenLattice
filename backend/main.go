package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/TaylorBeck/openlattice/backend/internal/ingest"
	"github.com/TaylorBeck/openlattice/backend/internal/state"
	"github.com/TaylorBeck/openlattice/backend/pb"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all for local prototyping
	},
}

// ClientManager handles connected WebSocket clients in a thread-safe manner.
type ClientManager struct {
	clients map[*websocket.Conn]bool
	mutex   sync.RWMutex
}

// NewClientManager initializes an empty ClientManager.
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[*websocket.Conn]bool),
	}
}

// AddClient registers a new WebSocket client.
func (cm *ClientManager) AddClient(conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.clients[conn] = true
}

// RemoveClient unregisters and closes a WebSocket client connection.
func (cm *ClientManager) RemoveClient(conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	if _, ok := cm.clients[conn]; ok {
		delete(cm.clients, conn)
		_ = conn.Close()
	}
}

// ClientCount returns the number of active connected clients.
func (cm *ClientManager) ClientCount() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return len(cm.clients)
}

// Broadcast sends a binary frame to all currently connected clients.
func (cm *ClientManager) Broadcast(messageType int, payload []byte) {
	cm.mutex.RLock()
	conns := make([]*websocket.Conn, 0, len(cm.clients))
	for conn := range cm.clients {
		conns = append(conns, conn)
	}
	cm.mutex.RUnlock()

	for _, conn := range conns {
		_ = conn.SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
		if err := conn.WriteMessage(messageType, payload); err != nil {
			log.Printf("[WebSocket] Broadcast error to %s: %v", conn.RemoteAddr(), err)
			cm.RemoveClient(conn)
		}
	}
}

// convertTracksToPB converts internal state Tracks to Protobuf Track messages.
func convertTracksToPB(tracks []state.Track) []*pb.Track {
	pbTracks := make([]*pb.Track, 0, len(tracks))
	for _, t := range tracks {
		pbTracks = append(pbTracks, &pb.Track{
			EntityId:       t.EntityID,
			Class:          pb.Track_TrackClass(t.Class),
			Latitude:       t.Latitude,
			Longitude:      t.Longitude,
			AltitudeMeters: t.AltitudeMeters,
			VelocityKnots:  t.VelocityKnots,
			HeadingDegrees: t.HeadingDegrees,
			TimestampMs:    t.TimestampMs,
		})
	}
	return pbTracks
}

// startBroadcaster runs a 100ms ticker (10Hz) that serializes the current track store
// into a pb.TrackUpdateStream binary payload and broadcasts it over WebSockets.
func startBroadcaster(ctx context.Context, store *state.TrackStore, manager *ClientManager) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if manager.ClientCount() == 0 {
				continue
			}

			tracks := store.GetAllTracks()
			streamMsg := &pb.TrackUpdateStream{
				Tracks: convertTracksToPB(tracks),
			}

			payload, err := proto.Marshal(streamMsg)
			if err != nil {
				log.Printf("[Broadcaster] Failed to marshal TrackUpdateStream: %v", err)
				continue
			}

			manager.Broadcast(websocket.BinaryMessage, payload)
		}
	}
}

// streamHandler returns an http.HandlerFunc that upgrades requests to WebSockets and registers them.
func streamHandler(manager *ClientManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("[WebSocket] Upgrade error: %v", err)
			return
		}

		manager.AddClient(conn)
		log.Printf("[WebSocket] Client connected: %s (Active: %d)", conn.RemoteAddr(), manager.ClientCount())

		// Read pump to detect disconnection and process control messages
		go func(c *websocket.Conn) {
			defer func() {
				manager.RemoveClient(c)
				log.Printf("[WebSocket] Client disconnected: %s (Active: %d)", c.RemoteAddr(), manager.ClientCount())
			}()
			for {
				if _, _, err := c.NextReader(); err != nil {
					break
				}
			}
		}(conn)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := state.NewTrackStore(60 * time.Second)
	defer store.Stop()

	// 1. Aerial Stream (ADS-B) Ingestion
	adsbMock := os.Getenv("ADSB_MOCK") == "true"
	adsbEndpoint := os.Getenv("ADSB_ENDPOINT")

	if adsbMock {
		log.Println("[Main] ADSB_MOCK=true detected. Starting local synthetic ADS-B TCP generator...")
		mockAddr, err := ingest.StartMockADSBServer(ctx, "127.0.0.1:0")
		if err != nil {
			log.Fatalf("[Main] Failed to start mock ADS-B server: %v", err)
		}
		adsbEndpoint = mockAddr
	} else if adsbEndpoint == "" {
		adsbEndpoint = ingest.DefaultADSBEndpoint
	}

	adsbClient := ingest.NewADSBClient(ingest.ADSBConfig{
		Endpoint:       adsbEndpoint,
		DialTimeout:    10 * time.Second,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     30 * time.Second,
	}, store)

	go func() {
		if err := adsbClient.Start(ctx); err != nil && ctx.Err() == nil {
			log.Printf("[Main] ADSB ingestion terminated: %v", err)
		}
	}()

	// 2. Maritime Stream (AIS) Ingestion
	aisKey := os.Getenv("AISSTREAM_API_KEY")
	aisMock := os.Getenv("AIS_MOCK") == "true" || (adsbMock && aisKey == "")
	aisEndpoint := os.Getenv("AIS_ENDPOINT")

	if aisMock {
		log.Println("[Main] Starting local synthetic AIS WebSocket generator...")
		mockWSAddr, err := ingest.StartMockAISServer(ctx, "127.0.0.1:0")
		if err != nil {
			log.Fatalf("[Main] Failed to start mock AIS server: %v", err)
		}
		aisEndpoint = mockWSAddr
	} else if aisEndpoint == "" {
		aisEndpoint = ingest.DefaultAISStreamEndpoint
	}

	aisClient := ingest.NewAISStreamClient(ingest.AISConfig{
		Endpoint:       aisEndpoint,
		APIKey:         aisKey,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     30 * time.Second,
	}, store)

	go func() {
		if err := aisClient.Start(ctx); err != nil && ctx.Err() == nil {
			log.Printf("[Main] AIS ingestion terminated: %v", err)
		}
	}()

	// 3. Periodic normalized track count logger (every 5 seconds)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				counts := store.GetCountsByClass()
				aerial := counts[state.TrackClassCivilianAir]
				maritime := counts[state.TrackClassCivilianMaritime]
				total := store.Count()
				log.Printf("[Lattice C2] Active Tracks: Aerial=%d, Maritime=%d, Total=%d",
					aerial, maritime, total)
			}
		}
	}()

	manager := NewClientManager()

	// 4. Start 100ms WebSocket broadcaster (10Hz)
	go startBroadcaster(ctx, store, manager)

	mux := http.NewServeMux()
	mux.HandleFunc("/stream", streamHandler(manager))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("OpenLattice Ingestion Engine starting on :%s", port)
	log.Fatal(server.ListenAndServe())
}
