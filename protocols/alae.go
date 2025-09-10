package protocols

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"ms-testing/models"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// DecodedData matches the structure of the JS decoded object
type DecodedData struct {
	BatteryVoltage    float64
	BatteryPercentage int
	Temperature       float64
	AckToken          int
	SosMode           bool
	TrackingState     bool
	Moving            bool
	PeriodicPos       bool
	PosOnDemand       bool
	OperatingMode     int

	Latitude  *float64
	Longitude *float64
	Accuracy  *float64
	Age       *int

	Bssid []string // For wifi type
	Rssi  []int    // For wifi or BLE

	MacAdr      []string // For BLE
	GpsTimeout  bool
	Shutdown    bool
	GeolocStart bool
	Heartbeat   bool
	ResetCause  *int
	FirmwareVer []byte
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

func Decoder(bytes []byte, port int) DecodedData {
	decoded := DecodedData{}

	if len(bytes) < 5 {
		return decoded
	}

	decoded.BatteryVoltage = float64(bytes[2])*0.0055 + 2.8
	decoded.BatteryPercentage = int(float64(bytes[2]) / 255.0 * 100)
	decoded.Temperature = float64(bytes[3])*0.5 - 44
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

	return decoded
}

// Signed int8 conversion
func signedByte(b byte) int {
	if b > 127 {
		return int(b) - 256
	}
	return int(b)
}

func StructToMap(data DecodedData) map[string]interface{} {
	out := make(map[string]interface{})
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)
	for i := 0; i < val.NumField(); i++ {
		key := typ.Field(i).Name
		value := val.Field(i).Interface()
		// If pointer, dereference
		if v, ok := value.(*float64); ok && v != nil {
			value = *v
		}
		if v, ok := value.(*int); ok && v != nil {
			value = *v
		}
		out[key] = value
	}
	return out
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
	var protocol string
	var lat *float64 = nil
	var lng *float64 = nil

	// Case A: DevEUI_uplink (semtech/raw lora type)
	if uplink, found := data["DevEUI_uplink"].(map[string]any); found {
		// Hardware
		if devEUI, ok := uplink["DevEUI"].(string); ok {
			hardwareID = devEUI
			signals["deviceName"] = devEUI
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

		decoded_map := StructToMap(decoded)

		positionTime = extractTimeField(uplink["Time"])


		for k, v := range decoded_map {
			switch k {
			case "gpsLatitude":
				if f, ok := v.(float64); ok {
					lat = &f
				}
			case "gpsLongitude":
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
