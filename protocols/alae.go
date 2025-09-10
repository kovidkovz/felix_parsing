package protocols

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"ms-testing/models"
	// "reflect"
	"strconv"
	"strings"
	"time"
)

// DecodedData matches the structure of the JS decoded object
type DecodedData struct {
	BatteryVoltage    float64  `json:"batteryVoltage"`
	BatteryLevel      int      `json:"batteryLevel"`
	TemperatureLevel  float64  `json:"temperatureLevel"`
	AckToken          int      `json:"ackToken"`
	SosMode           bool     `json:"sosMode"`
	TrackingState     bool     `json:"trackingState"`
	Moving            bool     `json:"moving"`
	PeriodicPos       bool     `json:"periodicPos"`
	PosOnDemand       bool     `json:"posOnDemand"`
	OperatingMode     int      `json:"operatingMode"`
	Latitude          *float64 `json:"latitude,omitempty"`
	Longitude         *float64 `json:"longitude,omitempty"`
	Accuracy          *float64 `json:"accuracy,omitempty"`
	Age               *int     `json:"age,omitempty"`
	Bssid             []string `json:"bssid,omitempty"`
	Rssi              []int    `json:"rssi,omitempty"`
	MacAdr            []string `json:"macAdr,omitempty"`
	GpsTimeout        bool     `json:"gpsTimeout,omitempty"`
	Shutdown          bool     `json:"shutdown,omitempty"`
	GeolocStart       bool     `json:"geolocStart,omitempty"`
	Heartbeat         bool     `json:"heartbeat,omitempty"`
	ResetCause        *int     `json:"resetCause,omitempty"`
	FirmwareVer       []byte   `json:"firmwareVer,omitempty"`
}

func formatMAC(bytes []byte) string {
	vals := make([]string, len(bytes))
	for i, b := range bytes {
		vals[i] = fmt.Sprintf("%02x", b)
	}
	return join(vals[:], ":") // No fmt.Sprintf needed here!
}

func join(a []string, sep string) string {
	out := ""
	for i := 0; i < len(a); i++ {
		if i > 0 {
			out += sep
		}
		out += a[i]
	}
	return out
}

