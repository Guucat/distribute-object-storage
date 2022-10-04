package locate

import (
	"distribute-object-system/common/rabbitmq"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var objects = make(map[string]int)
var mu sync.Mutex

func Locate(hash string) bool {
	mu.Lock()
	defer mu.Unlock()
	_, ok := objects[hash]
	return ok
}

func Add(hash string) {
	mu.Lock()
	defer mu.Unlock()
	objects[hash] = 1
}

func Del(hash string) {
	mu.Lock()
	mu.Unlock()
	delete(objects, hash)
}

func StartLocate() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		hash, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		exist := Locate(hash)
		if exist {
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}

// CollectObjects 缓存预热，将以存储的文件名加入缓存
func CollectObjects() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for _, file := range files {
		hash := filepath.Base(file)
		objects[hash] = 1
	}
}
