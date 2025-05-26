package xmqtt

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// CreateClient 创建客户端
func CreateClient(ctx context.Context, optionHandler func(*ClientCallBackOption, *Config)) {
	cfg := initConfig()
	// 根据配置循环创建监听的客户端
	for _, config := range cfg {
		option := new(ClientCallBackOption)
		// 调用外部进行处理的函数
		optionHandler(option, config)
		// 每个客户端在独立的协程中
		go (&Client{
			IsInit: true,
			Cfg:    config,
		}).SetMessageCallbackFunc(option.MessageCallbackFunc).SetOnConnectCallBackFunc(option.OnConnectCallBackFunc).Run(ctx)
	}
}

// 初始化配置
// 需要初始化的参数: debug模式 链接地址 clientId 订阅 topic Qos
func initConfig() []*Config {
	ctx := gctx.New()
	// 获取配置
	mqttCfg, err := g.Cfg().Get(ctx, "mqtt")
	if err != nil {
		panic("MQTT 配置初始化失败")
	}
	// 转换格式
	mqttCfgData := mqttCfg.MapStrVar()
	var c []*Config
	// 循环初始化数据
	for i, i2 := range mqttCfgData {
		v := i2.MapStrVar()
		// 设置默认值
		var ping int
		if v["ping"].IsEmpty() {
			ping = 30
		} else {
			ping = v["ping"].Int()
		}

		c = append(c, &Config{
			Name:         i,
			Debug:        v["debug"].Bool(),
			MqttUrl:      v["url"].String(),
			ClientId:     v["clientId"].String(),
			Subscribe:    v["subscribe"].String(),
			Qos:          byte(v["qos"].Int()),
			Username:     v["username"].String(),
			Password:     v["password"].String(),
			CleanSession: v["cleanSession"].Bool(),
			Ping:         ping,
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
