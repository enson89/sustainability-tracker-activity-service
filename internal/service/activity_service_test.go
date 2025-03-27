package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/enson89/sustainability-tracker-activity-service/internal/model"
	"github.com/enson89/sustainability-tracker-activity-service/internal/service"
)

// stubActivityRepo implements repository.ActivityRepository for testing.
type stubActivityRepo struct {
	activity *model.Activity
	err      error
}

func (s *stubActivityRepo) CreateActivity(activity *model.Activity) error {
	if s.err != nil {
		return s.err
	}
	activity.ID = 1
	activity.CreatedAt = time.Now()
	s.activity = activity
	return nil
}
func (s *stubActivityRepo) GetActivityByID(id int64) (*model.Activity, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.activity != nil && s.activity.ID == id {
		return s.activity, nil
	}
	return nil, errors.New("activity not found")
}
func (s *stubActivityRepo) UpdateActivity(activity *model.Activity) error {
	if s.err != nil {
		return s.err
	}
	s.activity = activity
	return nil
}
func (s *stubActivityRepo) DeleteActivity(id int64) error {
	if s.err != nil {
		return s.err
	}
	s.activity = nil
	return nil
}

// stubCache implements cache.ActivityCache.
type stubCache struct {
	stored map[int64]*model.Activity
}

func newStubCache() *stubCache {
	return &stubCache{stored: make(map[int64]*model.Activity)}
}

func (s *stubCache) CacheActivity(ctx context.Context, activity *model.Activity, expiration time.Duration) error {
	s.stored[activity.ID] = activity
	return nil
}

func (s *stubCache) GetCachedActivity(ctx context.Context, id int64) (*model.Activity, error) {
	if act, ok := s.stored[id]; ok {
		return act, nil
	}
	return nil, errors.New("not found")
}

func (s *stubCache) InvalidateActivityCache(ctx context.Context, id int64) error {
	delete(s.stored, id)
	return nil
}

func TestActivityService_CreateAndGet(t *testing.T) {
	repo := &stubActivityRepo{}
	cacheStub := newStubCache() // stubCache implements cache.ActivityCache
	svc := service.NewActivityService(repo, cacheStub)

	activity := &model.Activity{
		UserID:      1,
		Type:        "Test",
		Amount:      20.0,
		Description: "Testing",
	}
	err := svc.CreateActivity(activity)
	if err != nil {
		t.Fatalf("CreateActivity error: %v", err)
	}
	if activity.ID == 0 {
		t.Fatalf("expected activity ID to be set")
	}

	// Test retrieval (should hit the cache).
	got, err := svc.GetActivity(activity.ID)
	if err != nil {
		t.Fatalf("GetActivity error: %v", err)
	}
	if got.Type != "Test" {
		t.Errorf("expected activity type 'Test', got %s", got.Type)
	}
}
