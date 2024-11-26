package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	priceRoutes "github.com/dath-241/coin-price-be-go/services/price-service/routes"
	triggerRoutes "github.com/dath-241/coin-price-be-go/services/trigger-service/routes"
	"github.com/gin-gonic/gin"

	_ "github.com/dath-241/coin-price-be-go/docs"
	adminConfig "github.com/dath-241/coin-price-be-go/services/admin_service/config"
	adminMomo "github.com/dath-241/coin-price-be-go/services/admin_service/momo"
	adminRoutes "github.com/dath-241/coin-price-be-go/services/admin_service/routes"
	adminUtils "github.com/dath-241/coin-price-be-go/services/admin_service/utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/joho/godotenv"
)

// @title Coin-Price
// @version 1.0
// @description This is a sample server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
func main() {
	server := gin.Default()

	// Tải biến môi trường từ tệp .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Print("Price routes------------------------")
	priceRoutes.RegisterRoutes(server)

	log.Print("Trigger routes------------------------")
	triggerRoutes.SetupRoute(server)

	log.Print("Admin routes------------------------")
	// Kết nối MongoDB với retry
	maxRetries := 3
	retryDelay := 5 * time.Second
	if err := adminConfig.ConnectDatabaseWithRetry(maxRetries, retryDelay); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	// Đảm bảo ngắt kết nối khi server dừng
	defer adminConfig.DisconnectDatabase()

	// Bắt đầu routine dọn dẹp token hết hạn
	adminUtils.StartCleanupRoutine(1 * time.Minute)
	adminRoutes.SetupRouter(server)
	// Gọi hàm init trong package momo để khởi tạo các giá trị cần thiết
	adminMomo.Init()
	//r.GET("/blacklisted-tokens", utils.ListBlacklistedTokens)

	server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.Run(":8080")
	// Bắt tín hiệu tắt server để thực hiện cleanup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	log.Println("Server gracefully stopped.")

}
