package router

import (
	"bsnack/internal/handler"

	"github.com/gin-gonic/gin"
)

func Setup(transactionHandler *handler.TransactionHandler, reportHandler *handler.ReportHandler) *gin.Engine {

	r := gin.Default()

	r.POST("/transactions", transactionHandler.Create)

	r.GET("/reports/transactions", reportHandler.GetTransactionReport)

	return r
}
