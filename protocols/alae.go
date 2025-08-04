package protocols

import (
	"encoding/json"
	"fmt"
	"log"
	"ms-testing/models"
	"strings"
)

func getString(m map[string]interface{}, key string) string {
	val, ok := m[key]
	if ok {
		if key == "DevEUI" {
		return fmt.Sprintf("%v", strings.ToLower(val.(string)))
		} else {
			return fmt.Sprintf("%v", val)
		}
	}
	return ""
}

func ProcessAlaeMessage(msg []byte) []byte {
	var raw map[string]interface{}
	err := json.Unmarshal(msg, &raw)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		log.Printf("Invalid JSON: %s", string(msg))
		return nil
	}

	dataMap, ok := raw["data"].(map[string]interface{})
	if !ok {
		log.Println("Missing 'data' key in JSON")
		return nil
	}

	devUplink, ok := dataMap["DevEUI_uplink"].(map[string]interface{})
	if !ok {
		log.Println("DevEUI_uplink missing in data")
		return nil
	}

	payload, _ := devUplink["payload"].(map[string]interface{})
	gpsPrev, _ := payload["gpsPrevious"].(map[string]interface{})

	// Check if GPS data is present
	hasGPSData := payload["gpsLatitude"] != nil && payload["gpsLongitude"] != nil

	location := make(map[string]interface{})
	gps := make(map[string]interface{})

	if hasGPSData {
		location["lat"] = payload["gpsLatitude"]
		location["lng"] = payload["gpsLongitude"]
		location["accuracy"] = payload["horizontalAccuracy"]
		location["positionType"] = payload["rawPositionType"]

		gps = payload
	}

	signals := models.Signals{
		Battery:      payload["batteryLevel"],
		Temperature:  payload["temperatureMeasure"],
		ReceivedTime: getString(devUplink, "Time"),
		Location: map[string]interface{}{
			"lat":  gpsPrev["latitude"],
			"lng": gpsPrev["longitude"],
		},
	}

	final := models.FinalData{
		Data:         raw,
		HardwareID:   getString(devUplink, "DevEUI"),
		Timestamp:    getString(devUplink, "Time"),
		ServiceToken: getString(raw, "serviceToken"),
		Protocol:     getString(raw, "protocol"),
		GPS:          gps,
		Location:     location,
		Signals:      signals,
	}

	updatedJSON, err := json.Marshal(final)
	if err != nil {
		log.Println("Error marshaling processed data:", err)
		return nil
	}

	log.Println(string(updatedJSON))
	return updatedJSON
}
