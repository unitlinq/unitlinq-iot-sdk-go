package unitlinq

import (
	"fmt"
	"regexp"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fxamacker/cbor/v2"
)

// Create a new wait que for each request

var defaultResponseHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	// First check that the message is from generic response topic
	RBinding, err := regexp.Compile("device/" + clientID + "/api/response/cbor$")
	if err != nil {
		//logger.Log(err.Error(), "error")
		fmt.Println(err)
	}
	if !RBinding.MatchString(msg.Topic()) {
		//logger.Log("ENERGY: Wrong topic: "+msg.Topic(), "ERROR")
		return
	}
	// The message is for the default handler
	var data genericReposne
	err = cbor.Unmarshal(msg.Payload(), &data)
	if err != nil {
		fmt.Println(err)
	}
	val, ok := waitResponse[data.RequestID]
	if ok {
		val.Data = data
		if val.Callback != nil {
			// TODO Implement callback function logic here
			_ = 0
		}
		val.Channel <- true
	}
}

func createChannel(requestID int32) chan bool {
	temp := waitResponse[requestID]
	temp.Channel = make(chan bool)
	waitResponse[requestID] = temp
	return temp.Channel
}

func deleteChannel(requestID int32) {
	delete(waitResponse, requestID)
}
