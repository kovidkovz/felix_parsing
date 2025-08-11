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
    "id": "20635F03C1000352",
    "timestamp": "2025-08-10T14:41:52.364+00:00",
    "indoor": 1,
    "signals": {
      "dynamicMotionState": "MOVING",
      "messageType": "EVENT",
      "trackingMode": "MOTION_TRACKING",
      "battery": 86,
      "lastknownlat": 32.6291683,
      "lastknownlng": -83.593185,
      "lastpositionTime": "2025-08-10T14:35:56.468Z",
      "name": "Generated - 20635F03C1000352"
    },
    "wifiAccessPoints": [
      {}
    ]
  },
  "hardwareId": "20635f03c1000352",
  "time": 1754836913005,
  "messageId": "78bdc491-825d-47e3-89a5-124421e81306",
  "serviceToken": "aca80396-eda9-42ad-a4f0-005c88dc75de",
  "device_profile_name": "abeeway-compact-tracker"
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
