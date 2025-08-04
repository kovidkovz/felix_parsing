package main

import (
	"fmt"
	"ms-testing/protocols"
	// "gitlab.combain.com/traxmate/natstemplate/natstemplate"
)

// func main() {
// 	// initialize the nats connection for the producer
// 	_, js, err := natstemplate.NatsProducerInstance.InitializeConnection()
// 	if err != nil {
// 		log.Fatalf("Failed to initialize NATS producer: %v", err)
// 	}

// 	log.Println("NATS connection initialized successfully!")

// 	go natstemplate.Dynamic_Consumer_Creation("rawdata.", "rawdata_", "", js, processMessage)

// 	// Keep the service running
// 	select {}
// }

func main() {
	// Yahan tera JSON payload (shortened example)
	jsonPayload := `
{
  "data": {
    "DevEUI_uplink": {
      "Time": "2025-04-15T06:08:10.510+00:00",
      "DevEUI": "20635F03C1000493",
      "FPort": 18,
      "FCntUp": 8106,
      "LostUplinksAS": 0,
      "ADRbit": 1,
      "MType": 2,
      "FCntDn": 153,
      "payload_hex": "0520597d5040020600030305",
      "mic_hex": "f95d8fc5",
      "Lrcid": "000000CB",
      "LrrRSSI": -65,
      "LrrSNR": 9,
      "LrrESP": -65.514969,
      "SpFact": 8,
      "SubBand": "G0",
      "Channel": "LC5",
      "Lrrid": "10000506",
      "Late": 0,
      "Lrrs": {
        "Lrr": [
          {
            "Lrrid": "10000506",
            "Chain": 0,
            "LrrRSSI": -65,
            "LrrSNR": 9,
            "LrrESP": -65.514969
          }
        ]
      },
      "DevLrrCnt": 1,
      "CustomerID": "100000184",
      "CustomerData": {
        "loc": null,
        "alr": {
          "pro": "ABEE/APY",
          "ver": "1"
        },
        "tags": [],
        "doms": [],
        "name": "ACES-0493"
      },
      "BaseStationData": {
        "doms": [],
        "name": "ACES-GW1"
      },
      "ModelCfg": "1:TPX_470fcf3b-a998-4735-a915-a99dc2d3a08c",
      "DriverCfg": {
        "mod": {
          "pId": "abeeway",
          "mId": "compact-tracker",
          "ver": "1"
        },
        "app": {
          "pId": "abeeway",
          "mId": "asset-tracker",
          "ver": "2"
        },
        "id": "abeeway:asset-tracker:3"
      },
      "InstantPER": 0,
      "MeanPER": 0.019608,
      "DevAddr": "00DAAC2D",
      "TxPower": 18,
      "NbTrans": 1,
      "Frequency": 903.3,
      "DynamicClass": "A",
      "payload": {
        "messageType": "HEARTBEAT",
        "trackingMode": "MOTION_TRACKING",
        "batteryLevel": 89,
        "batteryStatus": "OPERATING",
        "ackToken": 5,
        "firmwareVersion": "2.6.0",
        "bleFirmwareVersion": "3.3.5",
        "resetCause": 64,
        "periodicPosition": "false",
        "temperatureMeasure": 19.2,
        "sosFlag": 0,
        "appState": 0,
        "dynamicMotionState": "STATIC",
        "onDemand": "false",
        "payload": "0520597d5040020600030305",
        "deviceConfiguration": {
          "mode": "MOTION_TRACKING"
        }
      },
      "points": {
        "batteryLevel": {
          "unitId": "%",
          "record": 89
        },
        "temperature": {
          "unitId": "Cel",
          "record": 19.2
        }
      },
      "downlinkUrl": "https://thingparkenterprise.us.actility.com/iot-flow/downlinkMessages/39eb958d-9184-4200-8ae4-27eefa53f8d6"
    }
  }
}`

// // Call function directly with []byte
// result := protocols.Parse_felix_data([]byte(jsonPayload))
// fmt.Println(string(result))


result := protocols.ProcessAlaeMessage([]byte(jsonPayload))
fmt.Println(string(result))

// result := protocols.Nats_message_handlers([]byte(jsonPayload))
// fmt.Println(string(result))

// result, err := FallbackLookup([]byte(jsonPayload))
// if err != nil {
//   fmt.Println("error", err)
// }
// fmt.Println(result)
}
