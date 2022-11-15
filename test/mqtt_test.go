package test

import (
	"fmt"
	mqtt2 "github.com/eclipse/paho.mqtt.golang"
	"gogs.mirlowz.com/x/gf-x-mqtt/mqtt"
	"testing"
)

var Client *mqtt.Client

func TestRun(t *testing.T) {
	Client = mqtt.CreateClient().SetMessageCallbackFunc(func(client *mqtt.Client, client2 mqtt2.Client, message mqtt2.Message) {
		fmt.Println(message.MessageID(), message.Topic(), string(message.Payload()))
		client.SendMsg("收到", "sdt/r/1")
	})

	Client.Run()
}
