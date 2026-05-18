package handler

import (
	"go-simaps/internal/config"
	"go-simaps/internal/service"
	"go-simaps/internal/apperror"
	"go-simaps/internal/constant"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service 		 service.AuthService
	AppDomain        string
	JWTAccessMaxAge  int
	JWTRefreshMaxAge int
}

func NewAuthHandler(service service.AuthService, config *config.Config) *AuthHandler {
	return &AuthHandler{
		service: service,
		AppDomain: config.AppDomain,
		JWTAccessMaxAge: config.JWTAccessMaxAge,
		JWTRefreshMaxAge: config.JWTRefreshMaxAge,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	url, state, err := h.service.GenerateLoginURL()
	if err != nil { c.Error(err); return }

	c.SetCookie(constant.OAuthState, state, 3600, "/", h.AppDomain, false, true)
	c.JSON(200, Response{
		Success: true,
		Message: "Login URL generated successfully",
		Data: url,
	})
}

func (h *AuthHandler) Callback(c *gin.Context) {
	queryState := c.Query("state")
	cookieState, err := c.Cookie(constant.OAuthState)
	c.SetCookie(constant.OAuthState, "", -1, "/", h.AppDomain, false, true)
	if err != nil || cookieState != queryState {
		c.Error(apperror.ErrUnauthorized("Invalid state parameter"))
		return
	}

	code := c.Query("code")
	if code == "" { c.Error(apperror.ErrBadRequest("Code is required")); return }

	ctx := c.Request.Context()
	accessToken, refreshToken, err := h.service.ExchangeToken(ctx, code)
	if err != nil { c.Error(err); return }

	c.SetCookie(constant.AccessToken, accessToken, h.JWTAccessMaxAge, "/", h.AppDomain, false, true)
	c.SetCookie(constant.RefreshToken, refreshToken, h.JWTRefreshMaxAge, "/", h.AppDomain, false, true)
	c.SetCookie(constant.Session, "true", h.JWTRefreshMaxAge, "/", h.AppDomain, false, false)

	c.JSON(200, Response{
		Success: true,
		Message: "Login successful",
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(constant.RefreshToken)
	if err != nil || refreshToken == "" {
		c.Error(apperror.ErrUnauthorized("Refresh token missing"))
		return
	}

	ctx := c.Request.Context()
	newAccessToken, err := h.service.RefreshToken(ctx, refreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookie(constant.AccessToken, newAccessToken, h.JWTAccessMaxAge, "/", h.AppDomain, false, true)
	c.JSON(200, Response{
		Success: true,
		Message: "Token refreshed successfully",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie(constant.AccessToken, "", -1, "/", h.AppDomain, false, true)
	c.SetCookie(constant.RefreshToken, "", -1, "/", h.AppDomain, false, true)
	c.SetCookie(constant.Session, "", -1, "/", h.AppDomain, false, false)

	c.JSON(200, Response{
		Success: true,
		Message: "Logged out successfully",
	})
}

// For Development Only
func (h *AuthHandler) GetDevSession(c *gin.Context) {
	userID := "TI_USER_17786367945577"

	ctx := c.Request.Context()
	accessToken, err := h.service.DevSession(ctx, userID)
	if err != nil { c.Error(err); return }

	c.SetCookie(constant.AccessToken, accessToken, h.JWTRefreshMaxAge, "/", h.AppDomain, false, true)
	c.SetCookie(constant.Session, "true", h.JWTRefreshMaxAge, "/", h.AppDomain, false, false)
	c.JSON(200, Response{
		Success: true,
		Message: "Session created successfully",
	})
}