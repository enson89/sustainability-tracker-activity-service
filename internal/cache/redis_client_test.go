package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/enson89/sustainability-tracker-activity-service/internal/cache"
	"github.com/enson89/sustainability-tracker-activity-service/internal/model"
)

func TestRedisClient_CacheAndGetActivity(t *testing.T) {
	// Start a miniredis server.
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	rc, err := cache.NewRedisClient(mr.Addr(), "")
	if err != nil {
		t.Fatalf("NewRedisClient error: %v", err)
	}

	ctx := context.Background()
	activity := &model.Activity{
		ID:          1,
		UserID:      1,
		Type:        "Test",
		Amount:      5.5,
		Description: "Unit test activity",
		CreatedAt:   time.Now(),
	}

	// Cache the activity.
	err = rc.CacheActivity(ctx, activity, 5*time.Minute)
	if err != nil {
		t.Errorf("CacheActivity error: %v", err)
	}

	// Retrieve the cached activity.
	cached, err := rc.GetCachedActivity(ctx, 1)
	if err != nil {
		t.Errorf("GetCachedActivity error: %v", err)
	}
	if cached.Type != "Test" {
		t.Errorf("expected activity type 'Test', got %s", cached.Type)
	}
}
