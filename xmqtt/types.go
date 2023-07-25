package xmqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MessageHandlerData struct {
	// XMQTT 本扩展 MQTT 操作对象
	XMQTT *Client
	// OMQTT 原始 MQTT 操作对象
	OMQTT mqtt.Client
	// Message 接收到的消息操作对象
	Message mqtt.Message
}

// GetMessageId 获取 messageId
func (t *MessageHandlerData) GetMessageId() uint16 {
	return t.Message.MessageID()
}

// GetTopic 获取 topic
func (t *MessageHandlerData) GetTopic() string {
	return t.Message.Topic()
}

// GetMsg 获取接收到的消息内容
func (t *MessageHandlerData) GetMsg() string {
	return string(t.Message.Payload())
}

// GetEvent 获取事件
func (t *MessageHandlerData) GetEvent() (eventName string, data any, err error) {
	// 是否为系统客户端连接事件
	if data, ok, err := IsSystemConnectEvent(t.GetMsg(), t.GetTopic()); ok && err == nil {
		if data.Event == "connected" {
			return ConnectEvent, data, nil
		}
		if data.Event == "disconnected" {
			return DisconnectEvent, data, nil
		}
	}
	return NullEvent, nil, nil
}

type MessageHandler func(handlerData *MessageHandlerData)

type SystemConnectEvent struct {
	Event    string `json:"event"`
	ClientId string `json:"clientId"`
}

// Client 客户端结构对象
type Client struct {
	// Cfg 配置
	Cfg *Config
	// Client 客户端
	Client *mqtt.Client
	// MessageCallbackFunc 接收到消息时触发函数
	MessageCallbackFunc MessageHandler
	// OnConnectCallBackFunc 连接成功时触发函数
	OnConnectCallBackFunc mqtt.OnConnectHandler
}

// ClientCallBackOption 客户端回调设置
type ClientCallBackOption struct {
	// MessageCallbackFunc 接收到消息时触发函数
	MessageCallbackFunc MessageHandler
	// OnConnectCallBackFunc 连接成功时触发函数
	OnConnectCallBackFunc mqtt.OnConnectHandler
}

type Config struct {
	Name         string `json:"name"`         // name 名称
	Debug        bool   `json:"debug"`        // debug 调试模式
	MqttUrl      string `json:"mqttUrl"`      // mqtt 链接地址
	ClientId     string `json:"clientId"`     // 客户端 id
	Subscribe    string `json:"subscribe"`    // 订阅地址
	Qos          byte   `json:"qos"`          // qos
	Username     string `json:"username"`     // 用户名
	Password     string `json:"password"`     // 密码
	Ping         int    `json:"ping"`         // ping 频率
	CleanSession bool   `json:"cleanSession"` // cleanSession
}
