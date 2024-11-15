package main

import (
	"log"

	"github.com/dath-241/coin-price-be-go/routes"
	services "github.com/dath-241/coin-price-be-go/services/trigger-service/services/alert"
	"github.com/dath-241/coin-price-be-go/services/trigger-service/utils"
	"github.com/joho/godotenv"
)

func main() {
	// Tải biến môi trường từ tệp .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Starting server...")
	interval, err := services.GetFundingRateInterval("BTCUSDT")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Funding rate interval for BTCUSDT: %s", interval)
	if err := utils.ConnectMongoDB("mongodb://localhost:27017"); err != nil {
		log.Fatal(err.Error())
	}

	r := routes.SetupRoute()
	r.Run(":3000")
}
