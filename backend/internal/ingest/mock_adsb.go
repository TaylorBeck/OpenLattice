package ingest

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"sync"
	"time"
)

type mockAircraft struct {
	icao      string
	lat       float64
	lon       float64
	altFeet   float64
	speed     float64
	heading   float64
	climbRate float64
}

// StartMockADSBServer starts a local TCP server that streams synthetic SBS-1 messages.
// It returns the bound address (e.g., "127.0.0.1:54321") or an error.
func StartMockADSBServer(ctx context.Context, listenAddr string) (string, error) {
	if listenAddr == "" {
		listenAddr = "127.0.0.1:0"
	}

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return "", fmt.Errorf("failed to bind mock ADSB server: %w", err)
	}

	actualAddr := listener.Addr().String()
	log.Printf("[MockADSB] Server listening on %s", actualAddr)

	StartMockADSBWithListener(ctx, listener)
	return actualAddr, nil
}

// StartMockADSBWithListener serves synthetic SBS-1 streams over any provided net.Listener.
func StartMockADSBWithListener(ctx context.Context, listener net.Listener) {
	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	go func() {
		var wg sync.WaitGroup
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					wg.Wait()
					return
				default:
					return
				}
			}

			wg.Add(1)
			go func(c net.Conn) {
				defer wg.Done()
				handleMockClient(ctx, c)
			}(conn)
		}
	}()
}

func handleMockClient(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	log.Printf("[MockADSB] Client connected from %s", conn.RemoteAddr())

	// Initialize synthetic aircraft fleet
	fleet := []*mockAircraft{
		{
			icao:      "4CA2D6",
			lat:       51.4700,
			lon:       -0.4543,
			altFeet:   36000,
			speed:     480.0,
			heading:   270.0,
			climbRate: 0.0,
		},
		{
			icao:      "A1B2C3",
			lat:       40.6413,
			lon:       -73.7781,
			altFeet:   14500,
			speed:     320.0,
			heading:   65.0,
			climbRate: 50.0,
		},
		{
			icao:      "3C6442",
			lat:       37.6213,
			lon:       -122.3790,
			altFeet:   4500,
			speed:     140.0,
			heading:   190.0,
			climbRate: -20.0,
		},
		{
			icao:      "7102EF",
			lat:       35.5494,
			lon:       139.7798,
			altFeet:   28000,
			speed:     460.0,
			heading:   310.0,
			climbRate: 10.0,
		},
	}

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	// Interleave MSG,3 and MSG,4 on alternate ticks
	isPositionTick := true

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now().UTC()
			dateStr := now.Format("2006/01/02")
			timeStr := now.Format("15:04:05.000")

			for _, ac := range fleet {
				// Kinematic update
				advanceAircraft(ac, 0.25)

				var msg string
				if isPositionTick {
					// MSG,3 airborne position
					msg = fmt.Sprintf("MSG,3,1,1,%s,1,%s,%s,%s,%s,,%.0f,,,%.4f,%.4f,,,,,,0\r\n",
						ac.icao, dateStr, timeStr, dateStr, timeStr, ac.altFeet, ac.lat, ac.lon)
				} else {
					// MSG,4 airborne velocity
					msg = fmt.Sprintf("MSG,4,1,1,%s,1,%s,%s,%s,%s,,,%.1f,%.1f,,,,,,,0\r\n",
						ac.icao, dateStr, timeStr, dateStr, timeStr, ac.speed, ac.heading)
				}

				if _, err := conn.Write([]byte(msg)); err != nil {
					log.Printf("[MockADSB] Client %s write error: %v", conn.RemoteAddr(), err)
					return
				}
			}

			isPositionTick = !isPositionTick
		}
	}
}

func advanceAircraft(ac *mockAircraft, dtSeconds float64) {
	// Speed in knots to degrees/sec approximation: 1 knot = 1 nautical mile/hr = 1/60 deg/hr
	degPerSec := (ac.speed / 3600.0) / 60.0 * dtSeconds
	rad := ac.heading * math.Pi / 180.0

	ac.lat += degPerSec * math.Cos(rad)
	ac.lon += degPerSec * math.Sin(rad) / math.Cos(ac.lat*math.Pi/180.0)

	// Keep coordinates within bounds
	if ac.lat > 85.0 {
		ac.lat = 85.0
		ac.heading = 180.0
	} else if ac.lat < -85.0 {
		ac.lat = -85.0
		ac.heading = 0.0
	}

	ac.altFeet += ac.climbRate * dtSeconds
	if ac.altFeet > 42000 {
		ac.climbRate = -ac.climbRate
	} else if ac.altFeet < 2000 {
		ac.climbRate = math.Abs(ac.climbRate)
	}
}
