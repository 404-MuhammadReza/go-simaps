package handler

import (
	"strings"

	"go-simaps/internal/model"
	"go-simaps/internal/service"
	"go-simaps/internal/apperror"
	"go-simaps/internal/constant"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type EmployeeHandler struct {
	service service.EmployeeService
}

func NewEmployeeHandler(service service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service}
}

func (h *EmployeeHandler) GetAll(c *gin.Context) {
    userRole := c.GetString("userRole")
    userEmail := c.GetString("userEmail")

    var response interface{}
    var err error

	ctx := c.Request.Context()
    if userRole == constant.Admin || userRole == constant.SuperAdmin {
        response, err = h.service.GetAll(ctx)
    } else {
		response, err = h.service.GetByEmail(ctx, userEmail)
	}
    if err != nil { c.Error(err); return }

    c.JSON(200, Response{
        Success: true,
        Message: "Employees retrieved successfully",
        Data:    response,
    })
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")

	var request model.EmployeeRequest
	err := c.ShouldBindJSON(&request)
	if err != nil { c.Error(apperror.ErrValidation(err)); return }

	ctx := c.Request.Context()
	response, err := h.service.Create(ctx, userID, request)
	if err != nil { c.Error(err); return }

	c.JSON(201, Response{
		Success: true,
		Message: "Employee created successfully",
		Data: response,
	})
}

func (h *EmployeeHandler) CreateMany(c *gin.Context) {
	userID := c.GetString("userID")

	file, err := c.FormFile("file")
	if err != nil { c.Error(apperror.ErrBadRequest("File is required")); return }

	header := file.Header.Get("Content-Type")
	if header != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		c.Error(apperror.ErrBadRequest("Invalid file type, only .xlsx allowed"))
		return
	}

	excel, err := file.Open()
	if err != nil { c.Error(apperror.ErrBadRequest("Failed to open file")); return }
	defer excel.Close()

	ctx := c.Request.Context()
	response, err := h.service.CreateMany(ctx, userID, excel)
	if err != nil { c.Error(err); return }

	c.JSON(200, Response{
		Success: true,
		Message: "Bulk employee creation completed",
		Data: response,
	})
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	employeeID := c.Param("id")

	var request model.EmployeeRequest
	err := c.ShouldBindJSON(&request)
	if err != nil { c.Error(apperror.ErrValidation(err)); return }

	ctx := c.Request.Context()
	response, err := h.service.Update(ctx, userID, employeeID, request)
	if err != nil { c.Error(err); return }

	c.JSON(200, Response{
		Success: true,
		Message: "Employee updated successfully",
		Data: response,
	})
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	employeeID := c.Param("id")

	ctx := c.Request.Context()
	err := h.service.Delete(ctx, employeeID)
	if err != nil { c.Error(err); return }

	c.JSON(200, Response{
		Success: true,
		Message: "Employee deleted successfully",
	})
}

func (h *EmployeeHandler) UpdatePicture(c *gin.Context) {
	employeeID := c.Param("id")
	file, err := c.FormFile("file")
	if err != nil { c.Error(apperror.ErrBadRequest("File is required")); return }

	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		c.Error(apperror.ErrBadRequest("Invalid file type, only image files are allowed"))
		return
	}

	size := file.Size
	if size > 1048576 {
		c.Error(apperror.ErrBadRequest("File size exceeds 1MB limit"))
		return
	}

	picture, err := file.Open()
	if err != nil {
		c.Error(apperror.ErrBadRequest("Failed to open uploaded file"))
		return
	}
	defer picture.Close()

	ctx := c.Request.Context()
	response, err := h.service.UpdatePicture(ctx, employeeID, picture, size, contentType)
	if err != nil { c.Error(err); return }

	c.JSON(200, Response{
		Success: true,
		Message: "Employee picture updated successfully",
		Data: response,
	})
}

func (h *EmployeeHandler) DeletePicture(c *gin.Context) {
	employeeID := c.Param("id")

	ctx := c.Request.Context()
	response, err := h.service.DeletePicture(ctx, employeeID)
	if err != nil { c.Error(err); return }

	c.JSON(200, Response{
		Success: true,
		Message: "Employee picture deleted successfully",
		Data: response,
	})
}

func (h *EmployeeHandler) GetTemplate(c *gin.Context) {
	excel, err := h.service.GetTemplate()
	if err != nil { c.Error(err); return }

	c.Header("Content-Disposition", "attachment; filename=employee_create_template.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Writer.Write(excel.Bytes())
}

func (h *EmployeeHandler) GetReport(c *gin.Context) {
	var payload model.EmployeeBulkResponse
	err := c.ShouldBindJSON(&payload)
	if err != nil { c.Error(apperror.ErrBadRequest("Failed to bind payload")); return }

	excel, err := h.service.GetReport(payload)
	if err != nil { c.Error(err); return }

	c.Header("Content-Disposition", "attachment; filename=employee_create_report.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Writer.Write(excel.Bytes())
}

func (h *EmployeeHandler) Export(c *gin.Context) {
	ctx := c.Request.Context()
	excel, err := h.service.Export(ctx)
	if err != nil { c.Error(err); return }

	c.Header("Content-Disposition", "attachment; filename=employees_data.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Writer.Write(excel.Bytes())
}
