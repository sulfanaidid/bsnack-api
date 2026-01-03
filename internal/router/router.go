package router

import (
	"github.com/gin-gonic/gin"

	"bsnack/internal/handler"
)

func Setup(h *handler.TransactionHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/transactions", h.Create)
	return r
}
