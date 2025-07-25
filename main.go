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
{"data":{"data":{"end_device_ids":{"device_id":"00-sc-t1000-087e-09-a","application_ids":{"application_id":"fms-demo-smart-office"},"dev_eui":"2CF7F1C06490087E","join_eui":"25E7B83C131B45F8","dev_addr":"27FD483B"},"correlation_ids":["gs:uplink:01K0XA1KNS1FSDMGM5T2PBVGPD"],"received_at":"2025-07-24T04:11:14.189704872Z","uplink_message":{"session_key_id":"AZf6cOUr1v5a0EI6LSi5pQ==","f_port":5,"f_cnt":21192,"frm_payload":"CAAAAABogbJeVGwOKoCmxahGnSrYa72oRp0sDUC7ASkAAFM=","decoded_payload":{"err":0,"messages":[[{"measurementId":"4200","measurementValue":[],"motionId":0,"timestamp":1753330270000,"type":"Event Status"},{"measurementId":"5002","measurementValue":[{"mac":"54:6C:0E:2A:80:A6","rssi":"-59"},{"mac":"A8:46:9D:2A:D8:6B","rssi":"-67"},{"mac":"A8:46:9D:2C:0D:40","rssi":"-69"}],"motionId":0,"timestamp":1753330270000,"type":"BLE Scan"},{"measurementId":"4097","measurementValue":29.7,"motionId":0,"timestamp":1753330270000,"type":"Air Temperature"},{"measurementId":"4199","measurementValue":0,"motionId":0,"timestamp":1753330270000,"type":"Light"},{"measurementId":"3000","measurementValue":83,"motionId":0,"timestamp":1753330270000,"type":"Battery"}]],"payload":"08000000006881b25e546c0e2a80a6c5a8469d2ad86bbda8469d2c0d40bb0129000053","valid":true},"rx_metadata":[{"gateway_ids":{"gateway_id":"fms-iot-things","eui":"24E124FFFEF2F8C3"},"time":"2025-07-24T04:11:13.957Z","timestamp":3973492905,"rssi":-65,"channel_rssi":-65,"snr":14.8,"frequency_offset":"480","location":{"latitude":30.2263929475762,"longitude":-97.7599310874939,"altitude":187,"source":"SOURCE_REGISTRY"},"uplink_token":"ChwKGgoOZm1zLWlvdC10aGluZ3MSCCThJP/+8vjDEKnh2uYOGgwI4eSGxAYQkdeP0gMgqKiItdLhBioMCOHkhsQGEMDSqsgD","channel_index":5,"gps_time":"2025-07-24T04:11:13.957Z","received_at":"2025-07-24T04:11:13.954226734Z"}],"settings":{"data_rate":{"lora":{"bandwidth":125000,"spreading_factor":7,"coding_rate":"4/5"}},"frequency":"904900000","timestamp":3973492905,"time":"2025-07-24T04:11:13.957Z"},"received_at":"2025-07-24T04:11:13.978241341Z","confirmed":true,"consumed_airtime":"0.097536s","version_ids":{"brand_id":"sensecap","model_id":"sensecapt1000-tracker-ab","hardware_version":"1.0","firmware_version":"1.0","band_id":"US_902_928"},"network_ids":{"net_id":"000013","ns_id":"EC656E0000102C53","tenant_id":"fmsiotcloud","cluster_id":"nam1","cluster_address":"nam1.cloud.thethings.industries","tenant_address":"fmsiotcloud.nam1.cloud.thethings.industries"},"last_battery_percentage":{"f_cnt":21134,"value":100,"received_at":"2025-07-24T03:42:20.370961182Z"}}},"hardwareId":"","time":1753330274279,"messageId":"39109457-4015-4af3-b443-74383a962673","serviceToken":"0941d3a4-8b43-490d-b90f-161b0a04c1ac","protocol":"default-iot"},"hardwareId":"2CF7F1C06490087E","time":1753330274279,"messageId":"39109457-4015-4af3-b443-74383a962673","radioData":{"bluetoothBeacons":[{"macAddress":"54:6C:0E:2A:80:A6","signalStrength":"-59"},{"macAddress":"A8:46:9D:2A:D8:6B","signalStrength":"-67"},{"macAddress":"A8:46:9D:2C:0D:40","signalStrength":"-69"}]},"signals":{"batteryLevel":83,"err":0,"eventStatus":[],"light":0,"payload":"08000000006881b25e546c0e2a80a6c5a8469d2ad86bbda8469d2c0d40bb0129000053","received_time":1753330274189,"temperatureLevel":29.7,"valid":true},"protocol":"default-iot","serviceToken":"0941d3a4-8b43-490d-b90f-161b0a04c1ac","positionTime":1753330273954}
`

// // Call function directly with []byte
// result := protocols.Parse_felix_data([]byte(jsonPayload))
// fmt.Println(string(result))


result := protocols.Nats_message_handlers([]byte(jsonPayload))
fmt.Println(string(result))
}
