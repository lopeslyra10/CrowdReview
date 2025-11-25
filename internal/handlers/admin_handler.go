package handlers

import (
	"net/http"

	"crowdreview/internal/services"
	"crowdreview/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AdminHandler exposes admin-only endpoints.
type AdminHandler struct {
	service services.AdminService
}

func NewAdminHandler(service services.AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

func (h *AdminHandler) Insights(c *gin.Context) {
	insights, err := h.service.GetInsights(c.Request.Context())
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONSuccess(c, http.StatusOK, insights)
}

func (h *AdminHandler) Suspicious(c *gin.Context) {
	reviews, err := h.service.ListSuspicious(c.Request.Context())
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONSuccess(c, http.StatusOK, reviews)
}

type respondRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *AdminHandler) Respond(c *gin.Context) {
	var req respondRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	reviewID := c.Param("id")
	if err := h.service.Respond(c.Request.Context(), reviewID, req.Status); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSONSuccess(c, http.StatusOK, gin.H{"status": "updated"})
}
