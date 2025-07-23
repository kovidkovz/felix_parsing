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
	{"data":{"data":{"correlation_ids":["gs:uplink:01K0SGMBHNAADG57P80B9E8GZB"],"end_device_ids":{"application_ids":{"application_id":"fms-demo-smart-office"},"dev_addr":"27FD4831","dev_eui":"20CA0DBB117CDCB7","device_id":"thpo-dcb7","join_eui":"2CA2FA3AD4CB9A49"},"received_at":"2025-07-22T16:49:22.180616322Z","uplink_message":{"confirmed":true,"consumed_airtime":"0.185344s","decoded_payload":{"humidity":47,"motion_detected":false,"pressure":1890,"raw_payload":"00 e8 01 d6 00 63 00 00","status":{"humidity_valid":true,"pressure_valid":false,"temperature_valid":true},"temperature":{"celsius":23.2,"fahrenheit":73.75999999999999}},"f_cnt":67,"f_port":1,"frm_payload":"AOgB1gBjAAA=","network_ids":{"cluster_address":"nam1.cloud.thethings.industries","cluster_id":"nam1","net_id":"000013","ns_id":"EC656E0000102C53","tenant_address":"fmsiotcloud.nam1.cloud.thethings.industries","tenant_id":"fmsiotcloud"},"packet_error_rate":0.04761905,"received_at":"2025-07-22T16:49:21.974592647Z","rx_metadata":[{"channel_index":3,"channel_rssi":-48,"frequency_offset":"-44","gateway_ids":{"eui":"24E124FFFEF3EA11","gateway_id":"mobile-gateway"},"gps_time":"2025-07-22T16:49:21.919Z","received_at":"2025-07-22T16:49:21.935418851Z","rssi":-48,"snr":12,"time":"2025-07-22T16:49:21.919Z","timestamp":304716766,"uplink_token":"ChwKGgoObW9iaWxlLWdhdGV3YXkSCCThJP/+8+oREN63ppEBGgwIkYL/wwYQp8mj0AMgsLajlO+FASoMCJGC/8MGEMCnm7YD"}],"session_key_id":"AZgyxqYnXofMEhl8w4U48Q==","settings":{"data_rate":{"lora":{"bandwidth":125000,"coding_rate":"4/5","spreading_factor":9}},"frequency":"904500000","time":"2025-07-22T16:49:21.919Z","timestamp":304716766}}},"hardwareId":"","messageId":"7201525b-107c-4fd7-8646-200c66f357eb","protocol":"default-iot","serviceToken":"0941d3a4-8b43-490d-b90f-161b0a04c1ac","time":1753202962277},"hardwareId":"20CA0DBB117CDCB7","time":1753202962277,"messageId":"7201525b-107c-4fd7-8646-200c66f357eb","radioData":{},"signals":{"Battery":null,"Air Temperature":23,"received_time":1753202962180,"location":null,"Light":null,"Event Status":null,"Humidity":47,"DecodedPayload":null,"assistance_type":null},"protocol":"default-iot","serviceToken":"0941d3a4-8b43-490d-b90f-161b0a04c1ac","positionTime":1753202961935,"mokoSignals":null}
`

// Call function directly with []byte
// result := protocols.Parse_felix_data([]byte(jsonPayload))
// fmt.Println(string(result))


result := protocols.Nats_message_handlers([]byte(jsonPayload))
fmt.Println(string(result))
}
