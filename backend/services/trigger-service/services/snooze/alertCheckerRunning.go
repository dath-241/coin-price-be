package services

import (
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	ticker    *time.Ticker
	stop      chan bool
	isRunning bool
	mutex     sync.Mutex
)

func Run(c *gin.Context) {
	StartRunning()
	c.JSON(200, gin.H{"status": "Alert checker started"})
}
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
