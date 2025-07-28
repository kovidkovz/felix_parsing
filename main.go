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
      "gs:uplink:01K10YG9TRMPYD0ET9FF03TSA0"
    ],
    "end_device_ids": {
      "application_ids": {
        "application_id": "fms-demo-smart-office"
      },
      "dev_addr": "27FD4838",
      "dev_eui": "C4F58AFFFF9216CA",
      "device_id": "00-sc-m010-16ca-20-a",
      "join_eui": "70B3D57ED0026B87"
    },
    "received_at": "2025-07-25T14:06:30.443341246Z",
    "uplink_message": {
      "consumed_airtime": "0.061696s",
      "decoded_payload": {
        "battery_percent": {
          "battery": 78
        },
        "charging_status": "No charging",
        "event_type": "Light intensity over threshold",
        "hex_format_payload": "6368838f66f60e18000c",
        "light_intensity": 12,
        "port": 5,
        "temperature": "24°C",
        "time": "07/25/2025, 14:06:30",
        "timezone": -10
      },
      "f_cnt": 521,
      "f_port": 5,
      "frm_payload": "Y2iDj2b2DhgADA==",
      "network_ids": {
        "cluster_address": "nam1.cloud.thethings.industries",
        "cluster_id": "nam1",
        "net_id": "000013",
        "ns_id": "EC656E0000102C53",
        "tenant_address": "fmsiotcloud.nam1.cloud.thethings.industries",
        "tenant_id": "fmsiotcloud"
      },
      "received_at": "2025-07-25T14:06:30.232659552Z",
      "rx_metadata": [
        {
          "channel_index": 3,
          "channel_rssi": -53,
          "frequency_offset": "-679",
          "gateway_ids": {
            "eui": "24E124FFFEF2F8C3",
            "gateway_id": "fms-iot-things"
          },
          "gps_time": "2025-07-25T14:06:30.196Z",
          "location": {
            "altitude": 187,
            "latitude": 30.2263929475762,
            "longitude": -97.7599310874939,
            "source": "SOURCE_REGISTRY"
          },
          "received_at": "2025-07-25T14:06:30.207247711Z",
          "rssi": -53,
          "snr": 14.5,
          "time": "2025-07-25T14:06:30.196Z",
          "timestamp": 1535621151,
          "uplink_token": "ChwKGgoOZm1zLWlvdC10aGluZ3MSCCThJP/+8vjDEJ/wntwFGgsI5p6OxAYQmqfMbiCY8pzR2MMiKgsI5p6OxAYQgPK6XQ=="
        }
      ],
      "session_key_id": "AZg+jOkGDOl/MJPZNynsCA==",
      "settings": {
        "data_rate": {
          "lora": {
            "bandwidth": 125000,
            "coding_rate": "4/5",
            "spreading_factor": 7
          }
        },
        "frequency": "904500000",
        "time": "2025-07-25T14:06:30.196Z",
        "timestamp": 1535621151
      }
    }
  },
  "hardwareId": "",
  "messageId": "2dca370f-fe38-457b-8f54-eaa927bb29f2",
  "protocol": "default-iot",
  "serviceToken": "0941d3a4-8b43-490d-b90f-161b0a04c1ac",
  "time": 1753452390539
}

`

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
