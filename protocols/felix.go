package protocols

import (
	"encoding/json"
	"fmt"
	"log"
	"ms-testing/models"
	"strconv"
	"strings"
	"time"

	// "github.com/nats-io/nats.go"
)

func Parse_felix_data(msg []byte) []byte {
	fmt.Println("felix data received by the parser")

	var raw map[string]any
	if err := json.Unmarshal(msg, &raw); err != nil {
		log.Println("Error unmarshaling Felix data:", err)
		return nil
	}

	data, ok := raw["data"].(map[string]any)
	if !ok {
		log.Println("Error: data field missing or invalid")
		return nil
	}

	var hardwareID string
	var geo *models.GeoLocation
	var signals map[string]any
	var indoor *models.Indoor
	var positionTime int64

	var blue_wifi []map[string]any
	var receivedTime int64

	locationSet := new(bool)
	batterySet := new(bool)
	temperatureSet := new(bool)
	lightSet := new(bool)
	bluetoothSet := new(bool)

	signals = make(map[string]any)

	// Process BLE-type payload done by cisco spaces
	if notifications, ok := data["notifications"].([]any); ok && len(notifications) > 0 {
		if notification, ok := notifications[0].(map[string]any); ok {
			if geoMap, ok := notification["geoCoordinate"].(map[string]any); ok {
				lat, _ := geoMap["latitude"].(float64)
				lng, _ := geoMap["longitude"].(float64)
				unit, _ := geoMap["unit"].(string)
				geo = &models.GeoLocation{
					Lat:  lat,
					Lng:  lng,
					Unit: unit,
				}
			}

			hardwareID, _ = notification["deviceId"].(string)

			if tsStr, ok := notification["lastLocated"].(string); ok && tsStr != "" {
				layout := "2006-01-02 15:04:05.000-0700"
				if t, err := time.Parse(layout, tsStr); err == nil {
					positionTime = t.UnixNano() / int64(time.Millisecond)
				}
			}

			var locationHierarchy map[string]any
			if lhRaw, ok := notification["locationHierarchy"]; ok {
				if lhMap, ok := lhRaw.(map[string]any); ok {
					locationHierarchy = lhMap
				}
			}

			signals["deviceName"] = hardwareID

			indoor = &models.Indoor{
				Building:        "East Alpine Road 230",
				BuildingId:      11878987,
				FloorIndex:      0,
				FloorLabel:      "0",
				BuildingModelId: 159643361,
				LocationHierarchy: locationHierarchy,
			}
		}
	}

	if edids, ok := data["end_device_ids"].(map[string]any); ok {
		if devEUI, ok := edids["dev_eui"].(string); ok && devEUI != "" {
			hardwareID = devEUI // override if present
		}
		if devType, ok := edids["device_id"].(string); ok{
			fmt.Println("device type:", devType)
			signals["deviceName"] = devType
		}
	}

	// Extract uplink_message whether inside end_device_ids or directly inside data
	uplinkMsg := extractUplinkMessage(data)


	if uplinkMsg != nil {
		processUplinkMessage(uplinkMsg, &geo, &signals, &blue_wifi, &positionTime, locationSet, batterySet, temperatureSet, lightSet, bluetoothSet)
	}

	// Extract from normalized_payload if present (common case)
	processNormalizedPayload(uplinkMsg, &signals, batterySet, temperatureSet)

	if receivedAt, ok := data["received_at"].(string); ok {
		if parsedTime, err := time.Parse(time.RFC3339Nano, receivedAt); err == nil {
			millis := parsedTime.UnixNano() / int64(time.Millisecond)
			receivedTime = millis
		}
		signals["received_time"] = receivedTime
	}

	// If location missing, fallback to RadioData
	if geo == nil {
		log.Println("GeoLocation missing, preparing radiodata...")
		
		var radioData models.RadioData

		if !*bluetoothSet {
			radioData = models.RadioData{
				WifiAccessPoints: blue_wifi,
			}
		} else {
			radioData = models.RadioData{
				BluetoothBeacons: blue_wifi,
			}
		}
		
		incoming := models.IncomingData{
			HardwareID:   hardwareID,
			Time:         raw["time"],
			MessageID:    fmt.Sprintf("%v", raw["messageId"]),
			Protocol:     fmt.Sprintf("%v", raw["protocol"]),
			ServiceToken: fmt.Sprintf("%v", raw["serviceToken"]),
			Data:         msg,
			RadioData:    &radioData,
			Signals:      signals,
			PositionTime: positionTime,
		}

		result, err := json.Marshal(incoming)
		if err != nil {
			log.Println("Error marshaling fallback radiodata:", err)
			return nil
		}
		return result
	}
	
	incoming := models.IncomingData{
		HardwareID:   hardwareID,
		Time:         raw["time"],
		MessageID:    fmt.Sprintf("%v", raw["messageId"]),
		Protocol:     fmt.Sprintf("%v", raw["protocol"]),
		ServiceToken: fmt.Sprintf("%v", raw["serviceToken"]),
		GeoLocation:  geo,
		Data:         msg,
		Indoor:       indoor,
		Signals:      signals,
		PositionTime: positionTime,
	}

	result, err := json.Marshal(incoming)
	if err != nil {
		log.Println("Error marshaling final incoming data:", err)
		return nil
	}
	return result
}

