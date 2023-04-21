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

func NewDefaultEnergyStruct() EnergyStruct {
	return EnergyStruct{
		V1:          0,
		V2:          0,
		V3:          0,
		V12:         0,
		V23:         0,
		V31:         0,
		I1:          0,
		I2:          0,
		I3:          0,
		Pf1:         1,
		Pf2:         1,
		Pf3:         1,
		W1:          0,
		W2:          0,
		W3:          0,
		ImportkWh:   0,
		Exportkwh:   0,
		ImportkVArh: 0,
		ExportkVArh: 0,
		Freq:        0,
		Timestamp:   0,
	}
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

func NewDefaultNpFloat() NodeParamFloat {
	return NodeParamFloat{
		ParamID:   nil,
		Value:     0,
		Timestamp: 0,
	}
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

const serverCertEmbedded = `
-----BEGIN CERTIFICATE-----
MIIClTCCAhugAwIBAgIIQE4JchcUCD0wCgYIKoZIzj0EAwIwdzELMAkGA1UEBhMC
SU4xEDAOBgNVBAgTB0d1amFyYXQxDzANBgNVBAcTBkluZGlhbjEVMBMGA1UEChMM
RmxpbnQgRW5lcmd5MREwDwYDVQQLEwhVbml0Z3JpZDEbMBkGA1UEAxMSVW5pdGdy
aWQgUm9vdCBDQSAxMCAXDTIyMDkwOTAwMDAwMFoYDzIwNTIwOTA4MjM1OTU5WjB3
MQswCQYDVQQGEwJJTjEQMA4GA1UECBMHR3VqYXJhdDEPMA0GA1UEBxMGSW5kaWFu
MRUwEwYDVQQKEwxGbGludCBFbmVyZ3kxETAPBgNVBAsTCFVuaXRncmlkMRswGQYD
VQQDExJVbml0Z3JpZCBSb290IENBIDEwdjAQBgcqhkjOPQIBBgUrgQQAIgNiAARm
49Ttjevys9q21Ue3bN2LdjCbgSME7T2rONjBCH6onja6U2PfS4hw9ALyjqprBcx0
/qlmGpnAuiVEm4qJWn1CyNQfYEfZPpvuV0NkBON4Nx3GEzEgthvuMtij/NROxYaj
cjBwMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFHD+22jmlpoT7can4bXC1Oic
c8BLMAsGA1UdDwQEAwIBBjARBglghkgBhvhCAQEEBAMCAAcwHgYJYIZIAYb4QgEN
BBEWD3hjYSBjZXJ0aWZpY2F0ZTAKBggqhkjOPQQDAgNoADBlAjBl5LeydupNqGzs
YHObNZqZuHEz1XRH0siSs70ZGLHe3jbJxzDZF9lJeS8pNio3S6wCMQDYkWMf2wuR
wbd/MYG7mURoAkzwlIJXHN94USIQ9KQ9hWCKTl7UyzGChPKpunDAJe0=
-----END CERTIFICATE-----`

type OnConnectCallback func(Client)
type OnConnLostCallback func(Client, error)
