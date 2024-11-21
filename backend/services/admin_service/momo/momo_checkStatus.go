package momo

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// QueryPaymentStatusRequest defines the request structure for checking payment status
type QueryPaymentStatusRequest struct {
	PartnerCode string `json:"partnerCode"`
	RequestID   string `json:"requestId"`
	OrderID     string `json:"orderId"`
	Signature   string `json:"signature"`
	Lang        string `json:"lang"`
}

// QueryPaymentStatusResponse defines the response structure from MoMo API
type QueryPaymentStatusResponse struct {
	PartnerCode  string                 `json:"partnerCode"`
	RequestID    string                 `json:"requestId"`
	OrderID      string                 `json:"orderId"`
	Amount       float64                `json:"amount"`
	ResultCode   int                    `json:"resultCode"`
	Message      string                 `json:"message"`
	PaymentType  string                 `json:"payType"`
	PaymentInfo  map[string]interface{} `json:"promotionInfo"`
	ResponseTime int64                  `json:"responseTime"`
}

// GenerateSignature generates the HMAC_SHA256 signature for MoMo payment query
func GenerateQuerySignature(orderId, requestId string) string {
	data := fmt.Sprintf("accessKey=%s&orderId=%s&partnerCode=%s&requestId=%s", accessKeyEnv, orderId, partnerCodeEnv, requestId)
	h := hmac.New(sha256.New, []byte(secretKeyEnv))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// QueryPaymentStatus queries the payment status using MoMo's API
func QueryPaymentStatus(orderId, requestId, lang string) (QueryPaymentStatusResponse, error) {
	// Generate signature
	signature := GenerateQuerySignature(orderId, requestId)

	// Prepare the request payload
	req := QueryPaymentStatusRequest{
		PartnerCode: partnerCodeEnv,
		RequestID:   requestId,
		OrderID:     orderId,
		Signature:   signature,
		Lang:        lang,
	}

	// Marshal the request body
	jsonRequest, err := json.Marshal(req)
	if err != nil {
		log.Println("Error marshalling request:", err)
		return QueryPaymentStatusResponse{}, err
	}

	// URL cho MoMo API để kiểm tra trạng thái thanh toán
	url := fmt.Sprintf("%s/v2/gateway/api/query", baseURLEnv)

	// Make the HTTP request to MoMo API
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonRequest))
	if err != nil {
		log.Println("Error sending request:", err)
		return QueryPaymentStatusResponse{}, err
	}
	defer resp.Body.Close()

	// Parse the response
	var response QueryPaymentStatusResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("Error decoding response:", err)
		return QueryPaymentStatusResponse{}, err
	}

	// // Check if the result code is not successful
	// if response.ResultCode != 0 {
	// 	return response, fmt.Errorf("failed to query payment status: %s", response.Message)
	// }

	// Return the successful response
	return response, nil
}
