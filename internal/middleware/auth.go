package middleware

import (
	"go-simaps/internal/utils"
	"go-simaps/internal/config"
	"go-simaps/internal/apperror"
	"go-simaps/internal/constant"

	"github.com/gin-gonic/gin"
)

func Authentication(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(constant.AccessToken)
		if token == "" || err != nil {
			c.Error(apperror.ErrUnauthorized("Authentication required"))
			c.Abort()
			return
		}

		claims, err := utils.ValidateAccessToken(token, cfg.JWTSecret)
		if err != nil {
			c.Error(apperror.ErrUnauthorized("Invalid or expired token"))
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.UserEmail)
		c.Set("userRole", claims.UserRole)
		c.Next()
	}
}

func Authorization(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roleVal, ok := ctx.Get("userRole")
        if !ok {
            ctx.Error(apperror.ErrForbidden("Role not found in token"))
            ctx.Abort()
            return
        }

		userRole, ok := roleVal.(string)
        if !ok {
            ctx.Error(apperror.ErrForbidden("Invalid role format"))
            ctx.Abort()
            return
        }

		if userRole == requiredRole || userRole == constant.SuperAdmin {
            ctx.Next()
            return
        }

        ctx.Error(apperror.ErrForbidden("You don't have permission to access this resource"))
        ctx.Abort()
	}
}
