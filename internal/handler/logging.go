package handler

import (
	"context"

	"go-simaps/internal/model"
	"go-simaps/internal/service"
	"go-simaps/internal/apperror"

	"github.com/gin-gonic/gin"
)

type LoggingHandler struct {
	service service.LoggingService
}

func NewLoggingHandler(service service.LoggingService) *LoggingHandler {
	return &LoggingHandler{service}
}

func (h *LoggingHandler) LogUsage(c *gin.Context) {
	userID := c.GetString("userID")

	var request model.LogRequest
	err := c.ShouldBindJSON(&request)
	if err != nil { c.Error(apperror.ErrValidation(err)); return }

	ctx := context.Background()
	err = h.service.LogUsage(ctx, userID, request)
	if err != nil { c.Error(err); return }

	c.JSON(200, Response{
		Success: true,
		Message: "Usage logged successfully",
	})
}

func (h *LoggingHandler) GetUsageRecap(c *gin.Context) {
	ctx := context.Background()
	recaps, err := h.service.GetSummary(ctx)
	if err != nil { c.Error(err); return }

	c.JSON(200, Response{
		Success: true,
		Message: "Usage recap retrieved successfully",
		Data: recaps,
	})
}

func (h *LoggingHandler) GetFeatureDetails(c *gin.Context) {
	var query model.DetailRequest

	err := c.ShouldBindQuery(&query)
	if err != nil { c.Error(apperror.ErrValidation(err)); return }

	ctx := context.Background()
	logs, err := h.service.GetDetails(ctx, query)
	if err != nil { c.Error(err);  return }

	c.JSON(200, Response{
		Success: true,
		Message: "Feature details retrieved successfully",
		Data: logs,
	})
}
