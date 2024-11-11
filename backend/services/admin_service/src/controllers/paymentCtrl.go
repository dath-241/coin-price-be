package controllers

import (
	"context"
	"net/http"
	"time"
    "os"
	"log"
	"strconv"

	"backend/services/admin_service/src/config"
	"backend/services/admin_service/src/models"
	"backend/services/admin_service/src/momo"
	"backend/services/admin_service/src/middlewares"

	"github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateVIPPayment() func(*gin.Context) {
    return func(c *gin.Context) {
        //Lấy token từ header Authorization
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            return
        }

        // Kiểm tra tính hợp lệ của token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            return
        }

        userID, ok := claims["user_id"].(string)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
            return
        }

		currentVIP, ok := claims["role"].(string)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in token"})
            return
        }

        // Parse JSON body
        var paymentRequest struct {
            Amount   int    `json:"amount"`
            VIPLevel string `json:"vip_level"`
        }
        if err := c.ShouldBindJSON(&paymentRequest); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Kiểm tra dữ liệu hợp lệ
        if paymentRequest.Amount <= 0 || paymentRequest.VIPLevel == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment data"})
            return
        }

		vipLevels := map[string]int{
            "VIP-0": 0,
            "VIP-1": 1,
            "VIP-2": 2,
            "VIP-3": 3,
        }

        currentVIPLevel := vipLevels[currentVIP]
        requestedVIPLevel, exists := vipLevels[paymentRequest.VIPLevel]

        if !exists {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target VIP level"})
            return
        }

        if requestedVIPLevel <= currentVIPLevel {
            c.JSON(http.StatusBadRequest, gin.H{"error": "New VIP level must be higher than current level"})
            return
        }

		amountStr := strconv.Itoa(paymentRequest.Amount)
		orderInfo := "Upgrade " + paymentRequest.VIPLevel

        paymentURL, orderId, err := momo.CreateMoMoPayment(amountStr, paymentRequest.VIPLevel, orderInfo)
		if err != nil {
			log.Println("Error creating MoMo payment:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment request"})
			return
		}

		// Lưu thông tin đơn hàng vào MongoDB
		order := bson.M{
			"user_id":    userID,
			"vip_level":  paymentRequest.VIPLevel,
			"amount":     paymentRequest.Amount,
			"order_id":   orderId,
			"payment_url": paymentURL,
			"created_at": time.Now(),
		}

		collection := config.DB.Collection("OrderMoMo")

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

		_, err = collection.InsertOne(ctx, order)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
            return
        }

        // Trả về URL thanh toán và order ID
        c.JSON(http.StatusOK, gin.H{
            "payment_url": paymentURL,
            "order_id":    orderId,
        })
    }
}


