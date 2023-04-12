package unitlinq

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/fxamacker/cbor/v2"
)

func (c *Client) CheckRemoteConfig() (ConfigMetaDataResp, bool, error) {
	// Check metadata from the server
	// Generate a random int32 number
	to := false
	temp := configReq{
		RequestID: rand.Int31(),
	}
	if !responseSubscribed {
		return ConfigMetaDataResp{}, false, errors.New("please initialize and connect the client first")
	}
	// Response topic is subscribed, before pushing the request register the request with the response handler
	response := createChannel(temp.RequestID)
	payload, err := cbor.Marshal(temp)
	if err != nil {
		return ConfigMetaDataResp{}, false, err
	}
	c.MQTT.Publish("device/"+c.GetClientID()+"/api/configmetadata/cbor", 2, false, payload)
	select {
	case <-response:
	case <-time.After(c.timeout):
		to = true
	}
	if to {
		return ConfigMetaDataResp{}, false, errors.New("request timeout")
	}
	resp := waitResponse[temp.RequestID]
	// Parse the response
	var data ConfigMetaDataResp
	err = cbor.Unmarshal(resp.Data.Payload, &data)
	if err != nil {
		fmt.Println(err)
	}
	if !resp.Data.Success {
		return data, resp.Data.Success, errors.New(resp.Data.Error)
	}
	return data, resp.Data.Success, nil
}

func (c *Client) GetRemoteConfig() (ConfigResp, bool, error) {
	to := false
	temp := configReq{
		RequestID: rand.Int31(),
	}
	if !responseSubscribed {
		return ConfigResp{}, false, errors.New("please initialize and connect the client first")
	}
	// Response topic is subscribed, before pushing the request register the request with the response handler
	response := createChannel(temp.RequestID)
	payload, err := cbor.Marshal(temp)
	if err != nil {
		return ConfigResp{}, false, err
	}
	c.MQTT.Publish("device/"+c.GetClientID()+"/api/getconfig/cbor", 2, false, payload)
	select {
	case <-response:
	case <-time.After(c.timeout):
		to = true
	}
	if to {
		return ConfigResp{}, false, errors.New("request timeout")
	}
	resp := waitResponse[temp.RequestID]
	// Parse the response
	var data ConfigResp
	err = cbor.Unmarshal(resp.Data.Payload, &data)
	if err != nil {
		fmt.Println(err)
	}
	if !resp.Data.Success {
		return data, resp.Data.Success, errors.New(resp.Data.Error)
	}
	return data, resp.Data.Success, nil
}
