package handler

import (
	"go-simaps/internal/model"
	"go-simaps/internal/service"
	"go-simaps/internal/apperror"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetString("userID")

	ctx := c.Request.Context()
	response, err := h.service.GetByID(ctx, userID)
	if err != nil { c.Error(err); return }

	c.JSON(200, Response{
		Success: true,
		Message: "User retrieved successfully",
		Data: response,
	})
}

func (h *UserHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	response, err := h.service.GetAll(ctx)
	if err != nil { c.Error(err); return }

	c.JSON(200, Response{
		Success: true,
		Message: "Users retrieved successfully",
		Data: response,
	})
}

func (h *UserHandler) UpdateRole(c *gin.Context) {
	userID := c.Param("id")

	var request model.RoleRequest
	err := c.ShouldBindJSON(&request)
	if err != nil { c.Error(apperror.ErrValidation(err)); return }

	ctx := c.Request.Context()
	response, err := h.service.UpdateRole(ctx, userID, request)
	if err != nil { c.Error(err); return }

	c.JSON(200, Response{
		Success: true,
		Message: "User role updated successfully",
		Data:    response,
	})
}

func (h *UserHandler) Delete(c *gin.Context) {
	userID := c.Param("id")

	ctx := c.Request.Context()
	err := h.service.Delete(ctx, userID)
	if err != nil { c.Error(err); return }

	c.JSON(200, Response{
		Success: true,
		Message: "User deleted successfully",
	})
}