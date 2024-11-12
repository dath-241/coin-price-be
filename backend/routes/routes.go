package routers

import (
	servicesA "github.com/dath-241/coin-price-be-go/services/price-service/services/alert"
	servicesI "github.com/dath-241/coin-price-be-go/services/price-service/services/indicator"
	"github.com/gin-gonic/gin"
	 "github.com/dath-241/coin-price-be-go/services/price-service/services/snooze"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	alerts := router.Group("/api/v1/vip2")
	{
		alerts.POST("/alerts", servicesA.CreateAlert)
		alerts.GET("/alerts", servicesA.GetAlerts)
		alerts.GET("/alerts/:id", servicesA.GetAlert)
		alerts.DELETE("/alerts/:id", servicesA.DeleteAlert)
		alerts.GET("/symbol-alerts", servicesA.GetSymbolAlerts)
		alerts.POST("/alerts/symbol", servicesA.SetSymbolAlert)
		
	
        alerts.POST("/start-alert-checker", func(c *gin.Context) {
            services.StartRunning()
            c.JSON(200, gin.H{"status": "Alert checker started"})
        })

    
        alerts.POST("/stop-alert-checker", func(c *gin.Context) {
            services.StopRunning()
            c.JSON(200, gin.H{"status": "Alert checker stopped"})
        })

	}
	indicators := router.Group("/api/v1/vip3/indicators")
	{
		indicators.POST("/", servicesI.SetAdvancedIndicatorAlert)
	}
	return router
}
