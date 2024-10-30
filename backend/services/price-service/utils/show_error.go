package utils

import "github.com/gin-gonic/gin"

func ShowError(statusCode int64, message string, context *gin.Context) {
	context.JSON(int(statusCode), gin.H{"message": message})
}
