package main

import (
	"log"

	"github.com/dath-241/coin-price-be-go/routes"
	"github.com/dath-241/coin-price-be-go/services/trigger-service/utils"
	"github.com/joho/godotenv"
)

func main() {
	// Tải biến môi trường từ tệp .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := utils.ConnectMongoDB("mongodb://localhost:27017"); err != nil {
		log.Fatal(err.Error())
	}

	r := routes.SetupRoute()
	r.Run(":3000")
}
