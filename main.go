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
        "correlation_ids": [
          "gs:uplink:01K1FW9ZRC1S4720YKZACDEBJ0"
        ],
        "end_device_ids": {
          "application_ids": {
            "application_id": "fms-demo-smart-office"
          },
          "dev_addr": "27FD4835",
          "dev_eui": "2CF7F1213050002C",
          "device_id": "sensecap-th-01-01",
          "join_eui": "8000000000000009"
        },
        "received_at": "2025-07-31T09:16:42.847503112Z",
        "uplink_message": {
          "confirmed": true,
          "consumed_airtime": "0.077056s",
          "decoded_payload": {
            "err": 0,
            "messages": [
              {
                "battery": 100,
                "type": "upload_battery"
              },
              {
                "interval": 1800,
                "type": "upload_interval"
              },
              {
                "measurementId": 4097,
                "measurementValue": 25.2,
                "type": "report_telemetry"
              },
              {
                "measurementId": 4098,
                "measurementValue": 66.6,
                "type": "report_telemetry"
              }
            ],
            "payload": "00070064001E00010110706200000102102804010014BB",
            "valid": true
          },
          "f_cnt": 12251,
          "f_port": 2,
          "frm_payload": "AAcAZAAeAAEBEHBiAAABAhAoBAEAFLs=",
          "last_battery_percentage": {
            "f_cnt": 12225,
            "received_at": "2025-07-30T20:13:30.785113840Z",
            "value": 39.130436
          },
          "locations": {
            "user": {
              "altitude": 191,
              "latitude": 30.3390715684485,
              "longitude": -97.7663040161133,
              "source": "SOURCE_REGISTRY"
            }
          },
          "network_ids": {
            "cluster_address": "nam1.cloud.thethings.industries",
            "cluster_id": "nam1",
            "net_id": "000013",
            "ns_id": "EC656E0000102C53",
            "tenant_address": "fmsiotcloud.nam1.cloud.thethings.industries",
            "tenant_id": "fmsiotcloud"
          },
          "received_at": "2025-07-31T09:16:42.636982652Z",
          "rx_metadata": [
            {
              "channel_index": 3,
              "channel_rssi": -65,
              "gateway_ids": {
                "eui": "A840411EBD004150",
                "gateway_id": "treehouse-ttn"
              },
              "received_at": "2025-07-31T09:16:42.603350498Z",
              "rssi": -65,
              "snr": 9,
              "time": "2025-07-31T09:16:42.595838Z",
              "timestamp": 3720559371,
              "uplink_token": "ChsKGQoNdHJlZWhvdXNlLXR0bhIIqEBBHr0AQVAQi/aM7g0aDAj66KzEBhCYr6avAiD4xfqUpL45"
            }
          ],
          "session_key_id": "AZHQDenGwat+GpjajlWCCA==",
          "settings": {
            "data_rate": {
              "lora": {
                "bandwidth": 125000,
                "coding_rate": "4/5",
                "spreading_factor": 7
              }
            },
            "frequency": "904500000",
            "time": "2025-07-31T09:16:42.595838Z",
            "timestamp": 3720559371
          },
          "version_ids": {
            "band_id": "US_902_928",
            "brand_id": "sensecap",
            "firmware_version": "3.4",
            "hardware_version": "2.0",
            "model_id": "sensecap-air-th"
          }
        }
      },
      "hardwareId": "",
      "messageId": "026fd429-55d6-497d-90ac-f470d195e74e",
      "protocol": "default-iot",
      "serviceToken": "0941d3a4-8b43-490d-b90f-161b0a04c1ac",
      "time": 1753953402935
    }`

// // Call function directly with []byte
result := protocols.Parse_felix_data([]byte(jsonPayload))
fmt.Println(string(result))


// result := protocols.Nats_message_handlers([]byte(jsonPayload))
// fmt.Println(string(result))

// result, err := FallbackLookup([]byte(jsonPayload))
// if err != nil {
//   fmt.Println("error", err)
// }
// fmt.Println(result)
}
