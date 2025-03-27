package handler

import (
	"net/http"
	"strconv"

	"github.com/enson89/sustainability-tracker-activity-service/internal/model"
	"github.com/enson89/sustainability-tracker-activity-service/internal/service"
	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	service service.ActivityService
}

func NewActivityHandler(service service.ActivityService) *ActivityHandler {
	return &ActivityHandler{service: service}
}

type CreateActivityRequest struct {
	UserID      int64   `json:"user_id" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Description string  `json:"description"`
}

func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	activity := &model.Activity{
		UserID:      req.UserID,
		Type:        req.Type,
		Amount:      req.Amount,
		Description: req.Description,
	}
	if err := h.service.CreateActivity(activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, activity)
}

func (h *ActivityHandler) GetActivity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	activity, err := h.service.GetActivity(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activity)
}

type UpdateActivityRequest struct {
	Type        string  `json:"type" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Description string  `json:"description"`
}

func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req UpdateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	activity, err := h.service.GetActivity(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	activity.Type = req.Type
	activity.Amount = req.Amount
	activity.Description = req.Description

	if err := h.service.UpdateActivity(activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activity)
}

func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.DeleteActivity(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "activity deleted"})
}
