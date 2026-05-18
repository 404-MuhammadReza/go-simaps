package handler

import (
	"go-simaps/internal/model"
	"go-simaps/internal/service"
	"go-simaps/internal/apperror"

	"github.com/gin-gonic/gin"
)

type HospitalHandler struct {
	service service.HospitalService
}

func NewHospitalHandler(service service.HospitalService) *HospitalHandler {
	return &HospitalHandler{service}
}

func (h *HospitalHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	response, err := h.service.GetAll(ctx)
	if err != nil { c.Error(err); return }

	c.JSON(200, Response{
		Success: true,
		Message: "Hospitals retrieved successfully",
		Data: response,
	})
}

func (h *HospitalHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")

	var request model.HospitalRequest
	err := c.ShouldBindJSON(&request)
	if err != nil { c.Error(apperror.ErrValidation(err)); return }

	ctx := c.Request.Context()
	response, err := h.service.Create(ctx, userID, request)
	if err != nil { c.Error(err); return }

	c.JSON(201, Response{
		Success: true,
		Message: "Hospital created successfully",
		Data: response,
	})
}

func (h *HospitalHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	hospitalID := c.Param("id")

	var request model.HospitalRequest
	err := c.ShouldBindJSON(&request)
	if err != nil { c.Error(apperror.ErrValidation(err)); return }

	ctx := c.Request.Context()
	response, err := h.service.Update(ctx, userID, hospitalID, request)
	if err != nil { c.Error(err); return }

	c.JSON(200, Response{
		Success: true,
		Message: "Hospital updated successfully",
		Data: response,
	})
}

func (h *HospitalHandler) Delete(c *gin.Context) {
	hospitalID := c.Param("id")

	ctx := c.Request.Context()
	err := h.service.Delete(ctx, hospitalID)
	if err != nil { c.Error(err); return }

	c.JSON(200, Response{
		Success: true,
		Message: "Hospital deleted successfully",
	})
}

func (h *HospitalHandler) CreateMany(c *gin.Context) {
	userID := c.GetString("userID")
	file, err := c.FormFile("file")
	if err != nil { c.Error(apperror.ErrBadRequest("File is required")); return }

	header := file.Header.Get("Content-Type")
	if header != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		c.Error(apperror.ErrBadRequest("Invalid file type, only .xlsx allowed"))
		return
	}

	excel, err := file.Open()
	if err != nil { c.Error(apperror.ErrInternal("Failed to open file")); return }
	defer excel.Close()

	ctx := c.Request.Context()
	response, err := h.service.CreateMany(ctx, userID, excel)
	if err != nil { c.Error(err); return }

	c.JSON(200, Response{
		Success: true,
		Message: "Hospitals created successfully",
		Data: response,
	})
}

func (h *HospitalHandler) GetReport(c *gin.Context) {
	var payload model.HospitalBulkResponse
	err := c.ShouldBindJSON(&payload)
	if err != nil { c.Error(apperror.ErrBadRequest("Failed to bind payload")); return }

	report, err := h.service.GetReport(payload)
	if err != nil { c.Error(err); return }

	c.Header("Content-Disposition", "attachment; filename=hospital_create_report.xlsx")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", report.Bytes())
}

func (h *HospitalHandler) GetTemplate(c *gin.Context) {
	excel, err := h.service.GetTemplate()
	if err != nil { c.Error(err); return }

	c.Header("Content-Disposition", "attachment; filename=hospital_create_template.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excel.Bytes())
}

func (h *HospitalHandler) Export(c *gin.Context) {
	ctx := c.Request.Context()
	file, err := h.service.Export(ctx)
	if err != nil { c.Error(err); return }

	c.Header("Content-Disposition", "attachment; filename=hospitals.xlsx")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", file.Bytes())
}