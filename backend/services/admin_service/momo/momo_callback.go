package momo

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	//"github.com/dath-241/coin-price-be-go/services/admin_service/controllers"
	"github.com/gin-gonic/gin"
)

func MoMoCallback() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Lấy reponse của momo lưu vào callbackData.
		var callbackData map[string]interface{}
		if err := c.ShouldBindJSON(&callbackData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON payload",
			})
			return
		}

		// Lấy chữ ký từ MoMo gửi về 	
		signature := callbackData["signature"].(string)
		fmt.Println("Received signature:", signature)


		params := getCallbackParams(callbackData)

		// Xây dựng chuỗi raw signature theo thứ tự a-z
		rawSignature := buildRawSignature(accessKeyEnv, params)

		// In raw signature để kiểm tra
		fmt.Println("Raw signature string:", rawSignature)

		// Tính toán chữ ký HMAC_SHA256
		calculatedSignature := calculateSignature(rawSignature, secretKeyEnv)

		// In chữ ký tính toán để so sánh
		fmt.Println("Calculated signature:", calculatedSignature)

		// So sánh chữ ký nhận được và chữ ký đã tính toán
		if calculatedSignature != signature {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid signature",
			})
			return
		}

		// Kiểm tra kết quả thanh toán
		if params["resultCode"] == "0" {
			// Thanh toán thành công
			fmt.Println("Thanh toan thanh cong.")
			// Cap nhat VIP cho user dua tren orderID
			//controllers.ConfirmPaymentHandlerSuccess()
		} else {
			// Thanh toán thất bại
			fmt.Println("Thanh toan that bai.")
		}

		// Gửi phản hồi về momo
		// Xây dựng chuỗi raw signature responce theo thứ tự a-z
		rawSignatureResponce := buildRawSignatureResponce(accessKeyEnv, params)
		// Tính toán chữ ký HMAC_SHA256
		calculatedSignatureResponce := calculateSignature(rawSignatureResponce, secretKeyEnv)
		response := gin.H{
			"partnerCode":  partnerCodeEnv,
			"requestId":    params["requestId"],
			"orderId":      params["orderId"],
			"resultCode":   params["resultCode"],
			"message":      params["message"],
			"responseTime": params["responseTime"],
			"extraData":    params["extraData"],
			"signature":	calculatedSignatureResponce,
		}
		c.JSON(http.StatusOK, response)
	}
}

// Lấy các tham số từ callbackData và trả về đúng thứ tự a-z
func getCallbackParams(callbackData map[string]interface{}) map[string]string {
	return map[string]string{
		//"accessKey":    fmt.Sprintf("%v", os.Getenv("MOMO_ACCESS_KEY")),  // Thêm accessKey nếu cần
		"amount":       fmt.Sprintf("%v", int64(callbackData["amount"].(float64))),
		"extraData":    fmt.Sprintf("%v", callbackData["extraData"]),
		"message":      fmt.Sprintf("%v", callbackData["message"]),
		"orderId":      fmt.Sprintf("%v", callbackData["orderId"]),
		"orderInfo":    fmt.Sprintf("%v", callbackData["orderInfo"]),
		"orderType":    fmt.Sprintf("%v", callbackData["orderType"]),
		"partnerCode":  fmt.Sprintf("%v", callbackData["partnerCode"]),
		"payType":      fmt.Sprintf("%v", callbackData["payType"]),
		"requestId":    fmt.Sprintf("%v", callbackData["requestId"]),
		"responseTime": fmt.Sprintf("%v", int64(callbackData["responseTime"].(float64))),
		"resultCode":   fmt.Sprintf("%d", int(callbackData["resultCode"].(float64))),
		"transId":      fmt.Sprintf("%v", int64(callbackData["transId"].(float64))),
	}
}

// Xây dựng chuỗi raw signature theo thứ tự a-z
func buildRawSignature(accessKey string, params map[string]string) string {
	var rawSignature bytes.Buffer

	// Thêm các tham số vào rawSignature theo thứ tự a-z
	rawSignature.WriteString("accessKey=")
	rawSignature.WriteString(accessKey)
	rawSignature.WriteString("&amount=")
	rawSignature.WriteString(params["amount"])
	rawSignature.WriteString("&extraData=")
	rawSignature.WriteString(params["extraData"])
	rawSignature.WriteString("&message=")
	rawSignature.WriteString(params["message"])
	rawSignature.WriteString("&orderId=")
	rawSignature.WriteString(params["orderId"])
	rawSignature.WriteString("&orderInfo=")
	rawSignature.WriteString(params["orderInfo"])
	rawSignature.WriteString("&orderType=")
	rawSignature.WriteString(params["orderType"])
	rawSignature.WriteString("&partnerCode=")
	rawSignature.WriteString(params["partnerCode"])
	rawSignature.WriteString("&payType=")
	rawSignature.WriteString(params["payType"])
	rawSignature.WriteString("&requestId=")
	rawSignature.WriteString(params["requestId"])
	rawSignature.WriteString("&responseTime=")
	rawSignature.WriteString(params["responseTime"])
	rawSignature.WriteString("&resultCode=")
	rawSignature.WriteString(params["resultCode"])
	rawSignature.WriteString("&transId=")
	rawSignature.WriteString(params["transId"])

	return rawSignature.String()
}

// Xây dựng chuỗi raw signature cho responce theo thứ tự a-z
func buildRawSignatureResponce(accessKey string, params map[string]string) string {
	var rawSignature bytes.Buffer

	// Thêm các tham số vào rawSignature theo thứ tự a-z
	rawSignature.WriteString("accessKey=")
	rawSignature.WriteString(accessKey)
	rawSignature.WriteString("&extraData=")
	rawSignature.WriteString(params["extraData"])
	rawSignature.WriteString("&message=")
	rawSignature.WriteString(params["message"])
	rawSignature.WriteString("&orderId=")
	rawSignature.WriteString(params["orderId"])
	rawSignature.WriteString("&orderInfo=")
	rawSignature.WriteString(params["partnerCode"])
	rawSignature.WriteString("&payType=")
	rawSignature.WriteString(params["requestId"])
	rawSignature.WriteString("&responseTime=")
	rawSignature.WriteString(params["responseTime"])
	rawSignature.WriteString("&resultCode=")
	rawSignature.WriteString(params["resultCode"])

	return rawSignature.String()
}

// Tính toán chữ ký HMAC_SHA256
func calculateSignature(rawSignature, secretKey string) string {
	hmac := hmac.New(sha256.New, []byte(secretKey))
	hmac.Write([]byte(rawSignature))
	return hex.EncodeToString(hmac.Sum(nil))
}
