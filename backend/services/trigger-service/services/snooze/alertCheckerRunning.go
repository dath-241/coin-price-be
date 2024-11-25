package services

import (
	"log"
	"sync"
	"time"

	_ "github.com/dath-241/coin-price-be-go/services/trigger-service/models"
	"github.com/gin-gonic/gin"
)

var (
	ticker    *time.Ticker
	stop      chan bool
	isRunning bool
	mutex     sync.Mutex
)

// Run starts the alert checker.
// @Summary Start alert checker
// @Description Starts the alert checker to monitor for alerts
// @Tags Alert Running
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseAlertCheckerStatus "Alert checker started successfully"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/vip2/start-alert-checker [post]
func Run(c *gin.Context) {
	StartRunning()
	c.JSON(200, gin.H{"status": "Alert checker started"})
}

// Stop stops the alert checker.
// @Summary Stop alert checker
// @Description Stops the alert checker from monitoring for alerts
// @Tags Alert Running
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseAlertCheckerStatus "Alert checker stopped successfully"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/vip2/stop-alert-checker [post]
func Stop(c *gin.Context) {
	StopRunning()
	c.JSON(200, gin.H{"status": "Alert checker stopped"})
}

func StartRunning() {
	mutex.Lock()
	defer mutex.Unlock()

	if isRunning {
		log.Println("Alert checker is already running.")
		return
	}

	stop = make(chan bool)
	ticker = time.NewTicker(1 * time.Second)
	isRunning = true

	go func() {
		for {
			select {
			case <-ticker.C:
				CheckAndSendAlerts()
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func StopRunning() {
	mutex.Lock()
	defer mutex.Unlock()

	if !isRunning {
		return
	}

	stop <- true
	isRunning = false
}
