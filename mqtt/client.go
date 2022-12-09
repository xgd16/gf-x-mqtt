package mqtt

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"log"
	"os"
	"time"
)

var MqttList = map[string]*Client{}

func (t *Client) Run() {
	defer func() {
		if err := recover(); err != nil {
			g.Log().Error(gctx.New(), "MQTT服务出现恐慌:", err)
			time.Sleep(3 * time.Second)
		}

		t.Run()
	}()

	if t.Cfg.Debug {
		mqtt.DEBUG = log.New(os.Stdout, "", 0)
		mqtt.ERROR = log.New(os.Stdout, "", 0)
	}
	// 配置链接地址
	opts := mqtt.NewClientOptions().AddBroker(t.Cfg.MqttUrl).SetClientID(t.Cfg.ClientId).SetUsername(t.Cfg.Username).SetPassword(t.Cfg.Password)

	opts.SetKeepAlive(60 * time.Second)
	// 设置消息回调处理函数
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		if t.Cfg.Debug {
			fmt.Printf("TOPIC: %s\n", msg.Topic())
			fmt.Printf("MSG: %s\n", msg.Payload())
		}

		t.MessageCallbackFunc(t, client, msg)
	})
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(15 * time.Second)
	opts.SetOnConnectHandler(t.OnConnectCallBackFunc)
	opts.SetPingTimeout(1 * time.Second)
	// 创建客户端
	c := mqtt.NewClient(opts)
	// 输出启动信息
	fmt.Printf("MQTT 已启动 %s\n地址:%s\n订阅:%s\nQos:%d\n\n", t.Cfg.ClientId, t.Cfg.MqttUrl, t.Cfg.Subscribe, t.Cfg.Qos)
	// 开启链接
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	// 退出函数时断开链接
	defer c.Disconnect(250)
	// 订阅主题
	if t.Cfg.Subscribe != "false" {
		if token := c.Subscribe(t.Cfg.Subscribe, t.Cfg.Qos, nil); token.Wait() && token.Error() != nil {
			panic("订阅主题失败")
		}
	}

	// 写入客户端信息
	t.Client = &c

	MqttList[t.Cfg.Name] = t
	// 维持运行
	select {}
}

func (t *Client) SendMsg(msg any, topic string) {
	json, err := gjson.EncodeString(msg)

	if err != nil {
		g.Log().Error(gctx.New(), "mqtt 创建json出错", err.Error())
		return
	}

	token := (*t.Client).Publish(topic, t.Cfg.Qos, false, json)
	token.Wait()
}
