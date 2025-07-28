package main

import (
	"fmt"
	// "ms-testing/protocols"
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
    "wifiAccessPoints": [{
        "macAddress": "47:47:89:f1:c8:e3",
        "signalStrength": "-56"
    }, {
        "macAddress": "46:b8:ee:04:35:75",
        "signalStrength": "-56"
    }, {
        "macAddress": "a8:46:9d:2b:15:bb",
        "signalStrength": "-83"
    }]
}
`

// // Call function directly with []byte
// result := protocols.Parse_felix_data([]byte(jsonPayload))
// fmt.Println(string(result))


// result := protocols.Nats_message_handlers([]byte(jsonPayload))
// fmt.Println(string(result))

result, err := FallbackLookup([]byte(jsonPayload))
if err != nil {
  fmt.Println("error", err)
}
fmt.Println(result)
}
