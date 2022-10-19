package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

type MessageHandler func(*Client, mqtt.Client, mqtt.Message)

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

type Config struct {
	Debug     bool   `json:"debug"`     // debug 调试模式
	MqttUrl   string `json:"mqttUrl"`   // mqtt 链接地址
	ClientId  string `json:"clientId"`  // 客户端 id
	Subscribe string `json:"subscribe"` // 订阅地址
	Qos       byte   `json:"qos"`       // qos
}

// CreateClient 创建客户端
func CreateClient() *Client {
	cfg := initConfig()

	return &Client{
		Cfg: cfg,
	}
}

// 初始化配置
// 需要初始化的参数: debug模式 链接地址 clientId 订阅 topic Qos
func initConfig() *Config {
	ctx := gctx.New()

	debug, err := g.Cfg().Get(ctx, "mqtt.debug", false)
	mqttUrl, err := g.Cfg().Get(ctx, "mqtt.url")
	clientId, err := g.Cfg().Get(ctx, "mqtt.clientId")
	subscribe, err := g.Cfg().Get(ctx, "mqtt.subscribe")
	qos, err := g.Cfg().Get(ctx, "mqtt.qos", 0)

	if err != nil {
		panic("MQTT 配置初始化失败")
	}

	return &Config{
		Debug:     debug.Bool(),
		MqttUrl:   mqttUrl.String(),
		ClientId:  clientId.String(),
		Subscribe: subscribe.String(),
		Qos:       byte(qos.Int()),
	}
}

// SetMessageCallbackFunc 接收到消息时触发函数
func (t *Client) SetMessageCallbackFunc(fn MessageHandler) *Client {
	t.MessageCallbackFunc = fn

	return t
}

func (t *Client) SetOnConnectCallBackFunc(fn mqtt.OnConnectHandler) *Client {
	t.OnConnectCallBackFunc = fn

	return t
}
