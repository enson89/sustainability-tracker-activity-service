package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/enson89/sustainability-tracker-activity-service/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestErrorHandlerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, _ := zap.NewDevelopment()
	router := gin.New()
	router.Use(middleware.ErrorHandler(logger))

	// Create a route that returns an error.
	router.GET("/error", func(c *gin.Context) {
		c.Error(http.ErrAbortHandler)
	})

	req, _ := http.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}
