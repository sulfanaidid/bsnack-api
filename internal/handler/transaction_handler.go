package handler

import (
	"bsnack/internal/model"
	"bsnack/internal/service"
	customerror "bsnack/pkg/custom_error"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	Service *service.TransactionService
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var req model.CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.CreateTransaction(c.Request.Context(), req)
	if err != nil {
		switch err {
		case customerror.ErrItemsEmpty:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case customerror.ErrProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case customerror.ErrStockNotEnough:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, res)
}