// ConfirmPaymentHandler xác nhận thanh toán thành công từ MoMo
func ConfirmPaymentHandlerSuccess() gin.HandlerFunc {
    return func(c *gin.Context) {
        type ConfirmPaymentRequest struct {
            OrderID          string `json:"order_id" binding:"required"`
            TransactionStatus string `json:"transaction_status" binding:"required"`
        }
        var request ConfirmPaymentRequest

        // Bind JSON body vào ConfirmPaymentRequest
        if err := c.ShouldBindJSON(&request); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
            return
        }

        // Kiểm tra nếu transaction_status là success
        if request.TransactionStatus != "success" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Payment not successful"})
            return
        }
		
        collection := config.DB.Collection("OrderMoMo")
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        // Tìm đơn hàng bằng order_id
        var order models.Order // Định nghĩa model order của bạn (với các thông tin như user_id, vip_level, v.v...)
        err := collection.FindOne(context.Background(), bson.M{"order_id": request.OrderID}).Decode(&order)
        if err != nil {
            log.Println("Error finding order:", err)
            c.JSON(http.StatusNotFound, gin.H{"error": "Invalid order"})
            return
        }

        // Lấy thông tin userID và VIP level từ order
        userID := order.UserID
        newVIP := order.VipLevel

        // Cập nhật thông tin user (VIP level) trong collection User
        userCollection := config.DB.Collection("User")

        // Cập nhật VIP level cho người dùng
        update := bson.M{"$set": bson.M{"role": newVIP}}
        _, err = userCollection.UpdateOne(context.Background(), bson.M{"_id": userID}, update)
        if err != nil {
            log.Println("Error updating user role:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
            return
        }

		// Xử lý chuyển token cũ vào blacklist
		// // Parse token cũ và lấy expiration time để thêm vào blacklist
		// oldTokenString := c.GetHeader("Authorization")
		// oldToken, err := jwt.Parse(oldTokenString, func(token *jwt.Token) (interface{}, error) {
		// 	return []byte(os.Getenv("JWT_SECRET")), nil
		// })
		// if err != nil || !oldToken.Valid {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		// 	return
		// }

		// claims, ok := oldToken.Claims.(jwt.MapClaims)
		// if !ok {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		// 	return
		// }

		// // Lấy thời gian hết hạn từ token cũ
		// expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)

		// // Thêm token cũ vào danh sách blacklist
		// middlewares.BlacklistedTokens[oldTokenString] = expirationTime

		// Chuyển đổi userID từ ObjectID sang string
		userIDString := userID.Hex()

		// Tạo lại JWT token với role mới
		_, err = middlewares.GenerateJWT(userIDString, newVIP)
		if err != nil {
			log.Println("Error generating new token:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		// Cần cập nhật token mới 

        // Trả về kết quả xác nhận thanh toán thành công
        c.JSON(http.StatusOK, gin.H{
            "message": "Payment confirmed and VIP level upgraded",
        })
    }
}


// func ConfirmPaymentHandler() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		type ConfirmPaymentRequest struct {
// 			OrderID           string `json:"order_id"`
// 			TransactionStatus string `json:"transaction_status"`
// 		}
// 		var req ConfirmPaymentRequest

// 		// Sử dụng Gin để bind JSON vào struct
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			// Trả về lỗi nếu dữ liệu JSON không hợp lệ
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		// Kiểm tra nếu trạng thái giao dịch không phải là "success"
// 		if req.TransactionStatus != "success" {
// 			// Trả về lỗi nếu trạng thái thanh toán không thành công
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction status"})
// 			return
// 		}

// 		amount := 100000 // Giả sử số tiền cần xác nhận (có thể lấy từ cơ sở dữ liệu)
// 		requestType := "capture"

// 		// Gửi yêu cầu xác nhận thanh toán đến MoMo
// 		result, err := momo.ConfirmMoMoPayment(req.OrderID, amount, requestType)
// 		if err != nil {
// 			log.Println("Failed to confirm payment:", err)
// 			// Trả về lỗi nếu có vấn đề trong quá trình xác nhận thanh toán
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm payment"})
// 			return
// 		}

// 		// Trả về thông báo xác nhận thanh toán đã thành công
// 		c.JSON(http.StatusOK, gin.H{
// 			//"message": "Payment confirmed and VIP level upgraded",
// 			"result": result,
// 		})
// 	}
// }

// // HandleQueryPaymentStatus returns a gin.HandlerFunc for querying the payment status
// func HandleQueryPaymentStatus() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Define the struct for the incoming request body
// 		type QueryPaymentRequest struct {
// 			OrderID   string `json:"orderId"`
// 			RequestID string `json:"requestId"`
// 			Lang      string `json:"lang"`
// 		}

// 		// Bind the JSON body to the struct
// 		var req QueryPaymentRequest
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			// Return error if JSON is invalid
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		// Check if required parameters are missing
// 		if req.OrderID == "" || req.RequestID == "" {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
// 			return
// 		}

// 		// Default language is "vi" if not provided
// 		if req.Lang == "" {
// 			req.Lang = "vi"
// 		}

// 		// Call the function to query the payment status
// 		result, err := momo.QueryPaymentStatus(req.OrderID, req.RequestID, req.Lang)
// 		if err != nil {
// 			log.Println("Error querying payment status:", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		// Return the result to the client
// 		c.JSON(http.StatusOK, gin.H{
// 			"result": result,
// 		})
// 	}
// }