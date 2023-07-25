package xmqtt

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/text/gstr"
)

// IsSystemConnectEvent 是否为监听系统连接事件
func IsSystemConnectEvent(message, topic string) (*SystemConnectEvent, bool, error) {
	disconnected := gstr.Count(topic, "/disconnected") > 0
	connected := gstr.Count(topic, "/connected") > 0

	var event string
	if disconnected {
		event = "disconnected"
	}
	if connected {
		event = "connected"
	}
	// 将接受到的事件数据解析
	data, err := gjson.DecodeToJson(message)
	if err != nil {
		return nil, false, err
	}
	clientId := data.Get("clientid").String()
	// 返回结果
	return &SystemConnectEvent{
		Event:    event,
		ClientId: clientId,
	}, gstr.Count(topic, "$SYS/brokers") > 0 && (disconnected || connected), nil
}
