// Package locate 用于向数据服务节点群发定位消息以确定存储的对象的位置信息，并通过临时队列接收反馈
package locate

import (
	"distribute-object-system/common/rabbitmq"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Handler 只接受查询对象位置信息的Get请求，并将结果写入响应
func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	info := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])
	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	b, _ := json.Marshal(info)
	w.Write(b)
}

// Locate 回对象位置信息, 通过dataServer Exchange向数据服务节点群发对象的名字, 并创建临时消息队列接受消息
// 函数会阻塞1s, 1s后未收到位置信息则返回空字符串""
func Locate(name string) string {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(1 * time.Second)
		q.Close()
	}()
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

// Exist 存在名称为s的对象则返回true
func Exist(name string) bool {
	return Locate(name) != ""
}
