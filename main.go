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
      "Time": "2025-09-09T08:00:44.928+00:00",
      "DevEUI": "20635F0241001349",
      "FPort": 18,
      "FCntUp": 9635,
      "LostUplinksAS": 0,
      "ADRbit": 1,
      "MType": 2,
      "FCntDn": 99,
      "payload_hex": "052039848040020600030305",
      "mic_hex": "582fe3a8",
      "Lrcid": "000000CB",
      "LrrRSSI": -32,
      "LrrSNR": 9.25,
      "LrrESP": -32.48772,
      "SpFact": 7,
      "SubBand": "G0",
      "Channel": "LC0",
      "Lrrid": "10000511",
      "Late": 0,
      "LrrLAT": 32.592274,
      "LrrLON": -83.632484,
      "Lrrs": {
        "Lrr": [
          {
            "Lrrid": "10000511",
            "Chain": 0,
            "LrrRSSI": -32,
            "LrrSNR": 9.25,
            "LrrESP": -32.48772
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
        "name": "SB3-1349"
      },
      "BaseStationData": {
        "doms": [],
        "name": "ACES-GW3"
      },
      "ModelCfg": "1:TPX_470fcf3b-a998-4735-a915-a99dc2d3a08c",
      "DriverCfg": {
        "mod": {
          "pId": "abeeway",
          "mId": "smart-badge",
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
      "MeanPER": 0,
      "DevAddr": "057B1E9E",
      "TxPower": 10,
      "NbTrans": 1,
      "Frequency": 902.3,
      "DynamicClass": "A",
      "PayloadEncryption": 0,
      "payload": {
        "messageType": "HEARTBEAT",
        "trackingMode": "MOTION_TRACKING",
        "batteryLevel": 57,
        "batteryStatus": "OPERATING",
        "ackToken": 8,
        "firmwareVersion": "2.6.0",
        "bleFirmwareVersion": "3.3.5",
        "resetCause": 64,
        "periodicPosition": "false",
        "temperatureMeasure": 22.8,
        "sosFlag": 0,
        "appState": 0,
        "dynamicMotionState": "STATIC",
        "onDemand": "false",
        "payload": "052039848040020600030305",
        "deviceConfiguration": {
          "mode": "MOTION_TRACKING"
        }
      },
      "points": {
        "batteryLevel": {
          "unitId": "%",
          "record": 57
        },
        "temperature": {
          "unitId": "Cel",
          "record": 22.8
        }
      },
      "downlinkUrl": "https://thingparkenterprise.us.actility.com/iot-flow/downlinkMessages/39eb958d-9184-4200-8ae4-27eefa53f8d6"
    }
  },
  "ruleStates": {
    "164067": false,
    "164067-inside": false
  },
  "deviceId": 270510,
  "signals": {
    "Received Time": "2025-09-09T08:00:49.340Z",
    "Mode": "Heart beat",
    "Battery Level": 57,
    "Temperature": 22.8
  }
}
`


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
