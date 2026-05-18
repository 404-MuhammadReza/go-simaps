package middleware

import (
	"errors"

	"go-simaps/internal/apperror"

	"github.com/gin-gonic/gin"
)

type response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Detail  any    `json:"detail,omitempty"`
}

func Errors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var appErr *apperror.AppError
			if errors.As(err, &appErr) {
				c.AbortWithStatusJSON(appErr.Code, response{
					Success: false,
					Message: appErr.Message,
					Detail:  appErr.Detail,
				})
				return
			}

			c.AbortWithStatusJSON(500, response{
				Success: false,
				Message: "Internal Server Error",
				Detail:  err.Error(),
			})

			return
		}
	}
}
