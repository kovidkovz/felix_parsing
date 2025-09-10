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
// {
//   "data": {
//     "DevEUI_uplink": {
//       "Time": "2025-09-09T06:38:20.519+00:00",
//       "DevEUI": "20635F00C80002A6",
//       "FPort": 17,
//       "FCntUp": 3924,
//       "LostUplinksAS": 0,
//       "ADRbit": 1,
//       "MType": 2,
//       "FCntDn": 71,
//       "payload_hex": "032a9088a00a136d2dce26ab02",
//       "mic_hex": "17b3594a",
//       "Lrcid": "000000CB",
//       "LrrRSSI": -39,
//       "LrrSNR": 9.25,
//       "LrrESP": -39.48772,
//       "SpFact": 7,
//       "SubBand": "G0",
//       "Channel": "LC2",
//       "Lrrid": "10000511",
//       "Late": 0,
//       "LrrLAT": 32.592274,
//       "LrrLON": -83.632484,
//       "Lrrs": {
//         "Lrr": [
//           {
//             "Lrrid": "10000511",
//             "Chain": 0,
//             "LrrRSSI": -39,
//             "LrrSNR": 9.25,
//             "LrrESP": -39.48772
//           }
//         ]
//       },
//       "DevLrrCnt": 1,
//       "CustomerID": "100000184",
//       "CustomerData": {
//         "loc": null,
//         "alr": {
//           "pro": "ABEE/APY",
//           "ver": "1"
//         },
//         "tags": [],
//         "doms": [],
//         "name": "ACES-02A6-IT"
//       },
//       "BaseStationData": {
//         "doms": [],
//         "name": "ACES-GW3"
//       },
//       "ModelCfg": "1:TPX_470fcf3b-a998-4735-a915-a99dc2d3a08c",
//       "DriverCfg": {
//         "mod": {
//           "pId": "abeeway",
//           "mId": "indus-tracker",
//           "ver": "2"
//         },
//         "app": {
//           "pId": "abeeway",
//           "mId": "asset-tracker",
//           "ver": "2"
//         },
//         "id": "abeeway:asset-tracker:3"
//       },
//       "InstantPER": 0,
//       "MeanPER": 0.019608,
//       "DevAddr": "0408C1A3",
//       "TxPower": 10,
//       "NbTrans": 1,
//       "Frequency": 902.7,
//       "DynamicClass": "A",
//       "PayloadEncryption": 0,
//       "payload": {
//         "gpsLatitude": 32.5922048,
//         "gpsLongitude": -83.6326656,
//         "horizontalAccuracy": 7.84,
//         "messageType": "POSITION_MESSAGE",
//         "age": 80,
//         "trackingMode": "MOTION_TRACKING",
//         "batteryVoltage": 3.59,
//         "ackToken": 10,
//         "rawPositionType": "GPS",
//         "periodicPosition": "true",
//         "temperatureMeasure": 24.8,
//         "sosFlag": 0,
//         "appState": 1,
//         "dynamicMotionState": "STATIC",
//         "onDemand": "false",
//         "payload": "032a9088a00a136d2dce26ab02",
//         "deviceConfiguration": {
//           "mode": "MOTION_TRACKING"
//         }
//       },
//       "points": {
//         "temperature": {
//           "unitId": "Cel",
//           "record": 24.8
//         },
//         "batteryVoltage": {
//           "unitId": "V",
//           "record": 3.59
//         },
//         "location": {
//           "unitId": "GPS",
//           "record": [
//             -83.6326656,
//             32.5922048
//           ]
//         },
//         "accuracy": {
//           "unitId": "m",
//           "record": 7.84
//         },
//         "age": {
//           "unitId": "s",
//           "record": 80
//         }
//       },
//       "downlinkUrl": "https://thingparkenterprise.us.actility.com/iot-flow/downlinkMessages/39eb958d-9184-4200-8ae4-27eefa53f8d6"
//     }
//   },
//   "ruleStates": {
//     "164067": false,
//     "164067-inside": false
//   },
//   "deviceId": 272092,
//   "signals": {
//     "Received Time": "2025-09-09T06:38:21.142Z",
//     "Mode": "Position",
//     "Temperature": 24.8,
//     "Location": {
//       "lat": 32.5922048,
//       "lng": -83.6326656
//     },
//     "Position Time": "2025-09-09T06:37:01.122Z"
//   }
// }
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
	payloadHex := "032a9386a011136d31ce26ad02"
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
