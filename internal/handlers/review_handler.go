package handlers

import (
	"net/http"

	"crowdreview/internal/services"
	"crowdreview/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReviewHandler manages review endpoints.
type ReviewHandler struct {
	service services.ReviewService
}

func NewReviewHandler(service services.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

type createReviewRequest struct {
	CompanyID   string `json:"company_id" binding:"required"`
	Rating      int    `json:"rating" binding:"required"`
	Title       string `json:"title"`
	Content     string `json:"content" binding:"required"`
	GeoLocation string `json:"geo_location"`
}

func (h *ReviewHandler) Create(c *gin.Context) {
	var req createReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.JSONError(c, http.StatusUnauthorized, "missing user")
		return
	}
	userID := userIDVal.(uuid.UUID)
	companyID, err := uuid.Parse(req.CompanyID)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid company id")
		return
	}
	review, err := h.service.Create(c.Request.Context(), userID, companyID, services.CreateReviewInput{
		Rating:      req.Rating,
		Title:       req.Title,
		Content:     req.Content,
		IPAddress:   c.ClientIP(),
		GeoLocation: req.GeoLocation,
	})
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSONSuccess(c, http.StatusCreated, review)
}

func (h *ReviewHandler) ListByCompany(c *gin.Context) {
	companyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	reviews, err := h.service.ListByCompany(c.Request.Context(), companyID)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONSuccess(c, http.StatusOK, reviews)
}
