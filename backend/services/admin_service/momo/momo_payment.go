package momo

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/sony/sonyflake"
)

var (
	partnerCodeEnv string
	accessKeyEnv   string
	secretKeyEnv   string
	baseURLEnv     string
	redirectURLEnv string
	ipnURLEnv      string
)

func Init() {
	partnerCodeEnv 	= getEnv("MOMO_PARTNER_CODE")
	accessKeyEnv 	= getEnv("MOMO_ACCESS_KEY")
	secretKeyEnv 	= getEnv("MOMO_SECRET_KEY")
	baseURLEnv 		= getEnv("MOMO_BASE_URL")
	redirectURLEnv 	= getEnv("MOMO_REDIRECT_URL")
	ipnURLEnv 		= getEnv("MOMO_IPN_URL")
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

// Payload defines the data structure for a MoMo payment request
type payload struct {
	PartnerCode  string `json:"partnerCode"`
	AccessKey    string `json:"accessKey"`
	RequestID    string `json:"requestId"`
	Amount       string `json:"amount"`
	OrderID      string `json:"orderId"`
	OrderInfo    string `json:"orderInfo"`
	PartnerName  string `json:"partnerName"`
	StoreId      string `json:"storeId"`
	OrderGroupId string `json:"orderGroupId"`
	Lang         string `json:"lang"`
	AutoCapture  bool   `json:"autoCapture"`
	RedirectUrl  string `json:"redirectUrl"`
	IpnUrl       string `json:"ipnUrl"`
	ExtraData    string `json:"extraData"`
	RequestType  string `json:"requestType"`
	Signature    string `json:"signature"`
}

// payment handles the MoMo payment process
func CreateMoMoPayment(amount string, vipLevel string, orderInfo string) (string, string, error) {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})

	// Generate random orderID and requestID
	a, _ := flake.NextID()
	b, _ := flake.NextID()

	// URL cho MoMo API để kiểm tra trạng thái thanh toán
	url := baseURLEnv + "/v2/gateway/api/create"

	var endpoint = url
	var accessKey = accessKeyEnv
	var secretKey = secretKeyEnv
	var partnerCode = partnerCodeEnv
	var redirectUrl = redirectURLEnv
	var ipnUrl = ipnURLEnv
	var orderId = strconv.FormatUint(a, 16)
	var requestId = strconv.FormatUint(b, 16)
	var extraData = ""
	var partnerName = "MoMo Payment"
	var storeId = "Test Store"
	var orderGroupId = ""
	var autoCapture = true
	var lang = "vi"
	var requestType = "payWithMethod"
	
	// Check if any required environment variable is missing 
	if (endpoint == "" || 
		accessKey == "" || 
		secretKey == "" || 
		partnerCode == "" || 
		redirectUrl == "" || 
		ipnUrl == "") {
		log.Println("Missing required MoMo environment variables") 
		return "", "", nil
	}

	// Build the raw signature
	var rawSignature bytes.Buffer
	rawSignature.WriteString("accessKey=")
	rawSignature.WriteString(accessKey)
	rawSignature.WriteString("&amount=")
	rawSignature.WriteString(amount)
	rawSignature.WriteString("&extraData=")
	rawSignature.WriteString(extraData)
	rawSignature.WriteString("&ipnUrl=")
	rawSignature.WriteString(ipnUrl)
	rawSignature.WriteString("&orderId=")
	rawSignature.WriteString(orderId)
	rawSignature.WriteString("&orderInfo=")
	rawSignature.WriteString(orderInfo)
	rawSignature.WriteString("&partnerCode=")
	rawSignature.WriteString(partnerCode)
	rawSignature.WriteString("&redirectUrl=")
	rawSignature.WriteString(redirectUrl)
	rawSignature.WriteString("&requestId=")
	rawSignature.WriteString(requestId)
	rawSignature.WriteString("&requestType=")
	rawSignature.WriteString(requestType)

	// Generate HMAC SHA256 signature
	hmac := hmac.New(sha256.New, []byte(secretKey))
	hmac.Write(rawSignature.Bytes())
	signature := hex.EncodeToString(hmac.Sum(nil))

	// Create payload for the request
	payload := payload{
		PartnerCode:  partnerCode,
		AccessKey:    accessKey,
		RequestID:    requestId,
		Amount:       amount,
		RequestType:  requestType,
		RedirectUrl:  redirectUrl,
		IpnUrl:       ipnUrl,
		OrderID:      orderId,
		StoreId:      storeId,
		PartnerName:  partnerName,
		OrderGroupId: orderGroupId,
		AutoCapture:  autoCapture,
		Lang:         lang,
		OrderInfo:    orderInfo,
		ExtraData:    extraData,
		Signature:    signature,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshaling payload:", err)
		return "", "", err
	}

	// Send HTTP request to MoMo endpoint
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatalln("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Parse the response from MoMo
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Println("Error decoding response:", err)
		return "", "", err
	}

	// Check if MoMo response contains a valid payment URL
	if result["payUrl"] == nil {
		log.Println("MoMo payment URL not found or invalid response")
		return "", "", nil
	}

	return result["payUrl"].(string), orderId, nil
}