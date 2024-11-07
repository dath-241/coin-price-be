package main

import (
    "coin-price-admin/src/config"
    "coin-price-admin/src/routes"
    "coin-price-admin/src/utils"
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
