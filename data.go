package unitlinq

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/gofrs/uuid"
)

// Push Energy readings to the unitlinq platfrom.
func (c *Client) PushEnergyData(data EnergyStruct) Token {
	var datapoint pushEnergyStruct
	var token Token
	datapoint.V1 = data.V1
	datapoint.V2 = data.V2
	datapoint.V3 = data.V3
	datapoint.V12 = data.V12
	datapoint.V23 = data.V23
	datapoint.V31 = data.V31
	datapoint.I1 = data.I1
	datapoint.I2 = data.I2
	datapoint.I3 = data.I3
	datapoint.Pf1 = data.Pf1
	datapoint.Pf2 = data.Pf2
	datapoint.Pf3 = data.Pf3
	datapoint.W1 = data.W1
	datapoint.W2 = data.W2
	datapoint.W3 = data.W3
	datapoint.ImportkWh = data.ImportkWh
	datapoint.Exportkwh = data.Exportkwh
	datapoint.ImportkVArh = data.ImportkVArh
	datapoint.ExportkVArh = data.ExportkVArh
	datapoint.Freq = data.Freq
	datapoint.Timestamp = data.Timestamp
	devID := c.ClientID.Bytes()
	datapoint.DeviceID = devID[:]
	encoded, err := cbor.Marshal(datapoint)
	if err != nil {
		fmt.Println(err.Error())
	}
	token = c.MQTT.Publish("device/data/energy/"+c.ClientID.String(), 2, true, encoded)
	return token
}

func (c *Client) PushFloatNP(datapoint NodeParamFloat) Token {
	devID := c.ClientID.Bytes()
	temp := pushNodeParamFloat{
		DeviceID:  devID[:],
		ParamID:   datapoint.ParamID,
		Value:     datapoint.Value,
		Timestamp: datapoint.Timestamp,
	}
	token := c.MQTT.Publish("device/data/param/"+c.ClientID.String(), 2, true, temp)
	return token
}

// Act as a gateway, Send energy data on behalf of other node. Please note that to act as a gateway for other device, both devices must be registered
// with the unitlinq platform
func (c *Client) EnergyGateWay(data EnergyStruct, node uuid.UUID) Token {
	var datapoint pushEnergyStruct
	var token Token
	datapoint.V1 = data.V1
	datapoint.V2 = data.V2
	datapoint.V3 = data.V3
	datapoint.V12 = data.V12
	datapoint.V23 = data.V23
	datapoint.V31 = data.V31
	datapoint.I1 = data.I1
	datapoint.I2 = data.I2
	datapoint.I3 = data.I3
	datapoint.Pf1 = data.Pf1
	datapoint.Pf2 = data.Pf2
	datapoint.Pf3 = data.Pf3
	datapoint.W1 = data.W1
	datapoint.W2 = data.W2
	datapoint.W3 = data.W3
	datapoint.ImportkWh = data.ImportkWh
	datapoint.Exportkwh = data.Exportkwh
	datapoint.ImportkVArh = data.ImportkVArh
	datapoint.ExportkVArh = data.ExportkVArh
	datapoint.Freq = data.Freq
	datapoint.Timestamp = data.Timestamp
	devID := node.Bytes()
	datapoint.DeviceID = devID[:]
	encoded, err := cbor.Marshal(datapoint)
	if err != nil {
		fmt.Println(err.Error())
	}
	token = c.MQTT.Publish("device/"+c.GetClientID()+"/gateway/energy", 2, true, encoded)
	return token
}

// Act as a gateway, send parameter data on behalf of other node. Please note that to act as a gateway for other device, both devices must be registered
// with the unitlinq platform
func (c *Client) NPFloatGateWay(data NodeParamFloat, node uuid.UUID) Token {
	devID := node.Bytes()
	temp := pushNodeParamFloat{
		DeviceID:  devID[:],
		ParamID:   data.ParamID,
		Value:     data.Value,
		Timestamp: data.Timestamp,
	}
	encoded, err := cbor.Marshal(temp)
	if err != nil {
		fmt.Println(err.Error())
	}
	token := c.MQTT.Publish("device/"+c.GetClientID()+"/gateway/param", 2, true, encoded)
	return token
}
