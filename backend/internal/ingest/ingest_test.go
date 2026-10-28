package ingest_test

import (
	"context"
	"math"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/TaylorBeck/openlattice/backend/internal/ingest"
	"github.com/TaylorBeck/openlattice/backend/internal/state"
)

const floatTolerance = 1e-4

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < floatTolerance
}

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
	case <-l.done:
		return nil, net.ErrClosed
	case c, ok := <-l.conns:
		if !ok {
			return nil, net.ErrClosed
		}
		return c, nil
	}
}

func (l *pipeListener) Close() error {
	l.once.Do(func() {
		close(l.done)
	})
	return nil
}

func (l *pipeListener) Addr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 80}
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

func TestParseSBS1_PositionMessage(t *testing.T) {
	line := "MSG,3,1,1,4CA2D6,1,2026/10/27,19:20:05.123,2026/10/27,19:20:05.123,,36000,,,51.4700,-0.4543,,,,,,0"
	track, ok, err := ingest.ParseSBS1(line)
	if err != nil {
		t.Fatalf("unexpected error parsing MSG,3: %v", err)
	}
	if !ok || track == nil {
		t.Fatalf("expected track to be parsed successfully")
	}

	if track.EntityID != "4CA2D6" {
		t.Errorf("expected EntityID 4CA2D6, got %s", track.EntityID)
	}
	if track.Class != state.TrackClassCivilianAir {
		t.Errorf("expected Class CivilianAir, got %v", track.Class)
	}
	expectedAltitudeMeters := 36000.0 * 0.3048
	if !almostEqual(track.AltitudeMeters, expectedAltitudeMeters) {
		t.Errorf("expected AltitudeMeters %f, got %f", expectedAltitudeMeters, track.AltitudeMeters)
	}
	if !almostEqual(track.Latitude, 51.4700) {
		t.Errorf("expected Latitude 51.4700, got %f", track.Latitude)
	}
	if !almostEqual(track.Longitude, -0.4543) {
		t.Errorf("expected Longitude -0.4543, got %f", track.Longitude)
	}
	if !track.HasPosition {
		t.Errorf("expected HasPosition to be true")
	}
	if track.HasVelocity {
		t.Errorf("expected HasVelocity to be false for MSG,3")
	}
}

func TestParseSBS1_VelocityMessage(t *testing.T) {
	line := "MSG,4,1,1,4CA2D6,1,2026/10/27,19:20:05.123,2026/10/27,19:20:05.123,,,450.0,180.0,,,,,,,0"
	track, ok, err := ingest.ParseSBS1(line)
	if err != nil {
		t.Fatalf("unexpected error parsing MSG,4: %v", err)
	}
	if !ok || track == nil {
		t.Fatalf("expected track to be parsed successfully")
	}

	if track.EntityID != "4CA2D6" {
		t.Errorf("expected EntityID 4CA2D6, got %s", track.EntityID)
	}
	if track.Class != state.TrackClassCivilianAir {
		t.Errorf("expected Class CivilianAir, got %v", track.Class)
	}
	if !almostEqual(track.VelocityKnots, 450.0) {
		t.Errorf("expected VelocityKnots 450.0, got %f", track.VelocityKnots)
	}
	if !almostEqual(track.HeadingDegrees, 180.0) {
		t.Errorf("expected HeadingDegrees 180.0, got %f", track.HeadingDegrees)
	}
	if !track.HasVelocity {
		t.Errorf("expected HasVelocity to be true")
	}
	if track.HasPosition {
		t.Errorf("expected HasPosition to be false for MSG,4")
	}
}

