package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dath-241/coin-price-be-go/services/admin_service/config"
	"github.com/dath-241/coin-price-be-go/services/admin_service/middlewares"
	"github.com/dath-241/coin-price-be-go/services/admin_service/models"
	"github.com/dath-241/coin-price-be-go/services/admin_service/momo"
	"github.com/dath-241/coin-price-be-go/services/admin_service/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

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

        // Kiểm tra tính hợp lệ của token
        claims, err := middlewares.VerifyJWT(tokenString) // true để chỉ định đây là AccessToken
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
                "error": "Invalid request data",
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

        // Kiểm tra cấp VIP
		vipLevels := map[string]int{
            "VIP-0": 0,
            "VIP-1": 1,
            "VIP-2": 2,
            "VIP-3": 3,
        }

        currentVIPLevel := vipLevels[currentVIP]
        requestedVIPLevel, exists := vipLevels[paymentRequest.VIPLevel]

        if !exists || (requestedVIPLevel <= currentVIPLevel) {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Invalid target VIP level",
            })
            return
        }

        if valid, expectedAmount := utils.IsValidUpgradeCost(currentVIP, paymentRequest.VIPLevel, paymentRequest.Amount); !valid {
            log.Printf("Invalid amount for upgrade: expected %d, got %d\n", expectedAmount, paymentRequest.Amount)

            c.JSON(http.StatusBadRequest, gin.H{
                "error": fmt.Sprintf("Invalid amount. Expected %d for upgrade from %s to %s", expectedAmount, currentVIP, paymentRequest.VIPLevel),
            })
            return
        }

		amountStr := strconv.Itoa(paymentRequest.Amount)
		orderInfo := "Upgrade " + paymentRequest.VIPLevel

        paymentURL, orderId, err := momo.CreateMoMoPayment(amountStr, paymentRequest.VIPLevel, orderInfo)
		if err != nil {
			log.Println("Error creating MoMo payment:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Internal Server Error",
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
                "error": "Internal Server Error",
            })
            return
        }

        // Trả về URL thanh toán và order ID
        c.JSON(http.StatusOK, gin.H{
            "payment_url": paymentURL,
            "order_id":    orderId,
        })

        // Bắt đầu timer 1h40p để tự động thay đổi trạng thái sau 100 phút
        go func() {
            timer := time.NewTimer(100 * time.Minute)
            <-timer.C

            // Cập nhật trạng thái của đơn hàng sau 100 phút
            updateOrderStatusToFailed(orderId)
        }()
    }
}

// Hàm cập nhật trạng thái đơn hàng thành "failed"
func updateOrderStatusToFailed(orderId string) {
    collection := config.DB.Collection("OrderMoMo")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    update := bson.M{"$set": bson.M{"transaction_status": "failed"}}
    _, err := collection.UpdateOne(ctx, bson.M{"order_id": orderId}, update)
    if err != nil {
        log.Printf("Error updating order status for orderId %s: %v", orderId, err)
    } else {
        log.Printf("Order %s status updated to failed due to timeout", orderId)
    }
}

// Xác nhận thanh toán thành công và cập nhật VIP cho user
func confirmPaymentHandlerSuccess(c *gin.Context, OrderID string, Role string){
		
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
            "error": "Internal Server Error",
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
            "error": "Internal Server Error",
        })
        return
    }

    // Chuyển đổi userID từ ObjectID sang string
    userIDString := userID.Hex()

    // Tạo lại Token với role mới
    token, err := middlewares.GenerateToken(userIDString, newVIP)
    if err != nil {
        log.Println("Error generating new access token:", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Internal Server Error",
        })
        return
    }

    if Role == "Admin" {
        c.JSON(http.StatusOK, gin.H{
            "message": "Payment confirmed and VIP level upgraded",
            "status": "0",
        })
        return
    }
 
    // Trả về kết quả xác nhận thanh toán thành công
    c.JSON(http.StatusOK, gin.H{
        "message": "Payment confirmed and VIP level upgraded",
        "status": "0",
        "token": token,
    })
}

// HandleQueryPaymentStatus godoc
// @Summary Check payment status and upgrade user's VIP level if successful
// @Description Queries the payment status and upgrades the user's VIP level if the payment is successful based on the order details
// @Tags Payment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param statusRequest body models.QueryPaymentRequest true "Order ID from the payment gateway"
// @Success 200 {object} models.ReponseQueryPaymentRequest "Payment confirmed and VIP level upgraded"
// @Failure 400 {object} models.ErrorResponse "Invalid order ID or missing parameters"
// @Failure 403 {object} models.ErrorResponse "You do not have permission to query this order"
// @Failure 404 {object} models.ErrorResponse "Order not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error during payment confirmation"
// @Router /api/v1/payment/status [post]
// Check payment status and upgrade user's VIP level if successful
func HandleQueryPaymentStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

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

		// Lấy token từ header Authorization
		tokenString := c.GetHeader("Authorization")

		var userIDFromToken, role string

		// Nếu token được cung cấp, xác thực và lấy thông tin user_id, role
		if tokenString != "" {
			claims, err := middlewares.VerifyJWT(tokenString) // true để chỉ định đây là AccessToken
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid or expired token",
				})
				return
			}
			userIDFromToken = claims.UserID
			role = claims.Role
		}

		// Nếu không có token và không phải admin, từ chối truy cập
		if tokenString == "" && role != "Admin" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token is required",
			})
			return
		}

        // Lấy thông tin đơn hàng từ MongoDB bằng OrderID
		collection := config.DB.Collection("OrderMoMo")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var order models.Order
		err := collection.FindOne(ctx, bson.M{"order_id": req.OrderID}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Order not found",
			})
			return
		}

        // Nếu không phải admin, kiểm tra quyền sở hữu đơn hàng
		if role != "Admin" {
			// Chuyển userIDFromToken từ string sang ObjectID
			userIDObj, err := primitive.ObjectIDFromHex(userIDFromToken)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid user ID in token",
				})
				return
			}

			if order.UserID != userIDObj {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "You do not have permission to query this order",
				})
				return
			}
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
                "error": "Internal Server Error",
            })
			return
		}

        // Kiểm tra trạng thái giao dịch từ MoMo
        if result.ResultCode == 0 {
            // Gọi hàm xử lý thành công
            confirmPaymentHandlerSuccess(c, req.OrderID, role)
        } else {
            c.JSON(http.StatusOK, gin.H{
                "message": "Transaction is not successful yet", 
                "status": result.Message,
            })
        }
	}
}