package state

import (
	"sync"
	"time"
)

// TrackClass categorizes entities in the tactical operating picture.
type TrackClass int32

const (
	TrackClassUnknown          TrackClass = 0
	TrackClassCivilianAir      TrackClass = 1
	TrackClassCivilianMaritime TrackClass = 2
	TrackClassThreatAir        TrackClass = 3
	TrackClassThreatMaritime   TrackClass = 4
)

func (c TrackClass) String() string {
	switch c {
	case TrackClassCivilianAir:
		return "CIVILIAN_AIR"
	case TrackClassCivilianMaritime:
		return "CIVILIAN_MARITIME"
	case TrackClassThreatAir:
		return "THREAT_AIR"
	case TrackClassThreatMaritime:
		return "THREAT_MARITIME"
	default:
		return "UNKNOWN"
	}
}

// Track represents a single spatial entity with kinematic and telemetry state.
type Track struct {
	EntityID       string     `json:"entity_id"`
	Class          TrackClass `json:"class"`
	Latitude       float64    `json:"latitude"`
	Longitude      float64    `json:"longitude"`
	AltitudeMeters float64    `json:"altitude_meters"`
	VelocityKnots  float64    `json:"velocity_knots"`
	HeadingDegrees float64    `json:"heading_degrees"`
	TimestampMs    int64      `json:"timestamp_ms"`
	LastSeen       time.Time  `json:"last_seen"`

	// Flags indicating which telemetry components are actively present in this update.
	HasPosition bool `json:"has_position,omitempty"`
	HasVelocity bool `json:"has_velocity,omitempty"`
}

// TrackStore is a thread-safe in-memory store for normalized tactical tracks.
type TrackStore struct {
	mu         sync.RWMutex
	tracks     map[string]Track
	ttl        time.Duration
	reapTicker *time.Ticker
	stopCh     chan struct{}
	stopped    bool
}

// NewTrackStore initializes a TrackStore with the specified TTL (defaults to 60s if <= 0).
func NewTrackStore(ttl time.Duration) *TrackStore {
	if ttl <= 0 {
		ttl = 60 * time.Second
	}
	reapInterval := ttl / 4
	if reapInterval < 10*time.Millisecond {
		reapInterval = 10 * time.Millisecond
	} else if reapInterval > 5*time.Second {
		reapInterval = 5 * time.Second
	}
	return NewTrackStoreWithInterval(ttl, reapInterval)
}

// NewTrackStoreWithInterval allows customizing both TTL and reaper frequency.
func NewTrackStoreWithInterval(ttl, reapInterval time.Duration) *TrackStore {
	if ttl <= 0 {
		ttl = 60 * time.Second
	}
	if reapInterval <= 0 {
		reapInterval = 5 * time.Second
	}

	store := &TrackStore{
		tracks:     make(map[string]Track),
		ttl:        ttl,
		reapTicker: time.NewTicker(reapInterval),
		stopCh:     make(chan struct{}),
	}

	go store.reaperLoop()
	return store
}

// UpsertTrack merges or inserts an incoming track update atomically.
func (s *TrackStore) UpsertTrack(update Track) {
	if update.EntityID == "" {
		return
	}

	now := time.Now()
	if update.LastSeen.IsZero() {
		update.LastSeen = now
	}
	if update.TimestampMs == 0 {
		update.TimestampMs = now.UnixMilli()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.tracks[update.EntityID]
	if !exists {
		s.tracks[update.EntityID] = update
		return
	}

	// Update class if specified
	if update.Class != TrackClassUnknown {
		existing.Class = update.Class
	}

	// Merge position fields if present
	if update.HasPosition || (!update.HasVelocity && (update.Latitude != 0 || update.Longitude != 0 || update.AltitudeMeters != 0)) {
		existing.Latitude = update.Latitude
		existing.Longitude = update.Longitude
		existing.AltitudeMeters = update.AltitudeMeters
		existing.HasPosition = true
	}

	// Merge velocity fields if present
	if update.HasVelocity || (!update.HasPosition && (update.VelocityKnots != 0 || update.HeadingDegrees != 0)) {
		existing.VelocityKnots = update.VelocityKnots
		existing.HeadingDegrees = update.HeadingDegrees
		existing.HasVelocity = true
	}

	existing.TimestampMs = update.TimestampMs
	existing.LastSeen = update.LastSeen

	s.tracks[update.EntityID] = existing
}

// GetTrack returns a track by EntityID.
func (s *TrackStore) GetTrack(entityID string) (Track, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	track, ok := s.tracks[entityID]
	return track, ok
}

// GetAllTracks returns a snapshot copy of all active tracks.
func (s *TrackStore) GetAllTracks() []Track {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Track, 0, len(s.tracks))
	for _, t := range s.tracks {
		result = append(result, t)
	}
	return result
}

// Count returns the number of active tracks in store.
func (s *TrackStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.tracks)
}

// GetCountsByClass returns track counts categorized by class.
func (s *TrackStore) GetCountsByClass() map[TrackClass]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	counts := make(map[TrackClass]int)
	for _, t := range s.tracks {
		counts[t.Class]++
	}
	return counts
}

// DeleteTrack removes a track by entity ID.
func (s *TrackStore) DeleteTrack(entityID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tracks, entityID)
}

// Stop terminates the background TTL reaper goroutine.
func (s *TrackStore) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.stopped {
		s.stopped = true
		s.reapTicker.Stop()
		close(s.stopCh)
	}
}

func (s *TrackStore) reaperLoop() {
	for {
		select {
		case <-s.stopCh:
			return
		case now := <-s.reapTicker.C:
			s.reapStale(now)
		}
	}
}

func (s *TrackStore) reapStale(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := now.Add(-s.ttl)
	for id, t := range s.tracks {
		if t.LastSeen.Before(cutoff) {
			delete(s.tracks, id)
		}
	}
}

