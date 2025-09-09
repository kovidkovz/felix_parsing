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
      "Time": "2025-09-09T06:38:20.519+00:00",
      "DevEUI": "20635F00C80002A6",
      "FPort": 17,
      "FCntUp": 3924,
      "LostUplinksAS": 0,
      "ADRbit": 1,
      "MType": 2,
      "FCntDn": 71,
      "payload_hex": "032a9088a00a136d2dce26ab02",
      "mic_hex": "17b3594a",
      "Lrcid": "000000CB",
      "LrrRSSI": -39,
      "LrrSNR": 9.25,
      "LrrESP": -39.48772,
      "SpFact": 7,
      "SubBand": "G0",
      "Channel": "LC2",
      "Lrrid": "10000511",
      "Late": 0,
      "LrrLAT": 32.592274,
      "LrrLON": -83.632484,
      "Lrrs": {
        "Lrr": [
          {
            "Lrrid": "10000511",
            "Chain": 0,
            "LrrRSSI": -39,
            "LrrSNR": 9.25,
            "LrrESP": -39.48772
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
        "name": "ACES-02A6-IT"
      },
      "BaseStationData": {
        "doms": [],
        "name": "ACES-GW3"
      },
      "ModelCfg": "1:TPX_470fcf3b-a998-4735-a915-a99dc2d3a08c",
      "DriverCfg": {
        "mod": {
          "pId": "abeeway",
          "mId": "indus-tracker",
          "ver": "2"
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
      "DevAddr": "0408C1A3",
      "TxPower": 10,
      "NbTrans": 1,
      "Frequency": 902.7,
      "DynamicClass": "A",
      "PayloadEncryption": 0,
      "payload": {
        "gpsLatitude": 32.5922048,
        "gpsLongitude": -83.6326656,
        "horizontalAccuracy": 7.84,
        "messageType": "POSITION_MESSAGE",
        "age": 80,
        "trackingMode": "MOTION_TRACKING",
        "batteryVoltage": 3.59,
        "ackToken": 10,
        "rawPositionType": "GPS",
        "periodicPosition": "true",
        "temperatureMeasure": 24.8,
        "sosFlag": 0,
        "appState": 1,
        "dynamicMotionState": "STATIC",
        "onDemand": "false",
        "payload": "032a9088a00a136d2dce26ab02",
        "deviceConfiguration": {
          "mode": "MOTION_TRACKING"
        }
      },
      "points": {
        "temperature": {
          "unitId": "Cel",
          "record": 24.8
        },
        "batteryVoltage": {
          "unitId": "V",
          "record": 3.59
        },
        "location": {
          "unitId": "GPS",
          "record": [
            -83.6326656,
            32.5922048
          ]
        },
        "accuracy": {
          "unitId": "m",
          "record": 7.84
        },
        "age": {
          "unitId": "s",
          "record": 80
        }
      },
      "downlinkUrl": "https://thingparkenterprise.us.actility.com/iot-flow/downlinkMessages/39eb958d-9184-4200-8ae4-27eefa53f8d6"
    }
  },
  "ruleStates": {
    "164067": false,
    "164067-inside": false
  },
  "deviceId": 272092,
  "signals": {
    "Received Time": "2025-09-09T06:38:21.142Z",
    "Mode": "Position",
    "Temperature": 24.8,
    "Location": {
      "lat": 32.5922048,
      "lng": -83.6326656
    },
    "Position Time": "2025-09-09T06:37:01.122Z"
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
