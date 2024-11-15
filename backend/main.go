package main

import (
	// "github.com/dath-241/coin-price-be-go/services/price-service/routes"
	"github.com/dath-241/coin-price-be-go/services/price-service/src/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	//Load env
	//Connect db
}

func main() {

	server := gin.Default()

	routes.RegisterRoutes(server)
	// routes.RegisterRoutes(server)

	server.Run(":8080")
	//Set up router, routs, server, websocket
}
