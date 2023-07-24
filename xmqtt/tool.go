package xmqtt

import "github.com/gogf/gf/v2/text/gstr"

// IsSystemConnectEvent 是否为监听系统连接事件
func IsSystemConnectEvent(topic string) bool {
	return gstr.Count(topic, "$SYS/brokers") > 0 &&
		(gstr.Count(topic, "disconnected") > 0 ||
			gstr.Count(topic, "connected") > 0)
}
