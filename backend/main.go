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
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all for local prototyping
	},
}

// ClientManager handles connected WebSocket clients
type ClientManager struct {
	clients map[*websocket.Conn]bool
	mutex   sync.RWMutex
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

	manager := &ClientManager{
		clients: make(map[*websocket.Conn]bool),
	}

	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket Upgrade Error:", err)
			return
		}

		manager.mutex.Lock()
		manager.clients[conn] = true
		manager.mutex.Unlock()

		log.Println("New tactical C2 client connected")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("OpenLattice Ingestion Engine starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
