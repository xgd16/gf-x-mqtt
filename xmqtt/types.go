package xmqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gogf/gf/v2/encoding/gjson"
	"sync"
)

// SafeMQTTList 安全 MQTT 列表
type SafeMQTTList struct {
	sync.RWMutex
	M map[string]*Client
}

// CreateSafeMQTTList 创建 安全 MQTT 列表
func CreateSafeMQTTList() *SafeMQTTList {
	return &SafeMQTTList{
		M: make(map[string]*Client),
	}
}

// Get 获取 MQTT 客户端对象
func (t *SafeMQTTList) Get(mqttName string) *Client {
	t.RLock()
	defer t.RUnlock()
	if data, ok := t.M[mqttName]; ok {
		return data
	}
	return &Client{IsInit: false}
}

// Set 设置 MQTT 客户端对象
func (t *SafeMQTTList) Set(mqttName string, client *Client) {
	t.Lock()
	defer t.Unlock()
	t.M[mqttName] = client
}

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
	} else {
		if err != nil {
			return "", nil, err
		}
	}
	return NullEvent, nil, nil
}

// EventHandlerData 事件处理数据
type EventHandlerData struct {
	MsgHandlerData *MessageHandlerData
	EventData      any
}

// SendMsg 发送消息
func (t *EventHandlerData) SendMsg(msg any, topic string, qos ...byte) error {
	return t.MsgHandlerData.XMQTT.SendMsg(msg, topic, qos...)
}

// GetJson 获取 JSON 对象
func (t *EventHandlerData) GetJson() (json *gjson.Json) {
	json, _ = gjson.DecodeToJson(t.MsgHandlerData.GetMsg())
	return
}

// MessageHandler 消息处理函数
type MessageHandler func(handlerData *MessageHandlerData)

// SystemConnectEvent 系统客户端连接事件
type SystemConnectEvent struct {
	Event    string `json:"event"`
	ClientId string `json:"clientId"`
}

// IsCurrentService 是否为当前服务
func (t *SystemConnectEvent) IsCurrentService(clientId string) bool {
	return t.ClientId == clientId
}

// Client 客户端结构对象
type Client struct {
	IsInit bool
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