func extractUplinkMessage(data map[string]any) map[string]any {
	if uplinkMsg, ok := data["uplink_message"].(map[string]any); ok {
		return uplinkMsg
	} else if edids, ok := data["end_device_ids"].(map[string]any); ok {
		if uplinkMsg, ok := edids["uplink_message"].(map[string]any); ok {
			return uplinkMsg
		}
	} else if ulinknrm, ok := data["uplink_normalized"].(map[string]any); ok {
		return ulinknrm
	}
	return nil
}

func processUplinkMessage(
	uplinkMsg map[string]any,
	geo **models.GeoLocation,
	signals *map[string]any,
	blue_wifi *[]map[string]any,
	positionTime *int64, locationSet *bool, batterySet *bool, temperatureSet *bool, lightSet *bool, bluetoothSet *bool,
) {

	// Step 1: Check locations["user"]
	if locs, ok := uplinkMsg["locations"].(map[string]any); ok {
		if userLoc, ok := locs["user"].(map[string]any); ok {
			lat, _ := userLoc["latitude"].(float64)
			lng, _ := userLoc["longitude"].(float64)
			source, _ := userLoc["source"].(string)
			alt, _ := userLoc["altitude"].(float64)
			*geo = &models.GeoLocation{
				Lat:    lat,
				Lng:    lng,
				Source: source,
				Alt:    alt,
			}
			*locationSet = true
		}
	}

	// Step 2: Process decoded_payload
	if decoded, ok := uplinkMsg["decoded_payload"].(map[string]any); ok {
		processDecodedPayload(decoded, geo, signals, blue_wifi, locationSet, batterySet, temperatureSet, lightSet, bluetoothSet)
	}

	// Step 3: If both locationSet == false AND bluetoothBeacons empty, fallback to rx_metadata
	if !*locationSet && (blue_wifi == nil || len(*blue_wifi) == 0) {
		if rxMeta, ok := uplinkMsg["rx_metadata"].([]any); ok && len(rxMeta) > 0 {
			if firstMeta, ok := rxMeta[0].(map[string]any); ok {
				if loc, ok := firstMeta["location"].(map[string]any); ok {
					lat, _ := loc["latitude"].(float64)
					lng, _ := loc["longitude"].(float64)
					alt, _ := loc["altitude"].(float64)
					source, _ := loc["source"].(string)
					*geo = &models.GeoLocation{
						Lat:    lat,
						Lng:    lng,
						Alt:    alt,
						Source: source,
					}
				}
			}
		}
	}

	// Step 4: Process positionTime as before
	if rxMeta, ok := uplinkMsg["rx_metadata"].([]any); ok && len(rxMeta) > 0 {
		if firstMeta, ok := rxMeta[0].(map[string]any); ok {

			// Step 4.1: Check time field
			if tsStr, ok := firstMeta["time"].(string); ok && tsStr != "" {
				if t, err := time.Parse(time.RFC3339Nano, tsStr); err == nil {
					*positionTime = t.UnixNano() / int64(time.Millisecond)
				}
			}else if rec_at, ok := firstMeta["received_at"].(string); ok && rec_at != "" {
				if t, err := time.Parse(time.RFC3339Nano, rec_at); err == nil {
					*positionTime = t.UnixNano() / int64(time.Millisecond)
				}
			}
		}
	}	
}

