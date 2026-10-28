package ingest

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/TaylorBeck/openlattice/backend/internal/state"
)

const (
	DefaultADSBEndpoint = "data.adsbhub.org:5002"
	FeetToMeters        = 0.3048
)

// ADSBConfig defines the connection parameters for the ADSBHub feed.
type ADSBConfig struct {
	Endpoint       string
	DialTimeout    time.Duration
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

// ADSBClient connects to an SBS-1 BaseStation TCP stream and upserts tracks into TrackStore.
type ADSBClient struct {
	config  ADSBConfig
	store   *state.TrackStore
	NetDial func(ctx context.Context, network, addr string) (net.Conn, error)
}

// NewADSBClient creates a new ADSB ingestion client.
func NewADSBClient(cfg ADSBConfig, store *state.TrackStore) *ADSBClient {
	if cfg.Endpoint == "" {
		cfg.Endpoint = DefaultADSBEndpoint
	}
	if cfg.DialTimeout <= 0 {
		cfg.DialTimeout = 10 * time.Second
	}
	if cfg.InitialBackoff <= 0 {
		cfg.InitialBackoff = 1 * time.Second
	}
	if cfg.MaxBackoff <= 0 {
		cfg.MaxBackoff = 30 * time.Second
	}

	return &ADSBClient{
		config: cfg,
		store:  store,
	}
}

// Start runs the TCP ingestion loop with automatic exponential backoff reconnection.
func (c *ADSBClient) Start(ctx context.Context) error {
	backoff := c.config.InitialBackoff

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		log.Printf("[ADSB] Connecting to %s...", c.config.Endpoint)
		var conn net.Conn
		var err error

		if c.NetDial != nil {
			conn, err = c.NetDial(ctx, "tcp", c.config.Endpoint)
		} else {
			var dialer net.Dialer
			dialer.Timeout = c.config.DialTimeout
			conn, err = dialer.DialContext(ctx, "tcp", c.config.Endpoint)
		}

		if err != nil {
			log.Printf("[ADSB] Connection failed: %v. Retrying in %v...", err, backoff)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
			backoff *= 2
			if backoff > c.config.MaxBackoff {
				backoff = c.config.MaxBackoff
			}
			continue
		}

		// Reset backoff upon successful connection
		backoff = c.config.InitialBackoff
		log.Printf("[ADSB] Connected successfully to %s", c.config.Endpoint)

		// Handle stream until disconnected or ctx canceled
		if err := c.consumeStream(ctx, conn); err != nil && ctx.Err() == nil {
			log.Printf("[ADSB] Stream disconnected: %v", err)
		}
	}
}

func (c *ADSBClient) consumeStream(ctx context.Context, conn net.Conn) error {
	defer conn.Close()

	// Ensure connection closes when context is cancelled
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			conn.Close()
		case <-done:
		}
	}()

	scanner := bufio.NewScanner(conn)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 64*1024)

	for scanner.Scan() {
		line := scanner.Text()
		track, ok, err := ParseSBS1(line)
		if err != nil {
			// Log malformed lines at debug level and continue
			continue
		}
		if ok && track != nil {
			c.store.UpsertTrack(*track)
		}
	}

	return scanner.Err()
}

// ParseSBS1 parses a single BaseStation SBS-1 CSV message line.
// Supports MSG,3 (airborne position) and MSG,4 (airborne velocity).
func ParseSBS1(line string) (*state.Track, bool, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, false, nil
	}

	parts := strings.Split(line, ",")
	if len(parts) < 5 {
		return nil, false, fmt.Errorf("malformed SBS-1 message: expected at least 5 fields, got %d", len(parts))
	}

	if parts[0] != "MSG" {
		return nil, false, nil // Not a message frame, safely ignore
	}

	msgType := strings.TrimSpace(parts[1])
	entityID := strings.TrimSpace(parts[4])
	if entityID == "" {
		return nil, false, fmt.Errorf("malformed SBS-1 message: empty entity ICAO hex ID")
	}

	timestampMs := parseSBS1Timestamp(parts)

	switch msgType {
	case "3": // Airborne Position Message
		if len(parts) < 16 {
			return nil, false, fmt.Errorf("malformed MSG,3: expected at least 16 fields, got %d", len(parts))
		}

		altFeetStr := strings.TrimSpace(parts[11])
		latStr := strings.TrimSpace(parts[14])
		lonStr := strings.TrimSpace(parts[15])

		if altFeetStr == "" || latStr == "" || lonStr == "" {
			return nil, false, fmt.Errorf("malformed MSG,3: missing coordinates or altitude")
		}

		altFeet, err := strconv.ParseFloat(altFeetStr, 64)
		if err != nil {
			return nil, false, fmt.Errorf("malformed MSG,3: invalid altitude %q: %w", altFeetStr, err)
		}

		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			return nil, false, fmt.Errorf("malformed MSG,3: invalid latitude %q: %w", latStr, err)
		}

		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			return nil, false, fmt.Errorf("malformed MSG,3: invalid longitude %q: %w", lonStr, err)
		}

		return &state.Track{
			EntityID:       entityID,
			Class:          state.TrackClassCivilianAir,
			Latitude:       lat,
			Longitude:      lon,
			AltitudeMeters: altFeet * FeetToMeters,
			TimestampMs:    timestampMs,
			HasPosition:    true,
		}, true, nil

	case "4": // Airborne Velocity Message
		if len(parts) < 14 {
			return nil, false, fmt.Errorf("malformed MSG,4: expected at least 14 fields, got %d", len(parts))
		}

		spdStr := strings.TrimSpace(parts[12])
		hdgStr := strings.TrimSpace(parts[13])

		if spdStr == "" || hdgStr == "" {
			return nil, false, fmt.Errorf("malformed MSG,4: missing speed or heading")
		}

		speed, err := strconv.ParseFloat(spdStr, 64)
		if err != nil {
			return nil, false, fmt.Errorf("malformed MSG,4: invalid ground speed %q: %w", spdStr, err)
		}

		hdg, err := strconv.ParseFloat(hdgStr, 64)
		if err != nil {
			return nil, false, fmt.Errorf("malformed MSG,4: invalid heading %q: %w", hdgStr, err)
		}

		return &state.Track{
			EntityID:       entityID,
			Class:          state.TrackClassCivilianAir,
			VelocityKnots:  speed,
			HeadingDegrees: hdg,
			TimestampMs:    timestampMs,
			HasVelocity:    true,
		}, true, nil

	default:
		// Other SBS-1 transmission types (e.g. MSG,1 callsign, MSG,5 surveillance alt)
		return nil, false, nil
	}
}

func parseSBS1Timestamp(parts []string) int64 {
	// Date generated is index 6, Time generated is index 7
	if len(parts) > 7 {
		dateStr := strings.TrimSpace(parts[6])
		timeStr := strings.TrimSpace(parts[7])
		if dateStr != "" && timeStr != "" {
			dateTimeStr := dateStr + " " + timeStr
			// Try with milliseconds first
			if t, err := time.Parse("2006/01/02 15:04:05.000", dateTimeStr); err == nil {
				return t.UnixMilli()
			}
			// Try standard seconds layout
			if t, err := time.Parse("2006/01/02 15:04:05", dateTimeStr); err == nil {
				return t.UnixMilli()
			}
		}
	}
	return time.Now().UnixMilli()
}
