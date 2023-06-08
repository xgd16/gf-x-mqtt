package test

import (
	"fmt"
	mqtt2 "github.com/eclipse/paho.mqtt.golang"
	"github.com/xgd16/gf-x-mqtt/mqtt"
	"testing"
)

func TestRun(t *testing.T) {
	mqtt.CreateClient(func(option *mqtt.ClientOption, config *mqtt.Config) {

		option.MessageCallbackFunc = func(client *mqtt.Client, client2 mqtt2.Client, message mqtt2.Message) {
			fmt.Println(message.MessageID(), message.Topic(), string(message.Payload()))
			client.SendMsg("收到", "sdt/c/1")
		}

	})

	select {}
}
