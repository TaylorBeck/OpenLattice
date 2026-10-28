package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/TaylorBeck/openlattice/backend/internal/state"
	"github.com/gorilla/websocket"
)

const (
	DefaultAISStreamEndpoint = "wss://stream.aisstream.io/v0/stream"
)

// AISConfig defines the parameters for AISStream.io ingestion.
type AISConfig struct {
	Endpoint       string
	APIKey         string
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

// AISStreamSubscription matches the required AISStream.io handshake contract.
type AISStreamSubscription struct {
	APIKey             string           `json:"APIKey"`
	BoundingBoxes      [][][2]float64   `json:"BoundingBoxes"`
	FilterMessageTypes []string         `json:"FilterMessageTypes"`
}

// AISStreamMetaData contains track metadata.
type AISStreamMetaData struct {
	MMSI      int64   `json:"MMSI"`
	ShipName  string  `json:"ShipName"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	TimeUTC   string  `json:"time_utc"`
}

// AISPositionReport contains kinematic telemetry for a vessel.
type AISPositionReport struct {
	Cog         float64 `json:"Cog"`
	Latitude    float64 `json:"Latitude"`
	Longitude   float64 `json:"Longitude"`
	Sog         float64 `json:"Sog"`
	TrueHeading float64 `json:"TrueHeading"`
	UserID      int64   `json:"UserID"`
	Valid       bool    `json:"Valid"`
}

// AISStreamInnerMessage holds the polymorphic AIS message body.
type AISStreamInnerMessage struct {
	PositionReport *AISPositionReport `json:"PositionReport,omitempty"`
}

// AISStreamEnvelope is the root container for AISStream WebSocket frames.
type AISStreamEnvelope struct {
	MessageType string                `json:"MessageType"`
	MetaData    AISStreamMetaData     `json:"MetaData"`
	Message     AISStreamInnerMessage `json:"Message"`
}

// AISStreamClient manages WebSocket connectivity to AISStream.io.
type AISStreamClient struct {
	config  AISConfig
	store   *state.TrackStore
	NetDial func(ctx context.Context, network, addr string) (net.Conn, error)
}

// NewAISStreamClient constructs a new AIS WebSocket ingestion client.
func NewAISStreamClient(cfg AISConfig, store *state.TrackStore) *AISStreamClient {
	if cfg.Endpoint == "" {
		cfg.Endpoint = DefaultAISStreamEndpoint
	}
	if cfg.InitialBackoff <= 0 {
		cfg.InitialBackoff = 1 * time.Second
	}
	if cfg.MaxBackoff <= 0 {
		cfg.MaxBackoff = 30 * time.Second
	}

	return &AISStreamClient{
		config: cfg,
		store:  store,
	}
}

// Start manages the connection lifecycle with permessage-deflate compression and handshake.
func (c *AISStreamClient) Start(ctx context.Context) error {
	backoff := c.config.InitialBackoff

	dialer := websocket.Dialer{
		EnableCompression: true,
		HandshakeTimeout:  10 * time.Second,
	}
	if c.NetDial != nil {
		dialer.NetDialContext = c.NetDial
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		log.Printf("[AIS] Connecting to %s (EnableCompression=true)...", c.config.Endpoint)
		conn, _, err := dialer.DialContext(ctx, c.config.Endpoint, http.Header{})
		if err != nil {
			log.Printf("[AIS] Connection error: %v. Retrying in %v...", err, backoff)
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

		// Reset backoff on successful connect
		backoff = c.config.InitialBackoff
		log.Printf("[AIS] Connected successfully to %s", c.config.Endpoint)

		if err := c.handleSession(ctx, conn); err != nil && ctx.Err() == nil {
			log.Printf("[AIS] Stream closed: %v", err)
		}
	}
}

func (c *AISStreamClient) handleSession(ctx context.Context, conn *websocket.Conn) error {
	defer conn.Close()

	// Transmit handshake subscription within 3 seconds
	handshakeDeadline := time.Now().Add(3 * time.Second)
	subMsg := AISStreamSubscription{
		APIKey: c.config.APIKey,
		BoundingBoxes: [][][2]float64{
			{{-90, -180}, {90, 180}},
		},
		FilterMessageTypes: []string{"PositionReport"},
	}

	if err := conn.SetWriteDeadline(handshakeDeadline); err != nil {
		return fmt.Errorf("failed to set handshake write deadline: %w", err)
	}

	if err := conn.WriteJSON(subMsg); err != nil {
		return fmt.Errorf("failed to send AISStream subscription handshake: %w", err)
	}

	// Reset write deadline
	_ = conn.SetWriteDeadline(time.Time{})
	log.Printf("[AIS] Subscription handshake successfully transmitted")

	// Ensure socket closes on context cancellation
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			conn.Close()
		case <-done:
		}
	}()

	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		// Accept both binary frames (as specified by AISStream) and text frames
		if messageType != websocket.BinaryMessage && messageType != websocket.TextMessage {
			continue
		}

		track, ok, err := ParseAISMessage(payload)
		if err != nil {
			// Skip unrecognized frames or malformed payloads
			continue
		}
		if ok && track != nil {
			c.store.UpsertTrack(*track)
		}
	}
}

// ParseAISMessage unmarshals a UTF-8 JSON AIS payload into a normalized track.
func ParseAISMessage(payload []byte) (*state.Track, bool, error) {
	if len(payload) == 0 {
		return nil, false, nil
	}

	var env AISStreamEnvelope
	if err := json.Unmarshal(payload, &env); err != nil {
		return nil, false, fmt.Errorf("malformed AIS JSON: %w", err)
	}

	var mmsi int64
	if env.MetaData.MMSI != 0 {
		mmsi = env.MetaData.MMSI
	} else if env.Message.PositionReport != nil && env.Message.PositionReport.UserID != 0 {
		mmsi = env.Message.PositionReport.UserID
	}

	if mmsi == 0 {
		return nil, false, fmt.Errorf("missing MMSI in AIS message")
	}

	entityID := strconv.FormatInt(mmsi, 10)

	lat := env.MetaData.Latitude
	lon := env.MetaData.Longitude
	var sog float64
	var heading float64

	if env.Message.PositionReport != nil {
		if lat == 0 && lon == 0 {
			lat = env.Message.PositionReport.Latitude
			lon = env.Message.PositionReport.Longitude
		}
		sog = env.Message.PositionReport.Sog
		heading = env.Message.PositionReport.TrueHeading
	}

	timestampMs := parseAISTimestamp(env.MetaData.TimeUTC)

	track := &state.Track{
		EntityID:       entityID,
		Class:          state.TrackClassCivilianMaritime,
		Latitude:       lat,
		Longitude:      lon,
		AltitudeMeters: 0.0, // Surface vessels are sea level
		VelocityKnots:  sog,
		HeadingDegrees: heading,
		TimestampMs:    timestampMs,
		HasPosition:    true,
		HasVelocity:    true,
	}

	return track, true, nil
}

func parseAISTimestamp(timeStr string) int64 {
	timeStr = strings.TrimSpace(timeStr)
	if timeStr != "" {
		// AISStream typically emits e.g. "2026-10-27 19:20:05.123456 +0000 UTC"
		layouts := []string{
			"2006-01-02 15:04:05.999999999 -0700 MST",
			"2006-01-02 15:04:05 -0700 MST",
			time.RFC3339Nano,
			time.RFC3339,
		}
		for _, layout := range layouts {
			if t, err := time.Parse(layout, timeStr); err == nil {
				return t.UnixMilli()
			}
		}
	}
	return time.Now().UnixMilli()
}
