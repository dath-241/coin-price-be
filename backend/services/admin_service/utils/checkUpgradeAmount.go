package utils

// Map định nghĩa bảng giá nâng cấp
var upgradeCosts = map[string]map[string]int{
    "VIP-0": {
        "VIP-1": 79000,
        "VIP-2": 129000,
        "VIP-3": 149000,
    },
    "VIP-1": {
        "VIP-2": 79000,
        "VIP-3": 99000,
    },
    "VIP-2": {
        "VIP-3": 79000,
    },
}

// Hàm kiểm tra số tiền
func IsValidUpgradeCost(currentVIP, targetVIP string, amount int) (bool, int) {
    // Kiểm tra cấp VIP hiện tại và cấp VIP mục tiêu có tồn tại không
    if costs, exists := upgradeCosts[currentVIP]; exists {
        if expectedAmount, valid := costs[targetVIP]; valid {
            return amount == expectedAmount, expectedAmount
        }
    }
    // Nếu không tìm thấy trong bảng giá
    return false, 0
}
