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
{"data":{"data":{"correlation_ids":["gs:uplink:01K0XRE1Z2A73R6BVWEAZ1NW7B"],"end_device_ids":{"application_ids":{"application_id":"fms-demo-smart-office"},"dev_addr":"27FD4836","dev_eui":"34814D249BFDFBB4","device_id":"thpo-fbb4","join_eui":"EEA2AC85876DC690"},"received_at":"2025-07-24T08:22:42.103615667Z","uplink_message":{"confirmed":true,"consumed_airtime":"0.370688s","decoded_payload":{"humidity":44.8,"motion_detected":false,"pressure":1890,"raw_payload":"01 03 01 c0 00 63 00 00","status":{"humidity_valid":true,"pressure_valid":false,"temperature_valid":true},"temperature":{"celsius":25.9,"fahrenheit":78.62}},"f_cnt":46885,"f_port":1,"frm_payload":"AQMBwABjAAA=","network_ids":{"cluster_address":"nam1.cloud.thethings.industries","cluster_id":"nam1","net_id":"000013","ns_id":"EC656E0000102C53","tenant_address":"fmsiotcloud.nam1.cloud.thethings.industries","tenant_id":"fmsiotcloud"},"packet_error_rate":0.04761905,"received_at":"2025-07-24T08:22:41.890670926Z","rx_metadata":[{"channel_index":4,"channel_rssi":-98,"frequency_offset":"-280","gateway_ids":{"eui":"24E124FFFEF2F8C3","gateway_id":"fms-iot-things"},"gps_time":"2025-07-24T08:22:41.845Z","location":{"altitude":187,"latitude":30.2263929475762,"longitude":-97.7599310874939,"source":"SOURCE_REGISTRY"},"received_at":"2025-07-24T08:22:41.866968142Z","rssi":-98,"snr":4,"time":"2025-07-24T08:22:41.845Z","timestamp":1881523211,"uplink_token":"ChwKGgoOZm1zLWlvdC10aGluZ3MSCCThJP/+8vjDEIuIl4EHGgwI0dqHxAYQgtqzqAMg+JWXnOGYCioMCNHah8QGEMDa9pID"}],"session_key_id":"AZbckHRsJYY9pAVVWzShqA==","settings":{"data_rate":{"lora":{"bandwidth":125000,"coding_rate":"4/5","spreading_factor":10}},"frequency":"904700000","time":"2025-07-24T08:22:41.845Z","timestamp":1881523211}}},"hardwareId":"","messageId":"845048e4-ae41-4976-9904-d7a6aa4a4752","protocol":"default-iot","serviceToken":"0941d3a4-8b43-490d-b90f-161b0a04c1ac","time":1753345362204},"location":{"lat":30.2263929475762,"lng":-97.7599310874939,"source":"SOURCE_REGISTRY","alt":187},"hardwareId":"34814D249BFDFBB4","time":1753345362204,"messageId":"845048e4-ae41-4976-9904-d7a6aa4a4752","signals":{"humidity":44.8,"motion_detected":false,"pressure":1890,"raw_payload":"01 03 01 c0 00 63 00 00","received_time":1753345362103,"status":{"humidity_valid":true,"pressure_valid":false,"temperature_valid":true},"temperatureLevel":25.9},"protocol":"default-iot","serviceToken":"0941d3a4-8b43-490d-b90f-161b0a04c1ac","positionTime":1753345361866}
`

// // Call function directly with []byte
// result := protocols.Parse_felix_data([]byte(jsonPayload))
// fmt.Println(string(result))


result := protocols.Nats_message_handlers([]byte(jsonPayload))
fmt.Println(string(result))
}
