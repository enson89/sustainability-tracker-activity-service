package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/enson89/sustainability-tracker-activity-service/internal/handler"
	"github.com/enson89/sustainability-tracker-activity-service/internal/model"
	"github.com/gin-gonic/gin"
)

// stubActivityService implements the service.ActivityService interface.
type stubActivityService struct{}

func (s *stubActivityService) CreateActivity(activity *model.Activity) error {
	activity.ID = 1
	return nil
}

func (s *stubActivityService) GetActivity(id int64) (*model.Activity, error) {
	return &model.Activity{ID: id, Type: "Test", Amount: 10, Description: "Testing"}, nil
}

func (s *stubActivityService) UpdateActivity(activity *model.Activity) error {
	return nil
}

func (s *stubActivityService) DeleteActivity(id int64) error {
	return nil
}

func TestCreateActivityHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stubSvc := &stubActivityService{}
	activityHandler := handler.NewActivityHandler(stubSvc)
	router := gin.New()
	router.POST("/activities", activityHandler.CreateActivity)

	payload := map[string]interface{}{
		"user_id":     1,
		"type":        "Test",
		"amount":      10,
		"description": "Testing activity",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/activities", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}
