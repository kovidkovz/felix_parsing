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
{"data":{"data":{"end_device_ids":{"device_id":"sensecap5-009e","application_ids":{"application_id":"fms-demo-smart-office"},"dev_eui":"2CF7F1C07030009E","join_eui":"A194E17CADD05677","dev_addr":"27FD4837"},"correlation_ids":["gs:uplink:01K10NDWANAABCGET8K1RSEC57"],"received_at":"2025-07-25T11:27:53.891925351Z","uplink_message":{"session_key_id":"AZg5IFTD2ryh5jQbe+asIg==","f_port":5,"f_cnt":3568,"frm_payload":"CAAAAABog2owqEadKthrwFRsDiqApsCoRp0rFXi8ARcAAFc=","decoded_payload":{"err":0,"messages":[[{"measurementId":"4200","measurementValue":[],"motionId":0,"timestamp":1753442864000,"type":"Event Status"},{"measurementId":"5002","measurementValue":[{"mac":"A8:46:9D:2A:D8:6B","rssi":"-64"},{"mac":"54:6C:0E:2A:80:A6","rssi":"-64"},{"mac":"A8:46:9D:2B:15:78","rssi":"-68"}],"motionId":0,"timestamp":1753442864000,"type":"BLE Scan"},{"measurementId":"4097","measurementValue":27.9,"motionId":0,"timestamp":1753442864000,"type":"Air Temperature"},{"measurementId":"4199","measurementValue":0,"motionId":0,"timestamp":1753442864000,"type":"Light"},{"measurementId":"3000","measurementValue":87,"motionId":0,"timestamp":1753442864000,"type":"Battery"}]],"payload":"080000000068836a30a8469d2ad86bc0546c0e2a80a6c0a8469d2b1578bc0117000057","valid":true},"rx_metadata":[{"gateway_ids":{"gateway_id":"fms-iot-things","eui":"24E124FFFEF2F8C3"},"time":"2025-07-25T11:27:53.655Z","timestamp":609021333,"rssi":-60,"channel_rssi":-60,"snr":10.8,"frequency_offset":"91","location":{"latitude":30.2263929475762,"longitude":-97.7599310874939,"altitude":187,"source":"SOURCE_REGISTRY"},"uplink_token":"ChwKGgoOZm1zLWlvdC10aGluZ3MSCCThJP/+8vjDEJXbs6ICGgwIudSNxAYQiYnaxgIgiJyA5NyuICoMCLnUjcQGEMCDqrgC","channel_index":2,"gps_time":"2025-07-25T11:27:53.655Z","received_at":"2025-07-25T11:27:53.662045357Z"}],"settings":{"data_rate":{"lora":{"bandwidth":125000,"spreading_factor":7,"coding_rate":"4/5"}},"frequency":"904300000","timestamp":609021333,"time":"2025-07-25T11:27:53.655Z"},"received_at":"2025-07-25T11:27:53.686057176Z","confirmed":true,"consumed_airtime":"0.097536s","version_ids":{"brand_id":"sensecap","model_id":"sensecapt1000-tracker-ab","hardware_version":"1.0","firmware_version":"1.0","band_id":"US_902_928"},"network_ids":{"net_id":"000013","ns_id":"EC656E0000102C53","tenant_id":"fmsiotcloud","cluster_id":"nam1","cluster_address":"nam1.cloud.thethings.industries","tenant_address":"fmsiotcloud.nam1.cloud.thethings.industries"},"last_battery_percentage":{"f_cnt":3417,"value":100,"received_at":"2025-07-25T08:57:47.685978822Z"}}},"hardwareId":"","time":1753442873992,"messageId":"af0d25a5-bff8-46ff-b84e-b7ae0ee50658","serviceToken":"0941d3a4-8b43-490d-b90f-161b0a04c1ac","protocol":"default-iot"},"hardwareId":"2CF7F1C07030009E","time":1753442873992,"messageId":"af0d25a5-bff8-46ff-b84e-b7ae0ee50658","radioData":{"bluetoothBeacons":[{"macAddress":"A8:46:9D:2A:D8:6B","signalStrength":"-64"},{"macAddress":"54:6C:0E:2A:80:A6","signalStrength":"-64"},{"macAddress":"A8:46:9D:2B:15:78","signalStrength":"-68"}]},"signals":{"batteryLevel":87,"device_name":"sensecap5-009e","err":0,"eventStatus":[],"light":0,"payload":"080000000068836a30a8469d2ad86bc0546c0e2a80a6c0a8469d2b1578bc0117000057","received_time":1753442873891,"temperatureLevel":27.9,"valid":true},"protocol":"default-iot","serviceToken":"0941d3a4-8b43-490d-b90f-161b0a04c1ac","positionTime":1753442873662}


`

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
