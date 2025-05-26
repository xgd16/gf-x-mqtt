package xmqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	ConnectEvent    = "ConnectEvent"
	DisconnectEvent = "DisconnectEvent"
	NullEvent       = "NullEvent"
)

// MqttList MQTT 客户端列表
var MqttList = CreateSafeMQTTList()

func (t *Client) Run(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			g.Log().Error(ctx, "MQTT服务出现恐慌:", err)
			time.Sleep(3 * time.Second)
		}
		t.Run(ctx)
	}()
	// 设置 debug
	if t.Cfg.Debug {
		mqtt.DEBUG = log.New(os.Stdout, "", 0)
		mqtt.ERROR = log.New(os.Stdout, "", 0)
	}
	// 配置链接地址
	opts := mqtt.NewClientOptions().AddBroker(t.Cfg.MqttUrl).SetClientID(t.Cfg.ClientId).SetUsername(t.Cfg.Username).SetPassword(t.Cfg.Password)

	opts.SetKeepAlive(time.Duration(t.Cfg.Ping) * time.Second)
	// 设置消息回调处理函数
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		if t.Cfg.Debug {
			fmt.Printf("订阅: %s\n消息: %s\n", msg.Topic(), msg.Payload())
		}
		t.MessageCallbackFunc(&MessageHandlerData{
			XMQTT:   t,
			OMQTT:   client,
			Message: msg,
		})
	})
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(15 * time.Second)
	opts.SetOnConnectHandler(t.OnConnectCallBackFunc)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetCleanSession(t.Cfg.CleanSession)
	// 创建客户端
	c := mqtt.NewClient(opts)
	// 输出启动信息
	fmt.Printf(`
MQTT 已启动
ClientId: %s
地址: %s
订阅: %s
Qos: %d
Ping: %d

`, t.Cfg.ClientId, t.Cfg.MqttUrl, t.Cfg.Subscribe, t.Cfg.Qos, t.Cfg.Ping)
	// 开启链接
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	// 退出函数时断开链接
	defer func() {
		g.Log().Error(ctx, t.Cfg.Name, "连接退出...")
		// 关闭连接
		c.Disconnect(250)
	}()
	// 订阅主题
	if t.Cfg.Subscribe != "false" {
		for _, item := range gstr.Split(t.Cfg.Subscribe, ",") {
			// 订阅
			subInfo := gstr.Split(item, ":")
			if token := c.Subscribe(subInfo[0], gconv.Byte(subInfo[1]), nil); token.Wait() && token.Error() != nil {
				panic("订阅主题失败")
			}
		}
	}
	// 写入客户端信息
	t.Client = &c
	// 写入全局
	MqttList.Set(t.Cfg.Name, t)
	// 维持运行
	select {
	case <-ctx.Done():
		return
	}
}

func (t *Client) SendMsg(msg any, topic string, qos ...byte) (err error) {
	// 如果没有初始化那么退出发送
	if !t.IsInit {
		err = fmt.Errorf("MQTT未初始化")
		return
	}
	// 设置 qos
	var qosNumber = t.Cfg.Qos
	// 如果在配置文件中配置了 那么使用配置文件中的配置
	if len(qos) >= 1 {
		qosNumber = qos[0]
	}
	// 将传入消息解析为json
	jsonByte, err := json.Marshal(msg)
	if err != nil {
		return
	}
	// 推送消息
	if token := (*t.Client).Publish(topic, qosNumber, false, string(jsonByte)); token.Wait() && token.Error() != nil {
		err = token.Error()
		return
	}
	return
}

func (t *Client) Json() *MqttResp {
	return CreateMqttResp(t)
}
