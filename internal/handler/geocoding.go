package handler

import (
	"go-simaps/internal/model"
	"go-simaps/internal/apperror"
	"go-simaps/internal/service"

	"github.com/gin-gonic/gin"
)

type GeocodingHandler struct {
	service service.GeocodingService
}

func NewGeocodingHandler(service service.GeocodingService) *GeocodingHandler {
	return &GeocodingHandler{service}
}

func (h *GeocodingHandler) ToCoordinate(c *gin.Context) {
	var request struct { Address string `json:"address" binding:"required"` }
	err := c.ShouldBindJSON(&request);
	if err != nil { c.Error(apperror.ErrValidation(err)); return }

	ctx := c.Request.Context()
	coordinate, err := h.service.ToCoordinate(ctx, request.Address)
	if err != nil { c.Error(err); return }

	c.JSON(200, Response{
		Success: true, 
		Message: "Coordinates retrieved successfully",
		Data: coordinate,
	})
}

func (h *GeocodingHandler) ToAddress(c *gin.Context) {
	var request model.Coordinate
	err := c.ShouldBindJSON(&request);
	if err != nil { c.Error(apperror.ErrValidation(err)); return }

	ctx := c.Request.Context()
	address, err := h.service.ToAddress(ctx, &request)
	if err != nil { c.Error(err); return }

	c.JSON(200, Response{
		Success: true, 
		Message: "Address retrieved successfully",
		Data: address,
	})
}