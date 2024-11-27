package utils

import (
	"testing"
	"time"
	"github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestStartCleanupRoutine(t *testing.T) {
	tests := []struct {
		name             string
		initialTokens    map[string]time.Time
		expectedTokens   map[string]bool // map các token mong đợi còn lại sau khi dọn dẹp
		interval         time.Duration   // khoảng thời gian giữa các lần kiểm tra
		sleepDuration    time.Duration   // thời gian chờ để đảm bảo hàm thực thi
	}{
		// Test Case 1: Kiểm tra xem token đã hết hạn có bị xóa hay không.
		{
			name: "Expired token should be cleaned up",
			initialTokens: map[string]time.Time{
				"expired_token": time.Now().Add(-time.Minute), // token hết hạn
				"valid_token":   time.Now().Add(10 * time.Minute), // token hợp lệ
			},
			expectedTokens: map[string]bool{
				"valid_token": true,
			},
			interval:      500 * time.Millisecond,
			sleepDuration: 1 * time.Second,
		},
		// Test Case 2: Kiểm tra xem token còn hạn có bị xóa không.
		{
			name: "Valid token should not be cleaned up",
			initialTokens: map[string]time.Time{
				"valid_token": time.Now().Add(10 * time.Minute), // token hợp lệ
			},
			expectedTokens: map[string]bool{
				"valid_token": true, // token hợp lệ phải còn lại
			},
			interval:      500 * time.Millisecond,
			sleepDuration: 1 * time.Second,
		},
		// Test Case 3: Kiểm tra xem khi không có token nào trong BlacklistedTokens, hàm có làm gì không.
		{
			name: "No token to clean up",
			initialTokens: map[string]time.Time{
				// Không có token
			},
			expectedTokens: map[string]bool{
				// Không có token nào nên không có gì thay đổi
			},
			interval:      500 * time.Millisecond,
			sleepDuration: 1 * time.Second,
		},
		// Test Case 4: Kiểm tra xem token hết hạn có bị xóa sau một khoảng thời gian không.
		{
			name: "Expired token should be cleaned after some time",
			initialTokens: map[string]time.Time{
				"expired_token": time.Now().Add(-time.Minute), // token hết hạn
			},
			expectedTokens: map[string]bool{
				// token hết hạn phải bị xóa sau khi chạy
			},
			interval:      500 * time.Millisecond,
			sleepDuration: 2 * time.Second, // Đảm bảo dọn dẹp sẽ diễn ra
		},
		// Test Case 5: Kiểm tra xem nhiều token hết hạn có bị xóa đồng thời không.
		{
			name: "Multiple expired tokens should be cleaned up",
			initialTokens: map[string]time.Time{
				"expired_token_1": time.Now().Add(-time.Minute), // token hết hạn
				"expired_token_2": time.Now().Add(-time.Minute), // token hết hạn
			},
			expectedTokens: map[string]bool{
				// Cả hai token hết hạn đều phải bị xóa
			},
			interval:      500 * time.Millisecond,
			sleepDuration: 2 * time.Second, // Đảm bảo dọn dẹp sẽ diễn ra
		},
	}

	// Lặp qua tất cả các test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Khởi tạo BlacklistedTokens với giá trị ban đầu
			middlewares.BlacklistedTokens = tt.initialTokens

			// Chạy StartCleanupRoutine với interval được định nghĩa trong test case
			go StartCleanupRoutine(tt.interval)

			// Chờ để đảm bảo dọn dẹp token hết hạn
			time.Sleep(tt.sleepDuration)

			// Kiểm tra kết quả sau khi dọn dẹp
			middlewares.BlacklistedTokensMutex.Lock()
			defer middlewares.BlacklistedTokensMutex.Unlock()

			// Kiểm tra từng token mong đợi
			for token, expectedExists := range tt.expectedTokens {
				_, exists := middlewares.BlacklistedTokens[token]
				if expectedExists {
					// Kiểm tra rằng token hợp lệ vẫn còn
					assert.True(t, exists, "Token should be present: %s | Expected: %v | Got: %v", token, expectedExists, exists)
				} else {
					// Kiểm tra rằng token hết hạn đã bị xóa
					assert.False(t, exists, "Token should be cleaned up: %s | Expected: %v | Got: %v", token, expectedExists, exists)
				}
			}
		})
	}
}
