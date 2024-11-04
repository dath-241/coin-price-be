package routers

import (
	"github.com/dath-241/coin-price-be-go/services/price-service/services/alert"
	
	"github.com/gin-gonic/gin"
)


func SetupRouter() *gin.Engine {
	router := gin.Default()




	alerts := router.Group("/api/v1/vip2/alerts")
    {
        alerts.POST("/", services.CreateAlert)
        alerts.GET("/", services.GetAlerts)
		alerts.GET("/:id", services.GetAlert)
        alerts.DELETE("/:id", services.DeleteAlert)
		alerts.GET("/symbol-alerts",services.GetSymbolAlerts)
    }
	return router
}