func processDecodedPayload(
	decoded map[string]any,
	geo **models.GeoLocation,
	signals *map[string]any,
	blue_wifi *[]map[string]any,
	locationSet, batterySet, temperatureSet, lightSet, bluetoothSet *bool,
) {
	// Basic key-value parsing
	for k, val := range decoded {
		switch k {
		case "battery_percent":
			switch v := val.(type) {
			case string:
				cleaned := strings.TrimSpace(strings.TrimSuffix(v, "%"))
				parsed, err := strconv.ParseFloat(cleaned, 64)
				if err != nil {
					(*signals)["batteryLevel"] = 0
				} else {
					(*signals)["batteryLevel"] = parsed
					*batterySet = true
				}
			default:
				(*signals)["batteryLevel"] = val
				*batterySet = true
			}

		case "messages", "position_data", "light_intensity", "ambient_temperature", "humidity", "relative_humidity", "illumination":
			// handled below, skip here
			continue
		
		case "temperature":
			if tempMap, ok := val.(map[string]any); ok {
				(*signals)["temperature"] = tempMap
				continue
			}
		case "battery":
			(*signals)["batteryLevel"] = val
			*batterySet = true
			continue

		default:
			(*signals)[k] = val
		}
	}

	// Parse messages
	if msgs, ok := decoded["messages"].([]any); ok && len(msgs) > 0 {
		if group, ok := msgs[0].([]any); ok {
			for _, m := range group {
				entry := m.(map[string]any)
				switch entry["measurementId"] {
				case "5002":
					if values, ok := entry["measurementValue"].([]any); ok {
						for _, b := range values {
							beacon := b.(map[string]any)
							macRaw := fmt.Sprintf("%v", beacon["mac"])
							rssiRaw := fmt.Sprintf("%v", beacon["rssi"])
							formattedMac := formatMac(macRaw)
							cleanedRSSI := cleanRSSI(rssiRaw)
							*blue_wifi = append(*blue_wifi, map[string]any{
								"macAddress":     formattedMac,
								"signalStrength": cleanedRSSI,
							})
						}
						*bluetoothSet = true
					}
				case "5001":
					if values, ok := entry["measurementValue"].([]any); ok {
						for _, b := range values {
							beacon := b.(map[string]any)
							macRaw := fmt.Sprintf("%v", beacon["mac"])
							rssiRaw := fmt.Sprintf("%v", beacon["rssi"])
							formattedMac := formatMac(macRaw)
							cleanedRSSI := cleanRSSI(rssiRaw)
							*blue_wifi = append(*blue_wifi, map[string]any{
								"macAddress":     formattedMac,
								"signalStrength": cleanedRSSI,
							})
						}
					}

				case "4097":
					if !*temperatureSet {
						(*signals)["temperatureLevel"] = entry["measurementValue"]
						*temperatureSet = true
					}

				case "3000":
					if !*batterySet {
						(*signals)["batteryLevel"] = entry["measurementValue"]
						*batterySet = true
					}

				case "4199":
					if !*lightSet {
						(*signals)["light"] = entry["measurementValue"]
						*lightSet = true
					}

				case "4200":
					(*signals)["eventStatus"] = entry["measurementValue"]
				}
			}
		}else if _, ok := msgs[0].(map[string]any); ok {
			// First group all entries by type
			grouped := map[string][]map[string]any{}
			for _, m := range msgs {
				entry := m.(map[string]any)
				typ, _ := entry["type"].(string)
				grouped[typ] = append(grouped[typ], entry)
			}

			// Now, build signals
			for typ, entries := range grouped {
				if typ == "upload_battery" {
					if battery, ok := entries[0]["battery"]; ok {
						(*signals)["batteryLevel"] = battery
					}
					continue
				}
				if len(entries) == 1 {
					entry := entries[0]
					if val, ok := entry["measurementValue"]; ok {
						(*signals)[typ] = val
					} else if val, ok := entry["interval"]; ok {
						(*signals)[typ] = val
					}
				} else {
					// Repeat:
					for _, entry := range entries {
						if measID, mok := entry["measurementId"]; mok {
							key := fmt.Sprintf("%s_%v", typ, measID)
							if val, ok := entry["measurementValue"]; ok {
								(*signals)[key] = val
							}
						}
					}
				}
			}
		}
	}

	// Position data parsing
	if values, ok := decoded["position_data"].([]any); ok {
		for _, b := range values {
			entry := b.(map[string]any)

			if latStr, ok := entry["latitude"].(string); ok {
				if lngStr, ok2 := entry["longitude"].(string); ok2 {
					lat := parseCoordinate(latStr)
					lng := parseCoordinate(lngStr)

					*geo = &models.GeoLocation{
						Lat:    lat,
						Lng:    lng,
						Alt:    0,
						Source: "GPS",
					}
					*locationSet = true
					continue
				}
			}else if macRaw, ok := entry["mac_address"].(string); ok {
				rssiRaw := fmt.Sprintf("%v", entry["rssi"])
				formattedMac := formatMac(macRaw)
				cleanedRSSI := cleanRSSI(rssiRaw)

				*blue_wifi = append(*blue_wifi, map[string]any{
					"macAddress":     formattedMac,
					"signalStrength": cleanedRSSI,
				})
				*bluetoothSet = true
			}
		}
	}

	// Fallback for AirTemperature
	if !*temperatureSet {
		if val, ok := decoded["temperature"]; ok {
			switch v := val.(type) {
			case map[string]any:
				if celsius, ok := v["celsius"].(float64); ok {
					(*signals)["temperatureLevel"] = celsius
					*temperatureSet = true
				}
				if fahrenheit, ok := v["fahrenheit"].(float64); ok {
					(*signals)["temperatureLevelF"] = fahrenheit
					*temperatureSet = true
				}
			case string:
				cleaned := strings.TrimSpace(strings.TrimSuffix(v, "°C"))
				parsed, err := strconv.ParseFloat(cleaned, 64)
				if err != nil {
					(*signals)["temperatureLevel"] = 0
				} else {
					(*signals)["temperatureLevel"] = parsed
				}
				*temperatureSet = true
			default:
				(*signals)["temperatureLevel"] = v
				*temperatureSet = true
			} 
		}else if val, ok := decoded["ambient_temperature"]; ok{
			(*signals)["temperatureLevel"] = val
			*temperatureSet = true
		}
	}

	// Fallback for Light
	if !*lightSet {
		if val, ok := decoded["light_intensity"]; ok {
			(*signals)["light"] = val
			*lightSet = true
		} else if val, ok := decoded["illumination"]; ok {
			(*signals)["light"] = val
			*lightSet = true
		}
	}

	// Humidity
	if val, ok := decoded["humidity"]; ok {
		(*signals)["humidity"] = val
	}else if val, ok := decoded["relative_humidity"]; ok{
		(*signals)["humidity"] = val
	}
}


	