func Decoder(bytes []byte, port int) []byte {
	decoded := DecodedData{}

	if len(bytes) < 5 {
		return nil
	}

	decoded.BatteryVoltage = float64(bytes[2])*0.0055 + 2.8
	decoded.BatteryLevel = int(float64(bytes[2]) / 255.0 * 100)
	decoded.TemperatureLevel = float64(bytes[3])*0.5 - 44
	decoded.AckToken = int(bytes[4] >> 4)

	decoded.SosMode = (bytes[1] & 0x10) != 0
	decoded.TrackingState = (bytes[1] & 0x08) != 0
	decoded.Moving = (bytes[1] & 0x04) != 0
	decoded.PeriodicPos = (bytes[1] & 0x02) != 0
	decoded.PosOnDemand = (bytes[1] & 0x01) != 0
	decoded.OperatingMode = int(bytes[1] >> 5)

	switch {
	// Position message & GPS type
	case bytes[0] == 0x03 && (bytes[4]&0x0F) == 0x00 && len(bytes) >= 13:
		latRawUint := (uint32(bytes[6]) << 16) | (uint32(bytes[7]) << 8) | uint32(bytes[8])
		latRawUint = latRawUint << 8
		latRaw := int32(latRawUint)
		latitude := float64(latRaw) / 10000000.0
		decoded.Latitude = &latitude

		lngRawUint := (uint32(bytes[9]) << 16) | (uint32(bytes[10]) << 8) | uint32(bytes[11])
		lngRawUint = lngRawUint << 8
		lngRaw := int32(lngRawUint)
		longitude := float64(lngRaw) / 10000000.0
		decoded.Longitude = &longitude

		accuracy := float64(bytes[12]) * 3.9
		decoded.Accuracy = &accuracy
		age := int(bytes[5]) * 8
		decoded.Age = &age

	// Position message & WiFi BSSID type
	case bytes[0] == 0x03 && (bytes[4]&0x0F) == 0x09 && len(bytes) >= 34:
		decoded.Bssid = []string{
			formatMAC(bytes[6:12]),
			formatMAC(bytes[13:19]),
			formatMAC(bytes[20:26]),
			formatMAC(bytes[27:33]),
		}
		decoded.Rssi = []int{
			signedByte(bytes[12]),
			signedByte(bytes[19]),
			signedByte(bytes[26]),
			signedByte(bytes[33]),
		}

		// Position message & BLE macaddr type
	case bytes[0] == 0x03 && (bytes[4]&0x0F) == 0x07 && len(bytes) >= 34:
		decoded.MacAdr = []string{
			formatMAC(bytes[6:12]),
			formatMAC(bytes[13:19]),
			formatMAC(bytes[20:26]),
			formatMAC(bytes[27:33]),
		}
		decoded.Rssi = []int{
			signedByte(bytes[12]),
			signedByte(bytes[19]),
			signedByte(bytes[26]),
			signedByte(bytes[33]),
		}

		// Position message & GPS timeout (failure)
	case bytes[0] == 0x03 && (bytes[4]&0x0F) == 0x01:
		decoded.GpsTimeout = true

		// Shutdown message
	case bytes[0] == 0x09:
		decoded.Shutdown = true

		// Geoloc start
	case bytes[0] == 0x0A:
		decoded.GeolocStart = true

		// Heartbeat
	case bytes[0] == 0x05:
		decoded.Heartbeat = true
		if len(bytes) >= 6 {
			tmp := int(bytes[5])
			decoded.ResetCause = &tmp
		}
		if len(bytes) >= 9 {
			decoded.FirmwareVer = bytes[6:9]
		}
	}

	bytes, err := json.Marshal(decoded)
	if err != nil {
		fmt.Println("Error marshalling", err)
	}

	return bytes
}

// Signed int8 conversion
func signedByte(b byte) int {
	if b > 127 {
		return int(b) - 256
	}
	return int(b)
}

func BytesToMap(data []byte) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal(data, &m)
	return m, err
}

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
	var lat *float64 = nil
	var lng *float64 = nil

	// Case A: DevEUI_uplink (semtech/raw lora type)
	if uplink, found := data["DevEUI_uplink"].(map[string]any); found {
		// Hardware
		if devEUI, ok := uplink["DevEUI"].(string); ok {
			hardwareID = devEUI
			signals["deviceName"] = devEUI
		}

		payloadHex, ok := uplink["payload_hex"].(string)
		if !ok {
			fmt.Println("payload_hex missing or not a string")
		}

		bytes, err := hex.DecodeString(payloadHex)
		if err != nil {
			fmt.Println("Error decoding hex:", err)
			return nil
		}

		decoded := Decoder(bytes, 17)

		decoded_map, err := BytesToMap(decoded)
		if err != nil {
			fmt.Println("Error:", err)
		}

		positionTime = extractTimeField(uplink["Time"])

		for k, v := range decoded_map {
			switch k {
			case "latitude":
				if f, ok := v.(float64); ok {
					lat = &f
				}
			case "longitude":
				if f, ok := v.(float64); ok {
					lng = &f
				}
			case "temperatureMeasure":
				signals["temperatureLevel"] = v
			default:
				signals[k] = v
			}
		}
		if lat != nil && lng != nil {
			geo = &models.GeoLocation{
				Lat: *lat,
				Lng: *lng,
			}
		}
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

	// Final assemble
	incoming := models.IncomingData{
		HardwareID:   hardwareID,
		Time:         raw["time"],
		MessageID:    fmt.Sprintf("%v", raw["messageId"]),
		Protocol:     fmt.Sprintf("%v", raw["device_profile_name"]),
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
