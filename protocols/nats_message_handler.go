package protocols

import (
	// "context"
	"encoding/json"
	"fmt"
	"log"
	// "net/http"
	// "net/url"
	// "strings"
	// "time"

	// "github.com/kovidkovz/natstemplate/natstemplate"
	// "gitlab.combain.com/go/keymate/cpslookup"
	// "gitlab.combain.com/go/keymate/logs"
	// "github.com/nats-io/nats.go"
	// "gitlab.combain.com/traxmate/natstemplate/natstemplate"
)

// func runMapping(src interface{}, dest interface{}) {
// 	bytes, err := json.Marshal(src)
// 	if err != nil {
// 		fmt.Println("Marshal error:", err)
// 		return
// 	}
// 	err = json.Unmarshal(bytes, dest)
// 	if err != nil {
// 		fmt.Println("Unmarshal error:", err)
// 	}
// }

type ParsedData struct {
	Location     *GeoLocation `json:"location,omitempty"`
	HardwareID   string       `json:"hardwareId,omitempty"`
	Time         int64        `json:"serverTime,omitempty"`
	MessageID    string       `json:"messageId,omitempty"`
	Protocol     string       `json:"protocol,omitempty"`
	ServiceToken string       `json:"serviceToken,omitempty"`
	PositionTime int64        `json:"deviceTime,omitempty"`
	Indoor       *Indoor      `json:"Indoor,omitempty"`
	EventStatus  interface{}  `json:"EventStatus,omitempty"`
	Battery      interface{}  `json:"Battery,omitempty"`
	Light        interface{}  `json:"Light,omitempty"`
	Temperature  interface{}  `json:"AirTemperature,omitempty"`
	Humidity     interface{}  `json:"Humidity,omitempty"`
}
type Indoor struct {
	Building        string `json:"building,omitempty"`
	BuildingId      int32  `json:"buildingId,omitempty"`
	FloorIndex      int    `json:"floorIndex"`
	FloorLabel      string `json:"floorLabel,omitempty"`
	BuildingModelId int32  `json:"buildingModelId,omitempty"`
}