func processNormalizedPayload(uplink map[string]any, signals *map[string]any, batterySet *bool, temperatureSet *bool) {
	// Case 1: normalized_payload is a slice
	if normalizedArray, ok := uplink["normalized_payload"].([]any); ok && len(normalizedArray) > 0 {
		if firstEntry, ok := normalizedArray[0].(map[string]any); ok {
			if !*batterySet{
				if val, ok := firstEntry["battery"]; ok {
				(*signals)["batteryLevel"] = val
			}
			}
			
			if airData, ok := firstEntry["air"].(map[string]any); ok {
				if !*temperatureSet{
					if val, ok := airData["temperature"]; ok {
					(*signals)["temperatureLevel"] = val
				}
				}
				
				if val, ok := airData["relativeHumidity"]; ok {
					(*signals)["humidity"] = val
				}
			}
		}
	}else if normalizedMap, ok := uplink["normalized_payload"].(map[string]any); ok {
		if !*batterySet{
			if val, ok := normalizedMap["battery"]; ok {
				(*signals)["batteryLevel"] = val
			}
		}
		
		if airData, ok := normalizedMap["air"].(map[string]any); ok {
			if !*temperatureSet{
				if val, ok := airData["temperature"]; ok {
					(*signals)["temperatureLevel"] = val
				}
			}
			if val, ok := airData["relativeHumidity"]; ok {
				(*signals)["humidity"] = val
			}
		}
	}
}


func formatMac(mac string) string {
    // Remove both ":" and "-" from input
    mac = strings.ReplaceAll(mac, ":", "")
    mac = strings.ReplaceAll(mac, "-", "")
    if len(mac) != 12 {
        return mac
    }
    return fmt.Sprintf("%s:%s:%s:%s:%s:%s",
        mac[0:2], mac[2:4], mac[4:6],
        mac[6:8], mac[8:10], mac[10:12])
}

func cleanRSSI(rssiRaw string) string {
	return strings.TrimSuffix(rssiRaw, "dBm")
}

func parseCoordinate(coord string) float64 {
	cleaned := strings.TrimSuffix(coord, "°")
	val, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0
	}
	return val
}
