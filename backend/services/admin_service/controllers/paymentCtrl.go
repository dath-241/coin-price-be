package controllers

import (
    "os"
	"context"
	"net/http"
	"time"
	"log"
	"strconv"

	"github.com/dath-241/coin-price-be-go/services/admin_service/config"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"github.com/dath-241/coin-price-be-go/services/admin_service/momo"
	"github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
)

// Khởi tạo thanh toán qua Momo
func CreateVIPPayment() func(*gin.Context) {
    return func(c *gin.Context) {
        //Lấy token từ header Authorization
        // tokenString := c.GetHeader("Authorization")
        // if tokenString == "" {
        //     c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
        //     return
        // }

        // Lấy token từ cookie
        tokenString, err := c.Cookie("accessToken")
        if err != nil || tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization token is required in cookies",
            })
            return
        }

        // Kiểm tra tính hợp lệ của token
        claims, err := middlewares.VerifyJWT(tokenString, true) // true để chỉ định đây là AccessToken
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": err.Error(),
            })
            return
        }

        // Lấy userID và role từ claims
        userID := claims.UserID
        currentVIP := claims.Role
        if userID == "" || currentVIP == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid token claims",
            })
            return
        }

        // Parse JSON body
        var paymentRequest struct {
            Amount   int    `json:"amount"`
            VIPLevel string `json:"vip_level"`
        }
        if err := c.ShouldBindJSON(&paymentRequest); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": err.Error(),
            })
            return
        }

        // Kiểm tra dữ liệu hợp lệ
        if paymentRequest.Amount <= 0 || paymentRequest.VIPLevel == "" {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Invalid request data",
            })
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
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Invalid target VIP level",
            })
            return
        }

        if requestedVIPLevel <= currentVIPLevel {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Invalid request data",
            })
            return
        }

		amountStr := strconv.Itoa(paymentRequest.Amount)
		orderInfo := "Upgrade " + paymentRequest.VIPLevel

        paymentURL, orderId, err := momo.CreateMoMoPayment(amountStr, paymentRequest.VIPLevel, orderInfo)
		if err != nil {
			log.Println("Error creating MoMo payment:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to create payment request",
            })
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
            "transaction_status": "pending",
		}

		collection := config.DB.Collection("OrderMoMo")
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

		_, err = collection.InsertOne(ctx, order)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to create user",
            })
            return
        }

        // Trả về URL thanh toán và order ID
        c.JSON(http.StatusOK, gin.H{
            "payment_url": paymentURL,
            "order_id":    orderId,
        })
    }
}

