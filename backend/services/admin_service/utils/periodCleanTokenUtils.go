package utils

import (
	"time"

	"github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"
)

// Hàm dọn dẹp token hết hạn
func StartCleanupRoutine(interval time.Duration) {
    // interval := 1 * time.Minute  // Đặt tần suất dọn dẹp là mỗi 1 phút
    ticker := time.NewTicker(interval)

    go func() {
        for range ticker.C {
            now := time.Now()
            middlewares.BlacklistedTokensMutex.Lock()
            for token, expTime := range middlewares.BlacklistedTokens {
                if now.After(expTime) {
                    delete(middlewares.BlacklistedTokens, token)
                }
            }
            middlewares.BlacklistedTokensMutex.Unlock()
        }
    }()
}


// func ListBlacklistedTokens(c *gin.Context) {
//     if len(middlewares.BlacklistedTokens) == 0 {
//         c.JSON(http.StatusOK, gin.H{"message": "No blacklisted tokens"})
//         return
//     }

//     // Tạo một slice để lưu các token và thời gian hết hạn
//     var tokens []gin.H
//     for token, expTime := range middlewares.BlacklistedTokens {
//         tokens = append(tokens, gin.H{
//             "token": token,
//             "expires_at": expTime.Format(time.RFC3339), // In ra thời gian hết hạn của token
//         })
//     }

//     // Trả về danh sách token bị blacklist
//     c.JSON(http.StatusOK, gin.H{
//         "blacklisted_tokens": tokens,
//     })
// }
