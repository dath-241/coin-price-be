package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	priceRoutes "github.com/dath-241/coin-price-be-go/services/price-service/routes"
	triggerRoutes "github.com/dath-241/coin-price-be-go/services/trigger-service/routes"
	triggerServiceAlert "github.com/dath-241/coin-price-be-go/services/trigger-service/services/alert"
	"github.com/gin-gonic/gin"

	triggerUtils "github.com/dath-241/coin-price-be-go/services/trigger-service/utils"
	adminRoutes "github.com/dath-241/coin-price-be-go/services/admin_service/routes"
	adminUtils "github.com/dath-241/coin-price-be-go/services/admin_service/utils"

	"github.com/joho/godotenv"
)

func main() {
	server := gin.Default()

	log.Print("Price routes------------------------")
	priceRoutes.RegisterRoutes(server)

	log.Print("Trigger routes------------------------")
	// Tải biến môi trường từ tệp .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Starting server...")
	interval, err := triggerServiceAlert.GetFundingRateInterval("BTCUSDT")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Funding rate interval for BTCUSDT: %s", interval)
	if err := triggerUtils.ConnectMongoDB("mongodb://localhost:27017"); err != nil {
		log.Fatal(err.Error())
	}

	triggerRoutes.SetupRoute(server)
	// triggerR.Run(":3000")
	// Nạp file .env vào môi trường
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Print("Admin routes------------------------")
	// Kết nối MongoDB với retry
	// maxRetries := 3
	// retryDelay := 5 * time.Second
	// if err := adminConfig.ConnectDatabaseWithRetry(maxRetries, retryDelay); err != nil {
	// 	log.Fatalf("Failed to connect to MongoDB: %v", err)
	// }
	// // Đảm bảo ngắt kết nối khi server dừng
	// defer adminConfig.DisconnectDatabase()

	// Bắt đầu routine dọn dẹp token hết hạn
	adminUtils.StartCleanupRoutine()
	adminRoutes.SetupRouter(server)

	// Gọi hàm init trong package momo để khởi tạo các giá trị cần thiết

	// adminMomo.Init()
	//r.GET("/blacklisted-tokens", utils.ListBlacklistedTokens)

	server.Run(":8080")
	// Bắt tín hiệu tắt server để thực hiện cleanup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	log.Println("Server gracefully stopped.")
	// // Chạy server trong một goroutine riêng
	// go func() {
	// 	log.Println("Server is running...")
	// 	if err := adminR.Run(":8082"); err != nil {
	// 		log.Printf("Server exited: %v", err)
	// 	}
	// }()

	// Chờ tín hiệu tắt từ hệ thống

	// go func() {
	//     <-quit
	//     log.Println("Shutting down server...")
	//     os.Exit(0)
	// }()

	// // Chạy server tại cổng 8082
	// if err := r.Run(":8082"); err != nil {
	//     log.Fatalf("Server encountered an error: %v", err)
	// }

}