// Xác nhận thanh toán thành công và cập nhật VIP cho user
func ConfirmPaymentHandlerSuccess() gin.HandlerFunc {
    return func(c *gin.Context) {
        type ConfirmPaymentRequest struct {
            OrderID          string `json:"order_id" binding:"required"`
            TransactionStatus string `json:"transaction_status" binding:"required"`
        }
        var request ConfirmPaymentRequest

        // Bind JSON body vào ConfirmPaymentRequest
        if err := c.ShouldBindJSON(&request); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Invalid request body",
            })
            return
        }

        // Kiểm tra nếu transaction_status là success
        if request.TransactionStatus != "success" {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Payment not successful",
            })
            return
        }
		
        collection := config.DB.Collection("OrderMoMo")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        // Tìm đơn hàng bằng order_id
        var order models.Order // Định nghĩa model order của bạn (với các thông tin như user_id, vip_level, v.v...)
        err := collection.FindOne(context.Background(), bson.M{"order_id": request.OrderID}).Decode(&order)
        if err != nil {
            log.Println("Error finding order:", err)
            c.JSON(http.StatusNotFound, gin.H{
                "error": "Invalid order",
            })
            return
        }

        // Cập nhật trạng thái giao dịch thành công
        update := bson.M{
            "$set": bson.M{
                "transaction_status": "success", // Cập nhật trạng thái giao dịch
            },
        }
        _, err = collection.UpdateOne(ctx, bson.M{"order_id": request.OrderID}, update)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to update transaction status",
            })
            return
        }

        // Lấy thông tin userID và VIP level từ order
        userID := order.UserID
        newVIP := order.VipLevel

        // Cập nhật thông tin user (VIP level) trong collection User
        userCollection := config.DB.Collection("User")

        // Cập nhật VIP level cho người dùng
        update = bson.M{"$set": bson.M{"role": newVIP}}
        _, err = userCollection.UpdateOne(context.Background(), bson.M{"_id": userID}, update)
        if err != nil {
            log.Println("Error updating user role:", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to update user role",
            })
            return
        }

		// Chuyển đổi userID từ ObjectID sang string
		userIDString := userID.Hex()

        // Tạo lại Access Token và Refresh Token với role mới
        accessToken, err := middlewares.GenerateAccessToken(userIDString, newVIP)
        if err != nil {
            log.Println("Error generating new access token:", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to generate access token",
            })
            return
        }

        refreshToken, err := middlewares.GenerateRefreshToken(userIDString, newVIP)
        if err != nil {
            log.Println("Error generating new refresh token:", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to generate refresh token",
            })
            return
        }

        // Lấy token cũ từ cookie hoặc header
        oldAccessToken, err := c.Cookie("accessToken")
        if err != nil {
            log.Println("Error retrieving old access token:", err)
        } else {
            // Xác thực Access Token cũ và thêm vào blacklist nếu hợp lệ
            accessClaims, err := middlewares.VerifyJWT(oldAccessToken, true)
            if err == nil {
                middlewares.BlacklistedTokens[oldAccessToken] = accessClaims.ExpiresAt.Time
            }
        }

        oldRefreshToken, err := c.Cookie("refreshToken")
        if err != nil {
            log.Println("Error retrieving old refresh token:", err)
        } else {
            // Xác thực Refresh Token cũ và thêm vào blacklist nếu hợp lệ
            refreshClaims, err := middlewares.VerifyJWT(oldRefreshToken, false)
            if err == nil {
                middlewares.BlacklistedTokens[oldRefreshToken] = refreshClaims.ExpiresAt.Time
            }
        }

        // Load biến môi trường cho tên miền cookie và thời gian sống
		cookieDomain := os.Getenv("COOKIE_DOMAIN")
		accessTokenTTL := os.Getenv("ACCESS_TOKEN_TTL") // Thời gian sống token
		refreshTokenTTL := os.Getenv("REFRESH_TOKEN_TTL")

		if cookieDomain == "" || accessTokenTTL == "" || refreshTokenTTL == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Environment variables are not set",
            })
			return
		}

		accessTokenTTLInt, err := strconv.Atoi(accessTokenTTL)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Invalid ACCESS_TOKEN_TTL format",
            })
            return
        }

        refreshTokenTTLInt, err := strconv.Atoi(refreshTokenTTL)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Invalid REFRESH_TOKEN_TTL format",
            })
            return
        }

        // Gửi token dưới dạng cookie
        c.SetCookie("accessToken", accessToken, accessTokenTTLInt, "/api/v1", cookieDomain, true, true)  // accessToken cookie
        c.SetCookie("accessToken", accessToken, accessTokenTTLInt, "/auth/logout", cookieDomain, true, true)  // accessToken cookie
        c.SetCookie("refreshToken", refreshToken, refreshTokenTTLInt, "/auth/refresh-token", cookieDomain, true, true) // refreshToken cookie
        c.SetCookie("refreshToken", refreshToken, refreshTokenTTLInt, "/auth/logout", cookieDomain, true, true) // refreshToken cookie


        // Trả về kết quả xác nhận thanh toán thành công
        c.JSON(http.StatusOK, gin.H{
            "message": "Payment confirmed and VIP level upgraded",
        })
    }
}

// HandleQueryPaymentStatus returns a gin.HandlerFunc for querying the payment status
func HandleQueryPaymentStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Define the struct for the incoming request body
		type QueryPaymentRequest struct {
			OrderID   string `json:"orderId"`
			RequestID string `json:"requestId"`
			Lang      string `json:"lang"`
		}

		// Bind the JSON body to the struct
		var req QueryPaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			// Return error if JSON is invalid
			c.JSON(http.StatusBadRequest, gin.H{
                "error": err.Error(),
            })
			return
		}

		// Check if required parameters are missing
		if req.OrderID == "" || req.RequestID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
                "error": "Missing required parameters",
            })
			return
		}

		// Default language is "vi" if not provided
		if req.Lang == "" {
			req.Lang = "vi"
		}

		// Call the function to query the payment status
		result, err := momo.QueryPaymentStatus(req.OrderID, req.RequestID, req.Lang)
		if err != nil {
			log.Println("Error querying payment status:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
                "error": err.Error(),
            })
			return
		}

		// Return the result to the client
		c.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	}
}