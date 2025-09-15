// package main

// import (
// 	"fmt"
// 	"ms-testing/protocols"
// 	// "gitlab.combain.com/traxmate/natstemplate/natstemplate"
// )

// // func main() {
// // 	// initialize the nats connection for the producer
// // 	_, js, err := natstemplate.NatsProducerInstance.InitializeConnection()
// // 	if err != nil {
// // 		log.Fatalf("Failed to initialize NATS producer: %v", err)
// // 	}

// // 	log.Println("NATS connection initialized successfully!")

// // 	go natstemplate.Dynamic_Consumer_Creation("rawdata.", "rawdata_", "", js, processMessage)

// // 	// Keep the service running
// // 	select {}
// // }

// func main() {
// 	// Yahan tera JSON payload (shortened example)
// 	jsonPayload := `
// {"data":{"DevEUI_uplink":{"Time":"2025-09-10T12:02:26.581+00:00","DevEUI":"20635F03C10001A7","FPort":18,"FCntUp":46584,"LostUplinksAS":0,"ADRbit":1,"MType":2,"FCntDn":804,"payload_hex":"052046768040020600030305","mic_hex":"207e9665","Lrcid":"000000CB","LrrRSSI":-91.000000,"LrrSNR":12.250000,"LrrESP":-91.251282,"SpFact":9,"SubBand":"G0","Channel":"LC0","Lrrid":"10000506","Late":0,"LrrLAT":32.630306,"LrrLON":-83.600533,"Lrrs":{"Lrr":[{"Lrrid":"10000506","Chain":0,"LrrRSSI":-91.000000,"LrrSNR":12.250000,"LrrESP":-91.251282},{"Lrrid":"10000513","Chain":0,"LrrRSSI":-73.000000,"LrrSNR":11.000000,"LrrESP":-73.331955}]},"DevLrrCnt":2,"CustomerID":"100000184","CustomerData":{"loc":null,"alr":{"pro":"ABEE/APY","ver":"1"},"tags":[],"doms":[],"name":"ACES-01A7"},"BaseStationData":{"doms":[],"name":"ACES-GW1"},"ModelCfg":"0","DriverCfg":{"mod":{"pId":"abeeway","mId":"compact-tracker","ver":"1"},"app":{"pId":"abeeway","mId":"asset-tracker","ver":"2"}},"InstantPER":0.000000,"MeanPER":0.019608,"DevAddr":"01128CE8","TxPower":18.000000,"NbTrans":1,"Frequency":902.3,"DynamicClass":"A","PayloadEncryption":0}},"hardwareId":"","time":1757505746947,"messageId":"5babe855-5f7e-4557-acc1-a74aa28be7fe","serviceToken":"aca80396-eda9-42ad-a4f0-005c88dc75de","device_profile_name":"abeeway-compact-tracker"}


// `

// // // Call function directly with []byte
// // result := protocols.Parse_felix_data([]byte(jsonPayload))
// // fmt.Println(string(result))

// result := protocols.ProcessAlaeMessage([]byte(jsonPayload))
// fmt.Println(string(result))

// // result := protocols.Nats_message_handlers([]byte(jsonPayload))
// // fmt.Println(string(result))

// // result, err := FallbackLookup([]byte(jsonPayload))
// // if err != nil {
// //   fmt.Println("error", err)
// // }
// // fmt.Println(result)
// }

package main

import (
	"encoding/hex"
	"fmt"
	"reflect"
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

	Bssid     []string // For wifi type
	Rssi      []int    // For wifi or BLE

	MacAdr    []string // For BLE
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

	decoded.SosMode = (bytes[1]&0x10) != 0
	decoded.TrackingState = (bytes[1]&0x08) != 0
	decoded.Moving = (bytes[1]&0x04) != 0
	decoded.PeriodicPos = (bytes[1]&0x02) != 0
	decoded.PosOnDemand = (bytes[1]&0x01) != 0
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

func main() {
	payloadHex := "05205b826040020600030305"
	bytes, err := hex.DecodeString(payloadHex)
	if err != nil {
		fmt.Println("Error decoding hex:", err)
		return
	}
	data := Decoder(bytes, 17)
	// fmt.Printf("%+v\n", data)

  out := StructToMap(data)

  fmt.Println("out", out)
	//   fmt.Printf("Battery Voltage: %.3f V\n", data.BatteryVoltage)
	//   fmt.Printf("Battery Percentage: %d%%\n", data.BatteryPercentage)
	//   fmt.Printf("Temperature: %.1f °C\n", data.Temperature)
	//   fmt.Printf("Ack Token: %d\n", data.AckToken)
	//   fmt.Printf("Sos Mode: %v\n", data.SosMode)
	//   fmt.Printf("Tracking State: %v\n", data.TrackingState)
	//   fmt.Printf("Moving: %v\n", data.Moving)
	//   fmt.Printf("Periodic Pos: %v\n", data.PeriodicPos)
	//   fmt.Printf("Pos On Demand: %v\n", data.PosOnDemand)
	//   fmt.Printf("Operating Mode: %d\n", data.OperatingMode)

	//   if data.Latitude != nil && data.Longitude != nil {
	//       fmt.Printf("Latitude: %.8f\n", *data.Latitude)
	//       fmt.Printf("Longitude: %.8f\n", *data.Longitude)
	//   }
	//   if data.Accuracy != nil {
	//       fmt.Printf("Accuracy: %.2f\n", *data.Accuracy)
	//   }
	//   if data.Age != nil {
	//       fmt.Printf("Age: %d\n", *data.Age)
	//   }
	// // ... repeat for any other pointer fields

	//	if data.Latitude != nil && data.Longitude != nil {
	//		fmt.Printf("Latitude: %.8f, Longitude: %.8f\n", *data.Longitude, *data.Longitude)
	//	}
}