func TestParseSBS1_MalformedPackets(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectError bool
		expectOk    bool
	}{
		{
			name:        "empty string",
			input:       "",
			expectError: false,
			expectOk:    false,
		},
		{
			name:        "whitespace only",
			input:       "   \t\n  ",
			expectError: false,
			expectOk:    false,
		},
		{
			name:        "non-MSG transmission type",
			input:       "CLK,1,1,1,4CA2D6",
			expectError: false,
			expectOk:    false,
		},
		{
			name:        "too few fields total",
			input:       "MSG,3,1",
			expectError: true,
			expectOk:    false,
		},
		{
			name:        "empty ICAO entity id",
			input:       "MSG,3,1,1,,1,2026/10/27,19:20:05.123,2026/10/27,19:20:05.123,,36000,,,51.4700,-0.4543,,,,,,0",
			expectError: true,
			expectOk:    false,
		},
		{
			name:        "MSG,3 too few fields",
			input:       "MSG,3,1,1,4CA2D6,1,2026/10/27,19:20:05.123,,36000",
			expectError: true,
			expectOk:    false,
		},
		{
			name:        "MSG,3 invalid non-numeric altitude",
			input:       "MSG,3,1,1,4CA2D6,1,2026/10/27,19:20:05.123,2026/10/27,19:20:05.123,,NOT_A_NUMBER,,,51.4700,-0.4543,,,,,,0",
			expectError: true,
			expectOk:    false,
		},
		{
			name:        "MSG,3 invalid non-numeric latitude",
			input:       "MSG,3,1,1,4CA2D6,1,2026/10/27,19:20:05.123,2026/10/27,19:20:05.123,,36000,,,BAD_LAT,-0.4543,,,,,,0",
			expectError: true,
			expectOk:    false,
		},
		{
			name:        "MSG,3 invalid non-numeric longitude",
			input:       "MSG,3,1,1,4CA2D6,1,2026/10/27,19:20:05.123,2026/10/27,19:20:05.123,,36000,,,51.4700,BAD_LON,,,,,,0",
			expectError: true,
			expectOk:    false,
		},
		{
			name:        "MSG,4 too few fields",
			input:       "MSG,4,1,1,4CA2D6,1,2026/10/27,19:20:05.123,,,450.0",
			expectError: true,
			expectOk:    false,
		},
		{
			name:        "MSG,4 invalid ground speed",
			input:       "MSG,4,1,1,4CA2D6,1,2026/10/27,19:20:05.123,2026/10/27,19:20:05.123,,,FAST,180.0,,,,,,,0",
			expectError: true,
			expectOk:    false,
		},
		{
			name:        "MSG,4 invalid heading",
			input:       "MSG,4,1,1,4CA2D6,1,2026/10/27,19:20:05.123,2026/10/27,19:20:05.123,,,450.0,NORTH,,,,,,,0",
			expectError: true,
			expectOk:    false,
		},
		{
			name:        "unhandled SBS-1 message type MSG,1 (callsign) ignored cleanly",
			input:       "MSG,1,1,1,4CA2D6,1,2026/10/27,19:20:05.123,2026/10/27,19:20:05.123,BAW123,,,,,,,,,,,0",
			expectError: false,
			expectOk:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			track, ok, err := ingest.ParseSBS1(tc.input)
			if tc.expectError && err == nil {
				t.Fatalf("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ok != tc.expectOk {
				t.Fatalf("expected ok=%v, got %v", tc.expectOk, ok)
			}
			if !tc.expectOk && track != nil {
				t.Fatalf("expected nil track when ok is false")
			}
		})
	}
}