type GeoLocation struct {
	Unit   string  `json:"unit,omitempty"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Source string  `json:"source,omitempty"`
	Alt    float64 `json:"alt,omitempty"`
}

type Signals struct {
	Parsed_data interface{} `json:"rawRequest,omitempty"`
	CpsRequest  interface{} `json:"cpsRequest,omitempty"`
	CpsResponse interface{} `json:"cpsResponse,omitempty"`
	Signals     interface{} `json:"signals,omitempty"`
}

func Nats_message_handlers(msg []byte) []byte {
	fmt.Printf("Received message: %s\n", string(msg))
	var raw map[string]interface{}

	var nonlookupresponse *Signals

	err := json.Unmarshal(msg, &raw)
	if err != nil {
		log.Println("JSON parse error:", err)
		return nil
	}

	// --- Extract protocol ---
	protocol, ok := raw["protocol"].(string)
	if !ok {
		log.Println("protocol not found or invalid type")
		return nil
	}

	// --- Extract hardwareId ---
	// hardwareId, ok := raw["hardwareId"].(string)
	// if !ok {
	// 	log.Println("hardwareId not found or invalid type")
	// 	return nil
	// }

	// --- For specific protocol types, skip if radioData exists ---
	if protocol == "abeeway-compact-tracker" || protocol == "cpsflex" || protocol == "default-iot" {
		if rd, ok := raw["radioData"].(map[string]interface{}); !ok || len(rd) == 0 {
			parsed := ParsedData{}

			// --- Parse location ---
			if locMap, ok := raw["location"].(map[string]interface{}); ok {
				location := &GeoLocation{}
				if lat, ok := locMap["lat"].(float64); ok {
					location.Lat = lat
				}
				if lng, ok := locMap["lng"].(float64); ok {
					location.Lng = lng
				}
				if alt, ok := locMap["alt"].(float64); ok {
					location.Alt = alt
				}
				if unit, ok := locMap["unit"].(string); ok {
					location.Unit = unit
				}
				if source, ok := locMap["source"].(string); ok {
					location.Source = source
				}
				parsed.Location = location
			}

			// --- Parse indoor ---
			if indoorMap, ok := raw["Indoor"].(map[string]interface{}); ok {
				indoor := &Indoor{}
				if bld, ok := indoorMap["building"].(string); ok {
					indoor.Building = bld
				}
				if bldId, ok := indoorMap["buildingId"].(float64); ok {
					indoor.BuildingId = int32(bldId)
				}
				if label, ok := indoorMap["floorLabel"].(string); ok {
					indoor.FloorLabel = label
				}
				if modelId, ok := indoorMap["buildingModelId"].(float64); ok {
					indoor.BuildingModelId = int32(modelId)
				}
				if index, ok := indoorMap["floorIndex"].(float64); ok {
					indoor.FloorIndex = int(index)
				}
				parsed.Indoor = indoor
			}

			// --- Generic fields ---
			if v, ok := raw["hardwareId"].(string); ok {
				parsed.HardwareID = v
			}
			if v, ok := raw["messageId"].(string); ok {
				parsed.MessageID = v
			}
			if v, ok := raw["serviceToken"].(string); ok {
				parsed.ServiceToken = v
			}
			if v, ok := raw["time"].(float64); ok {
				parsed.Time = int64(v)
			}
			if v, ok := raw["positionTime"].(float64); ok {
				parsed.PositionTime = int64(v)
			}

			// --- Parse "signals" block ---
			if sigMap, ok := raw["signals"].(map[string]interface{}); ok {
				if bat, ok := sigMap["Battery"]; ok {
					parsed.Battery = bat
				}
				if lit, ok := sigMap["Light"]; ok {
					parsed.Light = lit
				}
				if temp, ok := sigMap["Air Temperature"]; ok {
					parsed.Temperature = temp
				}
				if hum, ok := sigMap["Humidity"]; ok {
					parsed.Humidity = hum
				}
				if ev, ok := sigMap["Event Status"]; ok {
					parsed.EventStatus = ev
				}
			}

			// --- Parse "mokoSignals" if available ---
			if sigMap, ok := raw["mokoSignals"].(map[string]interface{}); ok {
				mokoSignals := make(map[string]interface{})
				for k, v := range sigMap {
					if k == "Battery" {
						mokoSignals["batteryLevel"] = v
					} else {
						mokoSignals[k] = v
					}

				// --- Generic fields ---
				if v, ok := raw["hardwareId"].(string); ok {
					mokoSignals["hardwareId"]= v
				}
				if v, ok := raw["messageId"].(string); ok {
					mokoSignals["messageId"]=v
				}
				if v, ok := raw["serviceToken"].(string); ok {
					mokoSignals["serviceToken"] = v
				}
				if v, ok := raw["time"].(float64); ok {
					mokoSignals["serverTime"] = v
				}
				if v, ok := raw["positionTime"].(float64); ok {
					mokoSignals["positionTime"] = v
				}
				}

				delete(raw, "mokoSignals")

				nonlookupresponse = &Signals{
					Parsed_data: raw,
					Signals:     mokoSignals,
				}

				noncpsbytes, err := json.Marshal(nonlookupresponse)
				if err != nil {
					log.Println("Error marshalling mokoSignals response:", err)
					return nil
				}

				return noncpsbytes
			}

			// --- Default response using ParsedData ---
			delete(raw, "signals") // prevent duplication

			nonlookupresponse = &Signals{
				Parsed_data: raw,
				Signals:     parsed,
			}

			noncpsbytes, err := json.Marshal(nonlookupresponse)
			if err != nil {
				log.Println("Error marshalling ParsedData:", err)
				return nil
			}

			return noncpsbytes
			}
	}

	radioData, ok := raw["radioData"].(map[string]interface{})
	if !ok {
		log.Println("radioData not found or invalid type")
		return nil
	}

	// Convert radioData to []byte
	bodyBytes, err := json.Marshal(radioData)
	if err != nil {
		log.Println("Error marshalling radioData:", err)
		return nil
	}

	fmt.Printf("radioData: %s\n", string(bodyBytes))
	// fmt.Println("api key is:", opts.Apikey)

	// _, response := handleLookupRequest(
	// 	context.TODO(), // ctx
	// 	nil,            // remoteIP
	// 	"",             // method
	// 	opts.Apikey,    //apikey
	// 	http.Header{},  // headers
	// 	url.Values{},   // values
	// 	"nats",         // flag (important)
	// 	bodyBytes,      // body []byte (NATS payload)
	// 	time.Now(),     // requestStartTime
	// )
	// fmt.Println("response:", response)

	// bodyByteslookup := response.BodyBytes(nil)

	// // Step 1: Unmarshal the BodyBytes from the LookupResponse
	// var responseMap map[string]interface{}
	// if err := json.Unmarshal(bodyByteslookup, &responseMap); err != nil {
	// 	log.Println("Error unmarshalling LookupResponse.BodyBytes:", err)
	// 	return
	// }

	// Step 2: Check if error exists in lookup response
	// _, hasError := responseMap["error"]

	// Step 3: Create final signals map
	finalSignals := make(map[string]interface{})

	// Inject battery, temperature, etc.
	if signalsData, ok := raw["signals"].(map[string]interface{}); ok {
		if battery, ok := signalsData["Battery"]; ok {
			finalSignals["Battery"] = battery
		}
		if temp, ok := signalsData["Air Temperature"]; ok {
			finalSignals["AirTemperature"] = temp
		}
		if serverTime, ok := signalsData["received_time"]; ok {
			finalSignals["serverTime"] = serverTime
		}
		if light, ok := signalsData["Light"]; ok {
			finalSignals["Light"] = light
		}
		if event_status, ok := signalsData["Event Status"]; ok {
			finalSignals["EventStatus"] = event_status
		}
		if humidity, ok := signalsData["Humidity"].(interface{}); ok {
			finalSignals["Humidity"] = humidity
		}
	}

	// Inject hardwareId and messageId
	if hwID, ok := raw["hardwareId"]; ok {
		finalSignals["hardwareId"] = hwID
	}
	if msgID, ok := raw["messageId"]; ok {
		finalSignals["messageId"] = msgID
	}
	if positionTime, ok := raw["positionTime"]; ok {
		finalSignals["deviceTime"] = positionTime
	}

	// If NO error in lookup response ➔ inject entire lookup response into signals
	// if !hasError {
	// 	for k, v := range responseMap {
	// 		finalSignals[k] = v
	// 	}
	// 	finalSignals["position"] = true
	// } else {
	// 	finalSignals["position"] = nil
	// }

	delete(raw, "signals")

	// Step 4: Prepare signals_response
	signals_response := &Signals{
		Parsed_data: raw,
		CpsRequest:  radioData,
		// CpsResponse: responseMap,
		Signals:     finalSignals, // Injected response map
	}

	jsonBytes, err := json.Marshal(signals_response)
	if err != nil {
		log.Println("Error marshalling signals_response:", err)
		return nil
	} 
	fmt.Println(string(jsonBytes))

	return jsonBytes
}
