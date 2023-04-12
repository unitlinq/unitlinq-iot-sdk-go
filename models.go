package unitlinq

import mqtt "github.com/eclipse/paho.mqtt.golang"

type Subscripiton struct {
	Topic    string
	Qos      byte
	Callback mqtt.MessageHandler
}

var subscriptionPool []Subscripiton
var responseSubscribed bool
var clientID string

type EnergyStruct struct {
	V1          float32
	V2          float32
	V3          float32
	V12         float32
	V23         float32
	V31         float32
	I1          float32
	I2          float32
	I3          float32
	Pf1         float32
	Pf2         float32
	Pf3         float32
	W1          float32
	W2          float32
	W3          float32
	ImportkWh   float32
	Exportkwh   float32
	ImportkVArh float32
	ExportkVArh float32
	Freq        float32
	Timestamp   int64
}

type pushEnergyStruct struct {
	DeviceID    []byte  `cbor:"nodeid"`
	V1          float32 `cbor:"v1"`
	V2          float32 `cbor:"v2"`
	V3          float32 `cbor:"v3"`
	V12         float32 `cbor:"v12"`
	V23         float32 `cbor:"v23"`
	V31         float32 `cbor:"v31"`
	I1          float32 `cbor:"i1"`
	I2          float32 `cbor:"i2"`
	I3          float32 `cbor:"i3"`
	Pf1         float32 `cbor:"pf1"`
	Pf2         float32 `cbor:"pf2"`
	Pf3         float32 `cbor:"pf3"`
	W1          float32 `cbor:"w1"`
	W2          float32 `cbor:"w2"`
	W3          float32 `cbor:"w3"`
	ImportkWh   float32 `cbor:"ikwh"`
	Exportkwh   float32 `cbor:"ekwh"`
	ImportkVArh float32 `cbor:"ikvarh"`
	ExportkVArh float32 `cbor:"ekvarh"`
	Freq        float32 `cbor:"freq"`
	Timestamp   int64   `cbor:"ts"`
}

type NodeParamFloat struct {
	ParamID   []byte
	Value     float32
	Timestamp int64
}

type genericReposne struct {
	RequestID int32  `cbor:"reqid"`
	Success   bool   `cbor:"success"`
	Error     string `cbor:"error"`
	ErrCode   int32  `cbor:"errcode"`
	Payload   []byte `cbor:"data"`
}

type pushNodeParamFloat struct {
	DeviceID  []byte  `cbor:"nid"`
	ParamID   []byte  `cbor:"pid"`
	Value     float32 `cbor:"v"`
	Timestamp int64   `cbor:"ts"`
}

type configReq struct {
	RequestID int32 `cbor:"reqid"`
}

type ConfigMetaDataResp struct {
	UpdatedOn int32 `cbor:"updatedon"`
}

type ConfigResp struct {
	Data []byte `cbor:"data"`
}

type ResponseCallback func(requestID int, payload []byte)

type responseObject struct {
	Channel  chan bool
	Callback ResponseCallback
	Data     genericReposne
}

var waitResponse = make(map[int32]responseObject)