func TestPartialStateMerging_EndToEnd(t *testing.T) {
	store := state.NewTrackStore(60 * time.Second)
	defer store.Stop()

	// Scenario 1: MSG,3 arrives, followed by MSG,4
	msg3 := "MSG,3,1,1,AABB01,1,2026/10/27,19:20:05.123,2026/10/27,19:20:05.123,,30000,,,42.3601,-71.0589,,,,,,0"
	t3, ok, err := ingest.ParseSBS1(msg3)
	if err != nil || !ok {
		t.Fatalf("failed to parse MSG,3: %v", err)
	}
	store.UpsertTrack(*t3)

	after3, found := store.GetTrack("AABB01")
	if !found {
		t.Fatalf("track AABB01 not found after MSG,3")
	}
	if !almostEqual(after3.Latitude, 42.3601) || !almostEqual(after3.Longitude, -71.0589) {
		t.Errorf("unexpected position after MSG,3: %+v", after3)
	}
	if after3.VelocityKnots != 0 || after3.HeadingDegrees != 0 {
		t.Errorf("velocity/heading should be zero before MSG,4: %+v", after3)
	}

	msg4 := "MSG,4,1,1,AABB01,1,2026/10/27,19:20:06.123,2026/10/27,19:20:06.123,,,480.0,210.0,,,,,,,0"
	t4, ok, err := ingest.ParseSBS1(msg4)
	if err != nil || !ok {
		t.Fatalf("failed to parse MSG,4: %v", err)
	}
	store.UpsertTrack(*t4)

	merged, found := store.GetTrack("AABB01")
	if !found {
		t.Fatalf("track AABB01 missing after MSG,4")
	}
	// Check all fields are merged intact
	if !almostEqual(merged.Latitude, 42.3601) || !almostEqual(merged.Longitude, -71.0589) {
		t.Errorf("position lost after MSG,4: %+v", merged)
	}
	if !almostEqual(merged.AltitudeMeters, 30000.0*0.3048) {
		t.Errorf("altitude lost after MSG,4: %+v", merged)
	}
	if !almostEqual(merged.VelocityKnots, 480.0) || !almostEqual(merged.HeadingDegrees, 210.0) {
		t.Errorf("velocity/heading not merged: %+v", merged)
	}

	// Scenario 2: MSG,4 arrives first, then MSG,3
	msg4First := "MSG,4,1,1,CCDD02,1,2026/10/27,19:20:05.123,2026/10/27,19:20:05.123,,,350.0,90.0,,,,,,,0"
	t4First, ok, err := ingest.ParseSBS1(msg4First)
	if err != nil || !ok {
		t.Fatalf("failed to parse MSG,4 first: %v", err)
	}
	store.UpsertTrack(*t4First)

	msg3Second := "MSG,3,1,1,CCDD02,1,2026/10/27,19:20:06.123,2026/10/27,19:20:06.123,,15000,,,34.0522,-118.2437,,,,,,0"
	t3Second, ok, err := ingest.ParseSBS1(msg3Second)
	if err != nil || !ok {
		t.Fatalf("failed to parse MSG,3 second: %v", err)
	}
	store.UpsertTrack(*t3Second)

	mergedReverse, found := store.GetTrack("CCDD02")
	if !found {
		t.Fatalf("track CCDD02 missing after reverse merge")
	}
	if !almostEqual(mergedReverse.VelocityKnots, 350.0) || !almostEqual(mergedReverse.HeadingDegrees, 90.0) {
		t.Errorf("velocity corrupted in reverse merge: %+v", mergedReverse)
	}
	if !almostEqual(mergedReverse.Latitude, 34.0522) || !almostEqual(mergedReverse.Longitude, -118.2437) {
		t.Errorf("position not set in reverse merge: %+v", mergedReverse)
	}
	if !almostEqual(mergedReverse.AltitudeMeters, 15000.0*0.3048) {
		t.Errorf("altitude not set in reverse merge: %+v", mergedReverse)
	}
}

func TestParseAISMessage_ValidJSON(t *testing.T) {
	rawJSON := `{
		"MessageType": "PositionReport",
		"MetaData": {
			"MMSI": 244670000,
			"ShipName": "ROTTERDAM",
			"latitude": 51.9244,
			"longitude": 4.4777,
			"time_utc": "2026-10-27 19:20:05.123456 +0000 UTC"
		},
		"Message": {
			"PositionReport": {
				"Cog": 123.4,
				"Latitude": 51.9244,
				"Longitude": 4.4777,
				"Sog": 14.5,
				"TrueHeading": 245.0,
				"UserID": 244670000,
				"Valid": true
			}
		}
	}`

	track, ok, err := ingest.ParseAISMessage([]byte(rawJSON))
	if err != nil {
		t.Fatalf("unexpected error parsing AIS JSON: %v", err)
	}
	if !ok || track == nil {
		t.Fatalf("expected track to be parsed from AIS JSON")
	}

	if track.EntityID != "244670000" {
		t.Errorf("expected EntityID 244670000, got %s", track.EntityID)
	}
	if track.Class != state.TrackClassCivilianMaritime {
		t.Errorf("expected Class CivilianMaritime, got %v", track.Class)
	}
	if !almostEqual(track.Latitude, 51.9244) {
		t.Errorf("expected Latitude 51.9244, got %f", track.Latitude)
	}
	if !almostEqual(track.Longitude, 4.4777) {
		t.Errorf("expected Longitude 4.4777, got %f", track.Longitude)
	}
	if track.AltitudeMeters != 0.0 {
		t.Errorf("expected AltitudeMeters 0 for maritime, got %f", track.AltitudeMeters)
	}
	if !almostEqual(track.VelocityKnots, 14.5) {
		t.Errorf("expected VelocityKnots 14.5, got %f", track.VelocityKnots)
	}
	if !almostEqual(track.HeadingDegrees, 245.0) {
		t.Errorf("expected HeadingDegrees 245.0, got %f", track.HeadingDegrees)
	}
	if !track.HasPosition || !track.HasVelocity {
		t.Errorf("expected HasPosition and HasVelocity to be true for AIS position report")
	}
}

