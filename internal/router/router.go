package router

import (
	"go-simaps/internal/config"
	"go-simaps/internal/handler"
	"go-simaps/internal/constant"
	"go-simaps/internal/apperror"
	"go-simaps/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Auth      *handler.AuthHandler
	User      *handler.UserHandler
	Employee  *handler.EmployeeHandler
	Hospital  *handler.HospitalHandler
	Logging   *handler.LoggingHandler
	Geocoding *handler.GeocodingHandler
}

func Setup(handler Handler, cfg *config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.Errors())

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/dev", handler.Auth.GetDevSession)
			auth.GET("/login", handler.Auth.Login)
			auth.GET("/callback", handler.Auth.Callback)
			auth.POST("/refresh", handler.Auth.Refresh)
			auth.POST("/logout", handler.Auth.Logout)
		}

		protect := api.Group("")
		protect.Use(middleware.Authentication(cfg))
		{
			protect.GET("/me", handler.User.GetMe)
			protect.GET("/employees", handler.Employee.GetAll)
			protect.GET("/hospitals", handler.Hospital.GetAll)
			protect.POST("/logs", handler.Logging.LogUsage)
			geocoding := protect.Group("/geocoding")
			{
				geocoding.POST("/coordinate", handler.Geocoding.ToCoordinate)
				geocoding.POST("/address", handler.Geocoding.ToAddress)
			}

			admin := protect.Group("/admin")
			{
				employees := admin.Group("/employees")
				employees.Use(middleware.Authorization(constant.Admin))
				{
					// employees.GET("", handler.Employee.GetAll)
					employees.POST("", handler.Employee.Create)
					employees.PUT("/:id", handler.Employee.Update)
					employees.DELETE("/:id", handler.Employee.Delete)

					employees.PUT("/:id/picture", handler.Employee.UpdatePicture)
					employees.DELETE("/:id/picture", handler.Employee.DeletePicture)

					employees.POST("/many", handler.Employee.CreateMany)
					employees.POST("/many/report", handler.Employee.GetReport)
					employees.GET("/many/template", handler.Employee.GetTemplate)
					employees.GET("/export", handler.Employee.Export)
				}

				hospitals := admin.Group("/hospitals")
				hospitals.Use(middleware.Authorization(constant.Admin))
				{
					hospitals.POST("", handler.Hospital.Create)
					hospitals.PUT("/:id", handler.Hospital.Update)
					hospitals.DELETE("/:id", handler.Hospital.Delete)

					hospitals.POST("/many", handler.Hospital.CreateMany)
					hospitals.POST("/many/report", handler.Hospital.GetReport)
					hospitals.GET("/many/template", handler.Hospital.GetTemplate)

					hospitals.GET("/export", handler.Hospital.Export)
				}

				user := admin.Group("/users")
				user.Use(middleware.Authorization(constant.SuperAdmin))
				{
					user.GET("", handler.User.GetAll)
					user.PUT("/:id", handler.User.UpdateRole)
					user.DELETE("/:id", handler.User.Delete)
				}

				log := admin.Group("/logs")
				log.Use(middleware.Authorization(constant.SuperAdmin))
				{
					log.GET("", handler.Logging.GetUsageRecap)
					log.GET("/details", handler.Logging.GetFeatureDetails)
					log.GET("/export", handler.Logging.ExportDetails)
				}
			}
		}
	}

	r.HandleMethodNotAllowed = true
	r.NoMethod(func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		c.Error(apperror.ErrMethod("Method " + method + " not allowed for path: " + path))
		c.Abort()
	})

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		c.Error(apperror.ErrNotFound("Route not found: " + path))
		c.Abort()
	})

	return r
}
