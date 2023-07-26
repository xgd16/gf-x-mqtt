package xmqtt

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

const (
	ConnectEvent    = "ConnectEvent"
	DisconnectEvent = "DisconnectEvent"
	NullEvent       = "NullEvent"
)

// MqttList MQTT 客户端列表
var MqttList = CreateSafeMQTTList()

var SendDefaultQos byte = 0

func (t *Client) Run() {
	defer func() {
		if err := recover(); err != nil {
			g.Log().Error(gctx.New(), "MQTT服务出现恐慌:", err)
			time.Sleep(3 * time.Second)
		}
		t.Run()
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
			fmt.Printf("TOPIC: %s\n", msg.Topic())
			fmt.Printf("MSG: %s\n", msg.Payload())
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
	fmt.Printf("MQTT 已启动 %s\n地址:%s\n订阅:%s\nQos:%d\nPing:%d\n\n", t.Cfg.ClientId, t.Cfg.MqttUrl, t.Cfg.Subscribe, t.Cfg.Qos, t.Cfg.Ping)
	// 开启链接
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	// 退出函数时断开链接
	defer func() {
		g.Log().Error(gctx.New(), t.Cfg.Name, "连接退出...")
		// 关闭连接
		c.Disconnect(250)
	}()
	// 订阅主题
	if t.Cfg.Subscribe != "false" {
		if token := c.Subscribe(t.Cfg.Subscribe, t.Cfg.Qos, nil); token.Wait() && token.Error() != nil {
			panic("订阅主题失败")
		}
	}
	// 写入客户端信息
	t.Client = &c
	// 写入全局
	MqttList.Set(t.Cfg.Name, t)
	// 维持运行
	select {}
}

func (t *Client) SendMsg(msg any, topic string, qos ...byte) error {
	ctx := gctx.New()
	// 如果没有初始化那么退出发送
	if !t.IsInit {
		logContext := "当前 MQTT 名称不存在 或 初始化失败"
		fmt.Println(logContext)
		g.Log().Warning(ctx, logContext)
		return nil
	}
	// 设置 qos
	var qosNumber = SendDefaultQos
	// 如果在配置文件中配置了 那么使用配置文件中的配置
	if len(qos) >= 1 {
		qosNumber = qos[0]
	}
	// 将传入消息解析为json
	json, err := gjson.EncodeString(msg)
	if err != nil {
		g.Log().Error(ctx, "mqtt 创建json出错", err.Error())
		return err
	}
	// 推送消息
	if token := (*t.Client).Publish(topic, qosNumber, false, json); token.Wait() && token.Error() != nil {
		err = token.Error()
		g.Log().Error(ctx, fmt.Sprintf("推送出现错误: %s", err))
		return err
	}
	return nil
}

func (t *Client) Json() *MqttResp {
	return CreateMqttResp(t)
}
