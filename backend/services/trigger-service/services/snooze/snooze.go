package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	models "github.com/dath-241/coin-price-be-go/services/trigger-service/models/alert"
	noify "github.com/dath-241/coin-price-be-go/services/trigger-service/services"
	services "github.com/dath-241/coin-price-be-go/services/trigger-service/services/alert"
	"github.com/dath-241/coin-price-be-go/services/trigger-service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"time"
)

func CheckPriceCondition(alert *models.Alert) bool {

	var Price float64
	var err error
	if alert.Type == "spot" {
		Price, err = services.GetSpotPrice(alert.Symbol)
	} else if alert.Type == "future" {
		Price, err = services.GetFuturePrice(alert.Symbol)
	} else if alert.Type == "funding_rate" {
		Price, err = services.GetFundingRate(alert.Symbol)
	}
	if err != nil {
		log.Printf("Error fetching price: %v", err)
		return false
	}
	if alert.Condition == "==" {
		if alert.Threshold == Price {
			return true
		}
	} else if alert.Condition == ">=" {
		if alert.Threshold <= Price {
			return true
		}
	} else if alert.Condition == "<=" {
		if alert.Threshold >= Price {
			return true
		}
	}
	return false
}

func checkRepeatCount(alert *models.Alert) bool {
	if alert.MaxRepeatCount > 0 && alert.RepeatCount >= alert.MaxRepeatCount {
		return false
	}
	return true
}

func CheckNewListingAndDelisting(alert *models.Alert) bool {
	newSymbols, delistedSymbols, err := services.FetchSymbolsFromBinance()
	if err != nil {
		log.Printf("Error fetching symbol data: %v", err)
		return false
	}
	if alert.Type == "new_listing" {
		for _, symbol := range newSymbols {
			if symbol == alert.Symbol {
				return true
			}
		}
	} else if alert.Type == "delisting" {
		for _, symbol := range delistedSymbols {
			if symbol == alert.Symbol {
				return true
			}
		}
	}
	return false
}

func FetchAlerts(alertID string) ([]models.Alert, error) {
	var results []models.Alert

	filter := bson.M{}
	if alertID != "" {
		objectID, err := primitive.ObjectIDFromHex(alertID)
		if err != nil {
			return nil, fmt.Errorf("invalid ID format")
		}
		filter["_id"] = objectID
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := utils.AlertCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func CheckNumberOfAlertSent(alert *models.Alert) bool {
	if alert.MaxRepeatCount > 0 && alert.RepeatCount >= alert.MaxRepeatCount {
		return false
	}
	return true
}
func CheckSnoozeCondition(alert *models.Alert) bool {
	currentTime := time.Now()

	switch alert.SnoozeCondition {
	case "Only once":
		if alert.RepeatCount > 0 {
			alert.IsActive = false
			log.Println("Đã hết số lần gửi", alert.ID.Hex())
			return false
		}
	case "Once a day":
		if currentTime.Sub(alert.UpdatedAt.Time()) < 24*time.Hour {
			log.Println("Đã gửi trong 1 ngày", alert.ID.Hex())
			return false
		}
	case "Once per 10 seconds":
		if currentTime.Sub(alert.UpdatedAt.Time()) < 10*time.Second {
			// log.Println("Chưa đủ 10 giây", alert.ID.Hex())
			return false
		}
	case "Once per 5 minutes":
		if currentTime.Sub(alert.UpdatedAt.Time()) < 5*time.Minute {
			log.Println("Chưa đủ 5 phut", alert.ID.Hex())
			return false
		}
	case "At Specific Time":
		start := alert.NextTriggerTime	
		if currentTime.Before(start) || currentTime.After(start) {
			if currentTime.Before(start) {
				log.Println("Chưa đến thời gian gửi", alert.ID.Hex())
			} else if currentTime.After(start) {
				log.Println("Đã hết thời gian gửi", alert.ID.Hex())
			}

			return false
		}
	case "Forever":
		return true
	}

	return true
}

func UpdateAlertAfterTrigger(alert *models.Alert) {
	currentTime := time.Now()

	alert.RepeatCount++
	if alert.MaxRepeatCount > 0 && alert.RepeatCount >= alert.MaxRepeatCount {
		alert.IsActive = false
	}
	switch alert.SnoozeCondition {
	case "Only once":
		alert.IsActive = false
	case "Once a day":
		alert.NextTriggerTime = currentTime.Add(24 * time.Hour)
	case "Once per 10 seconds":
		alert.NextTriggerTime = currentTime.Add(10 * time.Second)
	case "Once per 5 minutes":
		alert.NextTriggerTime = currentTime.Add(5 * time.Minute)
	case "At Specific Time":
		alert.IsActive = false
	case "Forever":
	}

	SaveAlert(alert)
}

func CheckAndSendAlerts() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"is_active": true}
	cursor, err := utils.AlertCollection.Find(ctx, filter)
	if err != nil {
		log.Println("Failed to fetch alerts:", err)
		return
	}
	defer cursor.Close(ctx)

	var alerts []models.Alert
	if err = cursor.All(ctx, &alerts); err != nil {
		return
	}

	for _, alert := range alerts {
		conditionMet := false
		if (alert.Type == "spot" || alert.Type == "future" || alert.Type == "funding_rate") && CheckPriceCondition(&alert) && checkRepeatCount(&alert) {
			conditionMet = true
		} else if (alert.Type == "new_listing" || alert.Type == "delisting") && CheckNewListingAndDelisting(&alert) && checkRepeatCount(&alert) {
			conditionMet = true
		}

		if conditionMet {
			if CheckSnoozeCondition(&alert) {
				noify.NotifyUserTriggers(alert.UserID)
				log.Println("Đã gửi thông báo!!!:", alert.ID.Hex(), alert.Type, alert.Symbol)
				UpdateAlertAfterTrigger(&alert)

				if !alert.IsActive {
					log.Println("Cảnh báo đã hết hiệu lực:", alert.ID.Hex())
					break
				}
			} else {
				// log.Println("Không đủ điều kiện snooze để gửi cảnh báo:", alert.ID.Hex())
			}
		} else {
			log.Println("Cảnh báo không đủ điều kiện kích hoạt:", alert.ID.Hex(), alert.Type, alert.Symbol)
		}
	}
}

// SaveAlert lưu hoặc cập nhật một cảnh báo trong cơ sở dữ liệu
func SaveAlert(alert *models.Alert) error {
	if alert.ID.IsZero() {
		alert.ID = primitive.NewObjectID()
		alert.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	} else {
		alert.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": alert.ID}
	update := bson.M{"$set": alert}
	result, err := utils.AlertCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("failed to save or update alert: " + err.Error())
	}

	if result.MatchedCount == 0 {
		_, err = utils.AlertCollection.InsertOne(ctx, alert)
		if err != nil {
			return errors.New("failed to insert new alert: " + err.Error())
		}
	}

	return nil
}
