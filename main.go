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
    "data": {
      "correlation_ids": [
        "gs:uplink:01K0Y4CA2KHRFHV9113D61QW6W"
      ],
      "end_device_ids": {
        "application_ids": {
          "application_id": "fms-demo-smart-office"
        },
        "dev_addr": "27FD483F",
        "dev_eui": "647FDA0000018263",
        "device_id": "vivid-8263",
        "join_eui": "647FDA8010000000"
      },
      "received_at": "2025-07-24T11:51:27.778023975Z",
      "uplink_message": {
        "consumed_airtime": "0.028288s",
        "decoded_payload": {
          "ambient_temperature": 25.3,
          "battery_voltage": 3.026,
          "port": "10",
          "raw": "[03, 67, 00, FD, 04, 68, 59, 00, BA, 0B, D2]",
          "relative_humidity": 44.5
        },
        "f_cnt": 1547,
        "f_port": 10,
        "frm_payload": "A2cA/QRoWQC6C9I=",
        "last_battery_percentage": {
          "f_cnt": 1495,
          "received_at": "2025-07-23T17:51:24.210628973Z",
          "value": 99.20949
        },
        "network_ids": {
          "cluster_address": "nam1.cloud.thethings.industries",
          "cluster_id": "nam1",
          "net_id": "000013",
          "ns_id": "EC656E0000102C53",
          "tenant_address": "fmsiotcloud.nam1.cloud.thethings.industries",
          "tenant_id": "fmsiotcloud"
        },
        "normalized_payload": [
          {
            "air": {
              "relativeHumidity": 44.5,
              "temperature": 25.3
            },
            "battery": 3.026
          }
        ],
        "received_at": "2025-07-24T11:51:27.572703341Z",
        "rx_metadata": [
          {
            "channel_index": 8,
            "channel_rssi": -52,
            "frequency_offset": "-5050",
            "gateway_ids": {
              "eui": "24E124FFFEF2F8C3",
              "gateway_id": "fms-iot-things"
            },
            "gps_time": "2025-07-24T11:51:27.535Z",
            "location": {
              "altitude": 187,
              "latitude": 30.2263929475762,
              "longitude": -97.7599310874939,
              "source": "SOURCE_REGISTRY"
            },
            "received_at": "2025-07-24T11:51:27.548670480Z",
            "rssi": -52,
            "snr": 14.2,
            "time": "2025-07-24T11:51:27.535Z",
            "timestamp": 1522298865,
            "uplink_token": "ChwKGgoOZm1zLWlvdC10aGluZ3MSCCThJP/+8vjDEPHf8dUFGgwIv7yIxAYQmerRkAIg6IrVgKeFDSoMCL+8iMQGEMDnjf8B"
          }
        ],
        "session_key_id": "AZdgUOigXtnCMxmdIcruZA==",
        "settings": {
          "data_rate": {
            "lora": {
              "bandwidth": 500000,
              "coding_rate": "4/5",
              "spreading_factor": 8
            }
          },
          "frequency": "904600000",
          "time": "2025-07-24T11:51:27.535Z",
          "timestamp": 1522298865
        },
        "version_ids": {
          "band_id": "US_902_928",
          "brand_id": "tektelic",
          "firmware_version": "2.1",
          "hardware_version": "D",
          "model_id": "t00061xx-vivid"
        }
      }
    },
    "hardwareId": "",
    "messageId": "6d2fda4a-4c40-4ce5-8781-cc8412b204f8",
    "protocol": "default-iot",
    "serviceToken": "0941d3a4-8b43-490d-b90f-161b0a04c1ac",
    "time": 1753357887872
  },
  "location": {
    "lat": 30.2263929475762,
    "lng": -97.7599310874939,
    "source": "SOURCE_REGISTRY",
    "alt": 187
  },
  "hardwareId": "647FDA0000018263",
  "time": 1753357887872,
  "messageId": "6d2fda4a-4c40-4ce5-8781-cc8412b204f8",
  "signals": {
    "batteryLevel": 3.026,
    "battery_voltage": 3.026,
    "humidity": 44.5,
    "port": "10",
    "raw": "[03, 67, 00, FD, 04, 68, 59, 00, BA, 0B, D2]",
    "received_time": 1753357887778,
    "temperatureLevel": 25.3
  },
  "protocol": "default-iot",
  "serviceToken": "0941d3a4-8b43-490d-b90f-161b0a04c1ac",
  "positionTime": 1753357887548
}
`

// // Call function directly with []byte
// result := protocols.Parse_felix_data([]byte(jsonPayload))
// fmt.Println(string(result))


result := protocols.Nats_message_handlers([]byte(jsonPayload))
fmt.Println(string(result))
}
