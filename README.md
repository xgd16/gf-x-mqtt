# [GF](https://goframe.org/pages/viewpage.action?pageId=1114119)-X-MQTT
### 配置
    [mqtt]
    debug = false   // 调试模式
    url = "tcp://192.168.0.220:1883"   // 连接地址
    clientId = "sdt_service_client"   // 客户端ID
    subscribe = "sdt/#"   // 订阅
    qos = 2   // qos