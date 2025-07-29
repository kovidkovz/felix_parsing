package models

import "encoding/json"

type IncomingData struct {
	Data         json.RawMessage `json:"data,omitempty"`
	GeoLocation  *GeoLocation    `json:"location,omitempty"`
	HardwareID   string          `json:"hardwareId,omitempty"`
	Time         interface{}     `json:"time,omitempty"`
	MessageID    string          `json:"messageId,omitempty"`
	RadioData    *RadioData      `json:"radioData,omitempty"`
	Signals      interface{}     `json:"signals,omitempty"`
	Protocol     string          `json:"protocol,omitempty"`
	ServiceToken string          `json:"serviceToken,omitempty"`
	Indoor       *Indoor         `json:"Indoor,omitempty"`
	PositionTime interface{}     `json:"positionTime,omitempty"`
	MokoSignals  interface{}     `json:"mokoSignals,omitempty"`
}

type RadioData struct {
	WifiAccessPoints []map[string]interface{} `json:"wifiAccessPoints,omitempty"`
	BluetoothBeacons []map[string]interface{} `json:"bluetoothBeacons,omitempty"`
	CellTowers       []map[string]interface{} `json:"cellTowers,omitempty"`
	Network          map[string]interface{}   `json:"network,omitempty"`
	GPS              *GPSData                 `json:"gps,omitempty"`
	Indoor           interface{}              `json:"indoor,omitempty"`
}

type GPSData struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Indoor struct {
	Building          string                 `json:"building,omitempty"`
	BuildingId        int32                  `json:"buildingId,omitempty"`
	FloorIndex        int                    `json:"floorIndex"`
	FloorLabel        string                 `json:"floorLabel,omitempty"`
	BuildingModelId   int32                  `json:"buildingModelId,omitempty"`
	LocationHierarchy map[string]interface{} `json:"locationHierarchy,omitempty"`
}

type GeoLocation struct {
	Unit   string  `json:"unit,omitempty"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Source string  `json:"source,omitempty"`
	Alt    float64 `json:"alt,omitempty"`
}

type FinalData struct {
	Data         map[string]interface{} `json:"data"`
	HardwareID   string                 `json:"hardwareId"`
	Timestamp    string                 `json:"timestamp"`
	GPS          map[string]interface{} `json:"gps"`
	Location     map[string]interface{} `json:"location"`
	Signals      Signals                `json:"signals"`
	ServiceToken string                 `json:"serviceToken"`
	Protocol     string                 `json:"protocol,omitempty"`
}

type Signals struct {
	Battery        interface{} `json:"Battery"`
	Temperature    interface{} `json:"Air Temperature"`
	ReceivedTime   interface{} `json:"received_time"`
	Location       interface{} `json:"location"`
	Light          interface{} `json:"Light"`
	EventStatus    interface{} `json:"Event Status"`
	Humidity       interface{} `json:"Humidity"`
	DecodedPayload interface{} `json:"DecodedPayload"`
	AssistanceType interface{} `json:"assistance_type"`
}
