package handlers

import (
	"net/http"

	"crowdreview/config"
	"crowdreview/internal/services"
	"crowdreview/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthHandler exposes auth endpoints.
type AuthHandler struct {
	auth   services.AuthService
	config config.Config
}

func NewAuthHandler(auth services.AuthService, cfg config.Config) *AuthHandler {
	return &AuthHandler{auth: auth, config: cfg}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	user, access, refresh, err := h.auth.Register(c.Request.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSONSuccess(c, http.StatusCreated, gin.H{
		"user":          user,
		"access_token":  access,
		"refresh_token": refresh,
	})
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	user, access, refresh, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		utils.JSONError(c, http.StatusUnauthorized, err.Error())
		return
	}
	utils.JSONSuccess(c, http.StatusOK, gin.H{
		"user":          user,
		"access_token":  access,
		"refresh_token": refresh,
	})
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := h.auth.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		utils.JSONError(c, http.StatusUnauthorized, "invalid refresh token")
		return
	}
	access, refresh, err := h.auth.Refresh(c.Request.Context(), userID)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONSuccess(c, http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}
