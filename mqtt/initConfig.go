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

type ClientOption struct {
	// MessageCallbackFunc 接收到消息时触发函数
	MessageCallbackFunc MessageHandler
	// OnConnectCallBackFunc 连接成功时触发函数
	OnConnectCallBackFunc mqtt.OnConnectHandler
}

type Config struct {
	Name      string `json:"name"`      // name 名称
	Debug     bool   `json:"debug"`     // debug 调试模式
	MqttUrl   string `json:"mqttUrl"`   // mqtt 链接地址
	ClientId  string `json:"clientId"`  // 客户端 id
	Subscribe string `json:"subscribe"` // 订阅地址
	Qos       byte   `json:"qos"`       // qos
	Username  string `json:"username"`  // 用户名
	Password  string `json:"password"`  // 密码
	Ping      int    `json:"ping"`      // ping 频率
}

// CreateClient 创建客户端
func CreateClient(optionHandler func(*ClientOption, *Config)) {
	cfg := initConfig()

	for _, config := range cfg {
		option := &ClientOption{}

		optionHandler(option, config)

		go (&Client{
			Cfg: config,
		}).SetMessageCallbackFunc(option.MessageCallbackFunc).SetOnConnectCallBackFunc(option.OnConnectCallBackFunc).Run()
	}

	select {}
}

// 初始化配置
// 需要初始化的参数: debug模式 链接地址 clientId 订阅 topic Qos
func initConfig() []*Config {
	ctx := gctx.New()

	mqttCfg, err := g.Cfg().Get(ctx, "mqtt")

	if err != nil {
		panic("MQTT 配置初始化失败")
	}

	mqttCfgData := mqttCfg.MapStrVar()

	var c []*Config

	for i, i2 := range mqttCfgData {
		v := i2.MapStrVar()

		var ping int

		if v["ping"].IsEmpty() {
			ping = 30
		} else {
			ping = v["ping"].Int()
		}

		c = append(c, &Config{
			Name:      i,
			Debug:     v["debug"].Bool(),
			MqttUrl:   v["url"].String(),
			ClientId:  v["clientId"].String(),
			Subscribe: v["subscribe"].String(),
			Qos:       byte(v["qos"].Int()),
			Username:  v["username"].String(),
			Password:  v["password"].String(),
			Ping:      ping,
		})
	}

	return c
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
