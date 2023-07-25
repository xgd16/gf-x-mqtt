# [GF](https://goframe.org/pages/viewpage.action?pageId=1114119)-X-MQTT

### 配置

```toml 
[mqtt]
    [mqtt.clientAdmin] # 此位置就是 clietName
        debug = false # 是否开启debug 
        url = "tcp://127.0.0.1:1884" # 连接目标
        clientId = "tuokeClient123" # 客户端id
        subscribe = "$SYS/brokers/emqx@172.17.0.2/clients/#" # 订阅频道 无需订阅 写 false
        qos = 0 # 协议质量 0 1 2
        username = "clientAdmin" # 用户名密码
        password = "clientAdmin321." # 密码
        cleanSession = false # 清空 session
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

### 快速接入示例

```go
package mqtt

import (
    "demo/src/service/mqtt/handler"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gctx"
    "github.com/xgd16/gf-x-mqtt/xmqtt"
)

// 注册 MQTT 处理
var register = map[string]func(data *xmqtt.EventHandlerData){
    xmqtt.ConnectEvent:    handler.Connect,   // 客户端连接事件
    xmqtt.DisconnectEvent: handler.Connect,   // 客户端断开连接事件
    xmqtt.NullEvent:       handler.NullEvent, // 没有事件时触发
}

func Service() {
    ctx := gctx.New()
    xmqtt.CreateClient(func(option *xmqtt.ClientCallBackOption, config *xmqtt.Config) {
        option.MessageCallbackFunc = func(data *xmqtt.MessageHandlerData) {
            // 获取 事件
            eventName, eventData, eventErr := data.GetEvent()
            if eventErr != nil {
                g.Log().Error(ctx, "MQTT 事件出错", eventErr)
                return
            }
            // 处理 事件
            register[eventName](&xmqtt.EventHandlerData{EventData: eventData, MsgHandlerData: data})
        }
    })
}

```

#### 处理事件

```go
package handler

import (
    "fmt"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/xgd16/gf-x-mqtt/xmqtt"
)

func NullEvent(data *xmqtt.EventHandlerData) {
    fmt.Println(data.MsgHandlerData.GetTopic(), data.MsgHandlerData.GetMsg())
    data.SendMsg(g.Map{
        "msg": data.GetJson().Get("msg").String() + "!!!!!!!!!!!!!!!",
    }, "a/1")
}
```

**PS: ``*xmqtt.EventHandlerData`` 中已实现 ``SendMsg``  操作 默认使用接收客户端用户进行发送操作 **

> *xmqtt.EventHandlerData 操作对象内的函数

SendMsg(msg any, topic string, qos ...byte) error 

**``*xmqtt.Client`` 中的 ``SendMsg`` 函数是此函数的原型 **

- ``GetJson`` 函数 获取订阅频道接收到的数据的json对象 **需要确保接收数据为 JSON**

GetJson() (json *gjson.Json)

### 细节操作

> 获取 ``MQTT`` 操作对象

```go
xmqtt.MqttList.Get("{配置里设置的MQTT名称}") // 获取到 *xmqtt.Client 操作对象
```

> *xmqtt.Client 操作对象内函数

- ``SendMsg`` 函数

SendMsg(msg any, topic string, qos ...byte) error

1. 参数 ``msg`` **要发送给客户端的数据输入任何类型会自动被转换成 json 数据发送给客户端**
2. 参数 ``topic`` **发送到那个订阅频段 例:** ``a/1``

3. qos **发送模式默认** ``0``



### 推荐服务端

> EMQX 免费好用的 MQTT 服务端

[EMQX: 大规模分布式 MQTT 消息服务器](https://www.emqx.io/zh)

> MQTTX 方便开发调试的客户端

[MQTTX：全功能 MQTT 客户端工具](https://mqttx.app/zh)

