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
{"data":{"data":{"notifications":[{"locationHierarchy":{"level":"ROOT","name":"FMS Office","id":"daca73b9-3dce-4a18-97eb-8a7f9d6b17be","child":{"level":"ORG","name":"FMS Office Network","id":"03d6ac81-c867-4295-a869-06d02f8559ca","child":{"level":"NETWORK","name":"FMS Office","id":"67e9439f-4102-46ad-917d-4fd095dbd230","child":{"level":"FLOOR","name":"FMS Office","id":"bfdf0ea8-f040-421f-84f6-0cea1d00e3f8"}}}},"floorUuid":"bfdf0ea8-f040-421f-84f6-0cea1d00e3f8","notificationType":"locationupdate","source":"UNKNOWN","deviceId":"c4:39:60:6f:e3:ed","geoCoordinate":{"unit":"DEGREES","latitude":30.226234204449092,"longitude":-97.7594865997888},"ssid":"","manufacturer":"","floorId":"","band":"","generatedBy":"dnl-data-processor","staticDevice":"false","timestamp":1753773195819,"eventId":1753773195931,"locComputeType":"RSSI","lastLocated":"2025-07-29 07:13:15.819+0000","ipAddress":[""],"userName":"","lastSeen":"2025-07-29 07:13:15.819+0000","apMacAddress":"e4:55:a8:92:1f:ff","subscriptionName":"Traxmatev3","locationMapHierarchy":"FMS Office Network-\u003eFMS Office-\u003eFMS Office","associated":false,"tenantId":"29010","locationCoordinate":{"unit":"FEET","x":156.27429,"y":32.013428,"z":0},"confidenceFactor":0,"subscriptionId":"89ecae31-5d2e-43a2-86fd-36103845e1df","entity":"BLE_TAGS"}]},"hardwareId":"","time":1753776724198,"messageId":"bc3797be-eac3-44f0-904a-b6808ed8ad6a","serviceToken":"0941d3a4-8b43-490d-b90f-161b0a04c1ac","protocol":"default-iot"},"location":{"unit":"DEGREES","lat":30.226234204449092,"lng":-97.7594865997888},"hardwareId":"c4:39:60:6f:e3:ed","time":1753776724198,"messageId":"bc3797be-eac3-44f0-904a-b6808ed8ad6a","signals":{"device_name":"c4:39:60:6f:e3:ed"},"protocol":"default-iot","serviceToken":"0941d3a4-8b43-490d-b90f-161b0a04c1ac","Indoor":{"building":"East Alpine Road 230","buildingId":11878987,"floorIndex":0,"floorLabel":"0","buildingModelId":159643361,"locationHierarchy":{"child":{"child":{"child":{"id":"bfdf0ea8-f040-421f-84f6-0cea1d00e3f8","level":"FLOOR","name":"FMS Office"},"id":"67e9439f-4102-46ad-917d-4fd095dbd230","level":"NETWORK","name":"FMS Office"},"id":"03d6ac81-c867-4295-a869-06d02f8559ca","level":"ORG","name":"FMS Office Network"},"id":"daca73b9-3dce-4a18-97eb-8a7f9d6b17be","level":"ROOT","name":"FMS Office"}},"positionTime":1753773195819}`

// // Call function directly with []byte
// result := protocols.Parse_felix_data([]byte(jsonPayload))
// fmt.Println(string(result))


result := protocols.Nats_message_handlers([]byte(jsonPayload))
fmt.Println(string(result))

// result, err := FallbackLookup([]byte(jsonPayload))
// if err != nil {
//   fmt.Println("error", err)
// }
// fmt.Println(result)
}
