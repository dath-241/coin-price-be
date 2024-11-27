package controllers

import (
	"context"
	"net/http"
	"time"
	"log"
	"strconv"

	"github.com/dath-241/coin-price-be-go/services/admin_service/config"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"github.com/dath-241/coin-price-be-go/services/admin_service/momo"
	"github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"
    "go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
)

// CreateVIPPayment godoc
// @Summary Initiate MoMo payment for VIP upgrade
// @Description Creates a MoMo payment request for upgrading the user's VIP level, validates the token, and stores the order details in the database
// @Tags Payment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param paymentRequest body models.CreateVIPPaymentRequest true "Payment request data"
// @Success 200 {object} models.CreateVIPPaymentReponse "Payment URL and Order ID"
// @Failure 400 {object} models.ErrorResponse "Invalid request data or missing parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized: Invalid or missing authorization token"
// @Failure 500 {object} models.ErrorResponse "Internal server error during payment creation"
// @Router /api/v1/payment/vip-upgrade [post]
// Khởi tạo thanh toán 
func CreateVIPPayment() func(*gin.Context) {
    return func(c *gin.Context) {
        //Lấy token từ header Authorization
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            return
        }

        // // Lấy token từ cookie
        // tokenString, err := c.Cookie("accessToken")
        // if err != nil || tokenString == "" {
        //     c.JSON(http.StatusUnauthorized, gin.H{
        //         "error": "Authorization token is required in cookies",
        //     })
        //     return
        // }

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
            "orderInfo":  orderInfo,
			"payment_url": paymentURL,
			"created_at": primitive.NewDateTimeFromTime(time.Now()),
            "updated_at": primitive.NewDateTimeFromTime(time.Now()),
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
func confirmPaymentHandlerSuccess(c *gin.Context, OrderID string){
		
    collection := config.DB.Collection("OrderMoMo")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Tìm đơn hàng bằng order_id
    var order models.Order // Định nghĩa model order của bạn (với các thông tin như user_id, vip_level, v.v...)
    err := collection.FindOne(context.Background(), bson.M{"order_id": OrderID}).Decode(&order)
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
            "updated_at": time.Now(),
        },
    }
    _, err = collection.UpdateOne(ctx, bson.M{"order_id": OrderID}, update)
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

    // Gửi token dưới dạng cookie
    setAuthCookies(c, accessToken, refreshToken, false, true)

    // Trả về kết quả xác nhận thanh toán thành công
    c.JSON(http.StatusOK, gin.H{
        "message": "Payment confirmed and VIP level upgraded",
        "status": "0",
        "token": accessToken,
    })
}

// HandleQueryPaymentStatus godoc
// @Summary Check payment status and upgrade user's VIP level if successful
// @Description Queries the payment status and upgrades the user's VIP level if the payment is successful based on the order details
// @Tags Payment
// @Accept json
// @Produce json
// @Param statusRequest body models.QueryPaymentRequest true "Order ID from the payment gateway"
// @Success 200 {object} models.ReponseQueryPaymentRequest "Payment confirmed and VIP level upgraded"
// @Failure 400 {object} models.ErrorResponse "Invalid order ID or missing parameters"
// @Failure 404 {object} models.ErrorResponse "Order not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error during payment confirmation"
// @Router /api/v1/payment/status [post]
// Check payment status and upgrade user's VIP level if successful
func HandleQueryPaymentStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Define the struct for the incoming request body

		// Bind the JSON body to the struct
		var req models.QueryPaymentRequest
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

        // Kiểm tra trạng thái giao dịch từ MoMo
        if result.ResultCode == 0 {
            // Gọi hàm xử lý thành công
            confirmPaymentHandlerSuccess(c, req.OrderID)
        } else {
            c.JSON(http.StatusOK, gin.H{
                "message": "Transaction is not successful yet", 
                "status": result.Message,
            })
        }
	}
}