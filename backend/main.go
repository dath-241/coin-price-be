package main

import (
    "backend/services/admin_service/src/config"
    "backend/services/admin_service/src/routes"
    "backend/services/admin_service/src/utils"
    "github.com/joho/godotenv"
    "log"

)

func main() {
    // Nạp file .env vào môi trường
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    config.ConnectDatabase()


    // Bắt đầu routine dọn dẹp token hết hạn
    utils.StartCleanupRoutine()
    r := routes.SetupRouter()
    //r.GET("/blacklisted-tokens", utils.ListBlacklistedTokens)

    r.Run(":8082") // Chạy server tại cổng 8080
}
