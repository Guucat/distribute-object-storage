// Package heartbeat 用于接受和处理来自数据服务节点dataServer的心跳消息
package heartbeat

import (
	"distribute-object-system/common/rabbitmq"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var dataServers = make(map[string]time.Time) // 保存数据服务节点的最新的心跳时间
var mutex sync.Mutex

// ListenHeartBeat 监听每一个来自数据服务节点的心跳信息，刷新过期时间
func ListenHeartBeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("apiServers")
	c := q.Consume()
	go removeExpiredDataServer()
	for msg := range c {
		dataServer, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}

// 删除心跳时间超过10s未发送心跳消息的过期数据服务节点
func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}

// GetDataServer 获取所有未过期的数据服务节点的监听地址
func GetDataServer() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds
}

// ChooseRandomDataServer 随机返回一个未过期的数据服务节点的监听地址(ip + port)
func ChooseRandomDataServer() string {
	ds := GetDataServer()
	n := len(ds)
	if n == 0 {
		return ""
	}
	return ds[rand.Intn(n)]
}
