package test

import (
	"fmt"
	"github.com/xgd16/gf-x-mqtt/xmqtt"
	"testing"
)

func TestRun(t *testing.T) {
	xmqtt.CreateClient(func(option *xmqtt.ClientCallBackOption, config *xmqtt.Config) {
		option.MessageCallbackFunc = func(data *xmqtt.MessageHandlerData) {
			fmt.Println(data.GetMessageId(), data.GetTopic(), data.GetMsg())
			//client.SendMsg("收到", "sdt/c/1")
		}
	})
	select {}
}
