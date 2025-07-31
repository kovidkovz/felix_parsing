package protocols

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
)

type finalResponse struct {
	Parsed_data any `json:"rawRequest,omitempty"`
	CpsRequest  any `json:"cpsRequest,omitempty"`
	CpsResponse any `json:"cpsResponse,omitempty"`
	Signals     any `json:"signals,omitempty"`
}

func Nats_message_handler(msg []byte) []byte {
	fmt.Printf("Received message: %s\n", string(msg))
	var raw map[string]any
	var nonlookupresponse *finalResponse
	// call the lookup handler for the positioning data
	err := json.Unmarshal(msg, &raw)
	if err != nil {
		log.Println("JSON parse error:", err)
		return nil
	}

	// consume serviceToken directly from main keys
	serviceToken, ok := raw["serviceToken"].(string)
	if !ok {
		log.Println("serviceToken not found or invalid type", serviceToken)
		return nil
	}

	// extract protocol
	device_type, ok := raw["protocol"].(string)
	if !ok {
		log.Println("protocol not found or invalid type", device_type)
		return nil
	}

	// extract harwareId
	hardwareId, ok := raw["hardwareId"].(string)
	if !ok {
		log.Println("hardwareId not found or invalid type", hardwareId)
		return nil
	}

	signals := make(map[string]any)

	// Inject generic data into the signals
	if hwID, ok := raw["hardwareId"]; ok {
		signals["hardwareId"] = hwID
	}
	if msgID, ok := raw["messageId"]; ok {
		signals["messageId"] = msgID
	}
	if positionTime, ok := raw["positionTime"]; ok {
		signals["positionTime"] = positionTime
	}
	if serverTime, ok := raw["serverTime"]; ok {
		signals["serverTime"] = serverTime
	}
	signals["serviceToken"] = serviceToken

	// inject all that is in the signals array into the final signals
	if sigMap, ok := raw["signals"].(map[string]any); ok {
		maps.Copy(signals, sigMap)
	}

	if device_type == "abeeway-compact-tracker" || device_type == "cpsflex" || device_type == "default-iot" {
		if rd, ok := raw["radioData"].(map[string]any); !ok || len(rd) == 0 {
			// signals := make(map[string]interface{})

			// --- Parse location ---
			if locMap, ok := raw["location"].(map[string]any); ok {
				signals["location"] = locMap
			}

			// --- Parse indoor ---
			if indoorMap, ok := raw["Indoor"].(map[string]any); ok {
				signals["Indoor"] = indoorMap
			}

			delete(raw, "signals") // prevent duplication

			nonlookupresponse = &finalResponse{
				Parsed_data: raw,
				Signals:     signals,
			}

			noncpsbytes, err := json.Marshal(nonlookupresponse)
			if err != nil {
				log.Println("Error marshalling ParsedData:", err)
				return nil
			}
			return noncpsbytes
		}
	}

	radioData, ok := raw["radioData"].(map[string]any)
	if !ok {
		log.Println("radioData not found or invalid type")
		return nil
	}

	delete(raw, "signals")

	// Step 4: Prepare finalresponse
	final_response := &finalResponse{
		Parsed_data: raw,
		CpsRequest:  radioData,
		CpsResponse: nil,
		Signals:     signals,
	}

	jsonBytes, err := json.Marshal(final_response)
	if err != nil {
		log.Println("Error marshalling signals_response:", err)
		return nil
	}
	return jsonBytes
}
