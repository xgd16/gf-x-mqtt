package mqtt

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"log"
	"os"
	"time"
)

func (t *Client) Run() {
	defer func() {
		if err := recover(); err != nil {
			g.Log().Error(gctx.New(), "MQTT服务出现恐慌:", err)
			time.Sleep(3 * time.Second)
			t.Run()
		}
	}()

	if t.Cfg.Debug {
		mqtt.DEBUG = log.New(os.Stdout, "", 0)
		mqtt.ERROR = log.New(os.Stdout, "", 0)
	}
	// 配置链接地址
	opts := mqtt.NewClientOptions().AddBroker(t.Cfg.MqttUrl).SetClientID(t.Cfg.ClientId)

	opts.SetKeepAlive(60 * time.Second)
	// 设置消息回调处理函数
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		if t.Cfg.Debug {
			fmt.Printf("TOPIC: %s\n", msg.Topic())
			fmt.Printf("MSG: %s\n", msg.Payload())
		}

		t.MessageCallbackFunc(client, msg)
	})
	opts.SetOnConnectHandler(t.OnConnectCallBackFunc)
	opts.SetPingTimeout(1 * time.Second)
	// 创建客户端
	c := mqtt.NewClient(opts)
	// 开启链接
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	// 退出函数时断开链接
	defer c.Disconnect(250)
	// 订阅主题
	if token := c.Subscribe(t.Cfg.Subscribe, t.Cfg.Qos, nil); token.Wait() && token.Error() != nil {
		panic("订阅主题失败")
	}
	// 写入客户端信息
	t.Client = &c
	// 维持运行
	select {}
	//// 发布消息
	//token := c.Publish("testtopic/1", 0, false, "Hello World")
	//token.Wait(
	//
	//time.Sleep(200 * time.Second)
	//
	//// 取消订阅
	//if token := c.Unsubscribe("testtopic/#"); token.Wait() && token.Error() != nil {
	//	fmt.Println(token.Error())
	//	os.Exit(1)
	//}
}

func (t *Client) SendMsg(msg string) {
	token := (*t.Client).Publish(t.Cfg.Topic, t.Cfg.Qos, false, msg)
	token.Wait()
}
