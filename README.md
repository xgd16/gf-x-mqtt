# [GF](https://goframe.org/pages/viewpage.action?pageId=1114119)-X-MQTT
### 配置
```toml
[mqtt]
    debug = false # 调试模式
    url = "tcp://192.168.0.220:1883" # 连接地址
    clientId = "sdt_service_client" # 客户端ID
    subscribe = "sdt/#" # 订阅
    qos = 2 # qos
    username = "service" # 用户名
    password = "486213" # 密码
```
### 代码演示
```go
package main

import (
    "fmt"
    "github.com/xgd16/gf-x-mqtt/xmqtt"
)

func main() {
    xmqtt.CreateClient(func(option *xmqtt.ClientCallBackOption, config *xmqtt.Config) {
        option.MessageCallbackFunc = func(data *xmqtt.MessageHandlerData) {
            fmt.Println(data.GetMessageId(), data.GetTopic(), data.GetMsg())
            //client.SendMsg("收到", "sdt/c/1")
        }
    })
    select {}
}
```

# 建议
    服务端订阅和客户端订阅分开防止串扰
### 示例
    将服务端订阅到 sdt/s/# 监听所有发送给服务端的消息
    用户端发送消息时 topic 发布到 sdt/s/1 上
    客户端监听到 sdt/c/1
    服务端向用户推送消息时 发布到 sdt/c/1 上