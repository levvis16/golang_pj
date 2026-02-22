package handler

import (
	"encoding/json"
	"net/http"

	"subscription-service/internal/logger"
	"subscription-service/internal/models"
	"subscription-service/internal/service"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
	logger  *logger.Logger
}

func NewSubscriptionHandler(service *service.SubscriptionService, logger *logger.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
		logger:  logger,
	}
}

func SetupRouter(h *SubscriptionHandler, logger *logger.Logger) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(LoggerMiddleware(logger))

	api := router.Group("/api/v1")
	{
		subscriptions := api.Group("/subscriptions")
		{
			subscriptions.POST("/", h.Create)
			subscriptions.GET("/", h.List)
			subscriptions.GET("/total", h.GetTotalCost)
			subscriptions.GET("/:id", h.GetByID)
			subscriptions.PUT("/:id", h.Update)
			subscriptions.DELETE("/:id", h.Delete)
		}
	}

	return router
}

func LoggerMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("incoming request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
		)
		c.Next()
	}
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	var req models.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	subscription, err := h.service.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	subscription, err := h.service.GetByID(id)
	if err != nil {
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

func (h *SubscriptionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var req models.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	subscription, err := h.service.Update(id, &req)
	if err != nil {
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *SubscriptionHandler) List(c *gin.Context) {
	filter := &models.SubscriptionFilter{
		UserID:      c.Query("user_id"),
		ServiceName: c.Query("service_name"),
		StartMonth:  c.Query("start_month"),
		EndMonth:    c.Query("end_month"),
	}

	subscriptions, err := h.service.List(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

func (h *SubscriptionHandler) GetTotalCost(c *gin.Context) {
	filter := &models.SubscriptionFilter{
		UserID:      c.Query("user_id"),
		ServiceName: c.Query("service_name"),
		StartMonth:  c.Query("start_month"),
		EndMonth:    c.Query("end_month"),
	}

	total, err := h.service.GetTotalCost(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": total})
}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

func respondWithJSON(c *gin.Context, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	c.Data(code, "application/json", response)
}