func TestParseAISMessage_Malformed(t *testing.T) {
	// Empty payload
	track, ok, err := ingest.ParseAISMessage([]byte(""))
	if err != nil || ok || track != nil {
		t.Errorf("empty payload should return nil, false, nil")
	}

	// Invalid JSON
	_, _, err = ingest.ParseAISMessage([]byte("NOT_JSON"))
	if err == nil {
		t.Errorf("expected error on invalid JSON")
	}

	// Missing MMSI
	missingMMSI := `{"MessageType":"PositionReport","MetaData":{},"Message":{}}`
	_, _, err = ingest.ParseAISMessage([]byte(missingMMSI))
	if err == nil {
		t.Errorf("expected error on missing MMSI")
	}
}

func TestADSBClient_IntegrationWithMock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	listener := newPipeListener()
	ingest.StartMockADSBWithListener(ctx, listener)

	store := state.NewTrackStore(60 * time.Second)
	defer store.Stop()

	client := ingest.NewADSBClient(ingest.ADSBConfig{
		Endpoint:       "127.0.0.1:5002",
		DialTimeout:    2 * time.Second,
		InitialBackoff: 50 * time.Millisecond,
		MaxBackoff:     200 * time.Millisecond,
	}, store)
	client.NetDial = listener.DialContext

	clientCtx, clientCancel := context.WithCancel(ctx)
	defer clientCancel()

	go func() {
		_ = client.Start(clientCtx)
	}()

	// Wait up to 2 seconds for tracks to be ingested
	deadline := time.Now().Add(2 * time.Second)
	var foundTracks []state.Track
	for time.Now().Before(deadline) {
		tracks := store.GetAllTracks()
		if len(tracks) >= 3 {
			foundTracks = tracks
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	if len(foundTracks) < 3 {
		t.Fatalf("expected at least 3 tracks from mock ADSB, got %d", len(foundTracks))
	}

	// Verify all tracks are CIVILIAN_AIR and have valid coordinates
	for _, tr := range foundTracks {
		if tr.Class != state.TrackClassCivilianAir {
			t.Errorf("track %s has unexpected class: %v", tr.EntityID, tr.Class)
		}
		if tr.Latitude == 0 && tr.Longitude == 0 {
			t.Errorf("track %s has unpopulated coordinates", tr.EntityID)
		}
	}
}

func TestAISClient_IntegrationWithMock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	listener := newPipeListener()
	ingest.StartMockAISWithListener(ctx, listener)

	store := state.NewTrackStore(60 * time.Second)
	defer store.Stop()

	client := ingest.NewAISStreamClient(ingest.AISConfig{
		Endpoint:       "ws://127.0.0.1/v0/stream",
		APIKey:         "TEST_API_KEY",
		InitialBackoff: 50 * time.Millisecond,
		MaxBackoff:     200 * time.Millisecond,
	}, store)
	client.NetDial = listener.DialContext

	clientCtx, clientCancel := context.WithCancel(ctx)
	defer clientCancel()

	go func() {
		_ = client.Start(clientCtx)
	}()

	// Wait up to 2 seconds for maritime tracks to be received
	deadline := time.Now().Add(2 * time.Second)
	var foundTracks []state.Track
	for time.Now().Before(deadline) {
		tracks := store.GetAllTracks()
		if len(tracks) >= 3 {
			foundTracks = tracks
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	if len(foundTracks) < 3 {
		t.Fatalf("expected at least 3 tracks from mock AIS, got %d", len(foundTracks))
	}

	// Verify all tracks are CIVILIAN_MARITIME
	for _, tr := range foundTracks {
		if tr.Class != state.TrackClassCivilianMaritime {
			t.Errorf("track %s has unexpected class: %v", tr.EntityID, tr.Class)
		}
		if tr.Latitude == 0 && tr.Longitude == 0 {
			t.Errorf("track %s has unpopulated coordinates", tr.EntityID)
		}
		if tr.AltitudeMeters != 0 {
			t.Errorf("track %s maritime should have 0 altitude, got %f", tr.EntityID, tr.AltitudeMeters)
		}
	}
}

