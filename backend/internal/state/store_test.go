package state_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/TaylorBeck/openlattice/backend/internal/state"
)

func TestTrackStore_BasicUpsertAndGet(t *testing.T) {
	store := state.NewTrackStore(60 * time.Second)
	defer store.Stop()

	track := state.Track{
		EntityID:       "4CA2D6",
		Class:          state.TrackClassCivilianAir,
		Latitude:       51.4700,
		Longitude:      -0.4543,
		AltitudeMeters: 10972.8,
		VelocityKnots:  450.0,
		HeadingDegrees: 180.0,
		TimestampMs:    1729000000000,
	}

	store.UpsertTrack(track)

	retrieved, found := store.GetTrack("4CA2D6")
	if !found {
		t.Fatalf("expected track to be found")
	}

	if retrieved.EntityID != track.EntityID {
		t.Errorf("expected EntityID %s, got %s", track.EntityID, retrieved.EntityID)
	}
	if retrieved.Latitude != track.Latitude {
		t.Errorf("expected Latitude %f, got %f", track.Latitude, retrieved.Latitude)
	}
	if retrieved.Longitude != track.Longitude {
		t.Errorf("expected Longitude %f, got %f", track.Longitude, retrieved.Longitude)
	}
	if retrieved.AltitudeMeters != track.AltitudeMeters {
		t.Errorf("expected AltitudeMeters %f, got %f", track.AltitudeMeters, retrieved.AltitudeMeters)
	}
	if retrieved.VelocityKnots != track.VelocityKnots {
		t.Errorf("expected VelocityKnots %f, got %f", track.VelocityKnots, retrieved.VelocityKnots)
	}
	if retrieved.HeadingDegrees != track.HeadingDegrees {
		t.Errorf("expected HeadingDegrees %f, got %f", track.HeadingDegrees, retrieved.HeadingDegrees)
	}
}

func TestTrackStore_PartialMerging(t *testing.T) {
	store := state.NewTrackStore(60 * time.Second)
	defer store.Stop()

	// 1. Position update arrives first (e.g. MSG,3)
	posUpdate := state.Track{
		EntityID:       "4CA2D6",
		Class:          state.TrackClassCivilianAir,
		Latitude:       51.4700,
		Longitude:      -0.4543,
		AltitudeMeters: 10972.8,
		HasPosition:    true,
	}
	store.UpsertTrack(posUpdate)

	// 2. Velocity update arrives second (e.g. MSG,4)
	velUpdate := state.Track{
		EntityID:       "4CA2D6",
		VelocityKnots:  480.0,
		HeadingDegrees: 270.0,
		HasVelocity:    true,
	}
	store.UpsertTrack(velUpdate)

	merged, ok := store.GetTrack("4CA2D6")
	if !ok {
		t.Fatalf("expected track 4CA2D6 to exist")
	}

	// Verify both position and velocity fields are merged
	if merged.Latitude != 51.4700 || merged.Longitude != -0.4543 || merged.AltitudeMeters != 10972.8 {
		t.Errorf("position fields corrupted after velocity update: %+v", merged)
	}
	if merged.VelocityKnots != 480.0 || merged.HeadingDegrees != 270.0 {
		t.Errorf("velocity fields not updated properly: %+v", merged)
	}
	if merged.Class != state.TrackClassCivilianAir {
		t.Errorf("expected Class CivilianAir, got %v", merged.Class)
	}

	// 3. Reverse order: test velocity first, then position on a new entity
	velFirst := state.Track{
		EntityID:       "BEEF01",
		Class:          state.TrackClassCivilianAir,
		VelocityKnots:  320.0,
		HeadingDegrees: 90.0,
		HasVelocity:    true,
	}
	store.UpsertTrack(velFirst)

	posSecond := state.Track{
		EntityID:       "BEEF01",
		Latitude:       40.7128,
		Longitude:      -74.0060,
		AltitudeMeters: 5000.0,
		HasPosition:    true,
	}
	store.UpsertTrack(posSecond)

	merged2, ok := store.GetTrack("BEEF01")
	if !ok {
		t.Fatalf("expected track BEEF01 to exist")
	}
	if merged2.VelocityKnots != 320.0 || merged2.HeadingDegrees != 90.0 {
		t.Errorf("velocity fields corrupted after position update: %+v", merged2)
	}
	if merged2.Latitude != 40.7128 || merged2.Longitude != -74.0060 || merged2.AltitudeMeters != 5000.0 {
		t.Errorf("position fields not updated properly: %+v", merged2)
	}
}

func TestTrackStore_TTLReaper(t *testing.T) {
	// 50ms TTL, reaped every 20ms
	store := state.NewTrackStoreWithInterval(50*time.Millisecond, 20*time.Millisecond)
	defer store.Stop()

	store.UpsertTrack(state.Track{
		EntityID: "EXPIRE_ME",
		Class:    state.TrackClassCivilianAir,
		Latitude: 10.0,
	})

	if _, ok := store.GetTrack("EXPIRE_ME"); !ok {
		t.Fatalf("track should exist initially")
	}

	// Wait past TTL + reaper cycle
	time.Sleep(120 * time.Millisecond)

	if _, ok := store.GetTrack("EXPIRE_ME"); ok {
		t.Errorf("expected track EXPIRE_ME to have been evicted by TTL reaper")
	}
}

func TestTrackStore_ConcurrentAccessRace(t *testing.T) {
	store := state.NewTrackStore(1 * time.Second)
	defer store.Stop()

	var wg sync.WaitGroup
	workers := 10
	iterations := 200

	for i := 0; i < workers; i++ {
		wg.Add(1)
		workerID := i
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				entityID := fmt.Sprintf("AC-%d", (workerID+j)%5)
				store.UpsertTrack(state.Track{
					EntityID:       entityID,
					Class:          state.TrackClassCivilianAir,
					Latitude:       float64(workerID) + float64(j)*0.01,
					Longitude:      float64(workerID) - float64(j)*0.01,
					AltitudeMeters: float64(j * 100),
					VelocityKnots:  float64(200 + j),
					HeadingDegrees: float64(j % 360),
				})

				_ = store.GetAllTracks()
				_, _ = store.GetTrack(entityID)
				_ = store.GetCountsByClass()
			}
		}()
	}

	wg.Wait()

	if store.Count() == 0 {
		t.Errorf("expected non-zero tracks count after concurrent upserts")
	}
}

