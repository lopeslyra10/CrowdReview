package handlers

import (
	"net/http"

	"crowdreview/internal/models"
	"crowdreview/internal/services"
	"crowdreview/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CompanyHandler manages company endpoints.
type CompanyHandler struct {
	service services.CompanyService
}

func NewCompanyHandler(service services.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: service}
}

func (h *CompanyHandler) List(c *gin.Context) {
	companies, err := h.service.List(c.Request.Context())
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONSuccess(c, http.StatusOK, companies)
}

func (h *CompanyHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	company, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		utils.JSONError(c, http.StatusNotFound, err.Error())
		return
	}
	utils.JSONSuccess(c, http.StatusOK, company)
}

func (h *CompanyHandler) Create(c *gin.Context) {
	var input models.Company
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	company, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSONSuccess(c, http.StatusCreated, company)
}

func (h *CompanyHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var input models.Company
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	company, err := h.service.Update(c.Request.Context(), id, input)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSONSuccess(c, http.StatusOK, company)
}
