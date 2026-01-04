package handler

import (
	"net/http"

	"bsnack/internal/service"
	customerror "bsnack/pkg/custom_error"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService service.ReportService
}

func NewReportHandler(reportService service.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

func (h *ReportHandler) GetTransactionReport(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")

	result, err := h.reportService.GetTransactionReport(c.Request.Context(), start, end)
	if err != nil {
		switch err {
		case customerror.ErrInvalidDateRange:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}
