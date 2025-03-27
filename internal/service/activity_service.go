package service

import (
	"context"
	"time"

	"github.com/enson89/sustainability-tracker-activity-service/internal/cache"
	"github.com/enson89/sustainability-tracker-activity-service/internal/model"
	"github.com/enson89/sustainability-tracker-activity-service/internal/repository"
)

type ActivityService interface {
	CreateActivity(activity *model.Activity) error
	GetActivity(id int64) (*model.Activity, error)
	UpdateActivity(activity *model.Activity) error
	DeleteActivity(id int64) error
}

type activityService struct {
	repo  repository.ActivityRepository
	cache cache.ActivityCache // Use the interface type here
}

// NewActivityService constructs an ActivityService with repository and cache.
func NewActivityService(repo repository.ActivityRepository, redisCache cache.ActivityCache) ActivityService {
	return &activityService{
		repo:  repo,
		cache: redisCache,
	}
}

func (s *activityService) CreateActivity(activity *model.Activity) error {
	if err := s.repo.CreateActivity(activity); err != nil {
		return err
	}
	ctx := context.Background()
	// Cache for 5 minutes (adjust expiration as needed)
	s.cache.CacheActivity(ctx, activity, 5*time.Minute)
	return nil
}

func (s *activityService) GetActivity(id int64) (*model.Activity, error) {
	ctx := context.Background()
	if cached, err := s.cache.GetCachedActivity(ctx, id); err == nil {
		return cached, nil
	}
	activity, err := s.repo.GetActivityByID(id)
	if err != nil {
		return nil, err
	}
	s.cache.CacheActivity(ctx, activity, 5*time.Minute)
	return activity, nil
}

func (s *activityService) UpdateActivity(activity *model.Activity) error {
	if err := s.repo.UpdateActivity(activity); err != nil {
		return err
	}
	ctx := context.Background()
	s.cache.CacheActivity(ctx, activity, 5*time.Minute)
	return nil
}

func (s *activityService) DeleteActivity(id int64) error {
	if err := s.repo.DeleteActivity(id); err != nil {
		return err
	}
	ctx := context.Background()
	s.cache.InvalidateActivityCache(ctx, id)
	return nil
}
