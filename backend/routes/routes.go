package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/api")
	authenticated.GET("/v1/funding-rate", getFundingRate)
}
