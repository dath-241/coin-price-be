package utils

import (
    "time"
	"backend/services/admin_service/src/middlewares"
)


// Hàm dọn dẹp token hết hạn
func StartCleanupRoutine() {
    ticker := time.NewTicker(1 * time.Minute)
    go func() {
        for range ticker.C {
            now := time.Now()
            for token, expTime := range  middlewares.BlacklistedTokens {
                if now.After(expTime) {
                    delete(middlewares.BlacklistedTokens, token)
                }
            }
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