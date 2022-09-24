package heartbeat

import (
	"distribute-object-system/common/rabbitmq"
	"os"
	"time"
)

// StartHeartbeat 通过消息队列RabbitMQ
// 每5s向接口服务apiServer发送本服务节点的监听地址(ip + port)
func StartHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	for {
		q.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}
