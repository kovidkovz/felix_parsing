package protocols

import (
	"encoding/json"
	"fmt"
	"log"
	"ms-testing/models"
	"strconv"
	"strings"
	"time"
)

func ProcessAlaeMessage(msg []byte) []byte {
	fmt.Println("abeeway data received by the parser")
	var raw map[string]any
	if err := json.Unmarshal(msg, &raw); err != nil {
		log.Println("Error unmarshaling Abeeway data:", err)
		return nil
	}

	data, ok := raw["data"].(map[string]any)
	if !ok {
		log.Println("Error: data field missing or invalid")
		return nil
	}

	var hardwareID string
	var geo *models.GeoLocation
	var signals = make(map[string]any)
	var radioData models.RadioData
	var positionTime int64
	var protocol string

	// Case A: DevEUI_uplink (semtech/raw lora type)
	if uplink, found := data["DevEUI_uplink"].(map[string]any); found {
		// Hardware
		if devEUI, ok := uplink["DevEUI"].(string); ok {
			hardwareID = devEUI
			signals["deviceName"] = devEUI
		}
		// Device Name
		if customerData, ok := uplink["CustomerData"].(map[string]any); ok {
			if name, ok := customerData["name"].(string); ok {
				signals["deviceName"] = name
			}
		}

		// Battery (Kabhi aata hai toh payload se nikalna pad sakta hai, yahan skip kar rahe)
		// Gateway radio params
		if v, ok := uplink["LrrRSSI"].(float64); ok {
			signals["gatewayRssi"] = v
		}
		if v, ok := uplink["LrrSNR"].(float64); ok {
			signals["gatewaySnr"] = v
		}
		if v, ok := uplink["MeanPER"].(float64); ok {
			signals["meanPER"] = v
		}

		// Location
		lat, _latOK := uplink["LrrLAT"].(float64)
		lng, _lngOK := uplink["LrrLON"].(float64)
		if _latOK && _lngOK {
			geo = &models.GeoLocation{
				Lat: lat,
				Lng: lng,
				// Source: "abeeway",  // Set if you want
			}
			positionTime = extractTimeField(uplink["Time"])
		}

		// Radio: Lrrs
		if lrrs, ok := uplink["Lrrs"].(map[string]any); ok {
			if lrrarr, ok := lrrs["Lrr"].([]any); ok && len(lrrarr) > 0 {
				var wifiAps []map[string]any
				for _, entry := range lrrarr {
					lrr, _ := entry.(map[string]any)
					if lrr == nil {
						continue
					}
					ap := map[string]any{}
					ap["stationId"] = lrr["Lrrid"]
					if v, ok := lrr["LrrRSSI"]; ok {
						ap["signalStrength"] = v
					}
					if v, ok := lrr["LrrSNR"]; ok {
						ap["snr"] = v
					}
					wifiAps = append(wifiAps, ap)
				}
				radioData.WifiAccessPoints = wifiAps
			}
		}
		protocol = "lorawan"
	}

	// Case B: location+signals (NAKED GPS/CLOUD API type)
	if loc, found := data["location"].(map[string]any); found {
		geo = &models.GeoLocation{
			Lat: getFloat(loc["lat"]),
			Lng: getFloat(loc["lng"]),
		}
		if pt, ok := data["positionTime"].(float64); ok {
			positionTime = int64(pt)
		}
		// signals:
		if sigs, ok := data["signals"].(map[string]any); ok {
			for k, v := range sigs {
				// Copy all in signals, format keys for battery etc.
				if k == "battery" {
					signals["batteryLevel"] = v
				} else if strings.Contains(strings.ToLower(k), "temperature") {
					signals["temperatureLevel"] = v
				} else {
					signals[k] = v
				}
			}
		}
		// Radio data
		if wifis, ok := data["wifiAccessPoints"].([]any); ok && len(wifis) > 0 {
			// Proper wifi AP format
			var arr []map[string]any
			for _, ap := range wifis {
				if apmap, ok := ap.(map[string]any); ok {
					arr = append(arr, apmap)
				}
			}
			if len(arr) > 0 {
				radioData.WifiAccessPoints = arr
			}
		}
		protocol = "cloud-gps"
	}

	// Case C: PURE signals, maybe "wifiAccessPoints" without location
	if sigs, ok := data["signals"].(map[string]any); ok && geo == nil {
		for k, v := range sigs {
			if k == "battery" {
				signals["batteryLevel"] = v
			} else if strings.Contains(strings.ToLower(k), "temperature") {
				signals["temperatureLevel"] = v
			} else {
				signals[k] = v
			}
		}
	}
	if ap, ok := data["wifiAccessPoints"].([]any); ok && radioData.WifiAccessPoints == nil {
		var arr []map[string]any
		for _, x := range ap {
			if m, ok := x.(map[string]any); ok {
				arr = append(arr, m)
			}
		}
		if len(arr) > 0 {
			radioData.WifiAccessPoints = arr
		}
	}

	// HardwareId fallback
	if hardwareID == "" {
		// Try request/id or data/id
		if req, ok := data["request"].(map[string]any); ok {
			if id, ok := req["id"].(string); ok {
				hardwareID = id
			}
		} else if id, ok := data["id"].(string); ok {
			hardwareID = id
		} else if rawhid, ok := raw["hardwareId"].(string); ok && rawhid != "" {
			hardwareID = rawhid
		}
	}
	// Lowercase
	hardwareID = strings.ToUpper(hardwareID)

	// Final assemble
	incoming := models.IncomingData{
		HardwareID:   hardwareID,
		Time:         raw["time"],
		MessageID:    fmt.Sprintf("%v", raw["messageId"]),
		Protocol:     protocol,
		ServiceToken: fmt.Sprintf("%v", raw["serviceToken"]),
		Data:         msg,
		GeoLocation:  geo,
		Signals:      signals,
		PositionTime: positionTime,
	}
	if len(radioData.WifiAccessPoints) > 0 {
		incoming.RadioData = &radioData
	}

	result, err := json.Marshal(incoming)
	if err != nil {
		log.Println("Error marshaling Abeeway result:", err)
		return nil
	}
	return result
}

// ---------- Helpers ----------

func getFloat(v interface{}) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	case int64:
		return float64(n)
	case string:
		f, _ := strconv.ParseFloat(n, 64)
		return f
	}
	return 0
}

func extractTimeField(x interface{}) int64 {
	// from "2025-08-11T10:01:54.277+00:00"
	if str, ok := x.(string); ok && str != "" {
		t, err := time.Parse(time.RFC3339Nano, str)
		if err == nil {
			return t.UnixNano() / int64(time.Millisecond)
		}
	}
	return 0
}
