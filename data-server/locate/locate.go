package locate

import (
	"distribute-object-system/common/rabbitmq"
	"distribute-object-system/common/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var objects = make(map[string]int)
var mu sync.Mutex

func Locate(hash string) int {
	mu.Lock()
	defer mu.Unlock()
	id, ok := objects[hash]
	if !ok {
		id = -1
	}
	return id
}

func Add(hash string, id int) {
	mu.Lock()
	defer mu.Unlock()
	objects[hash] = id
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
		id := Locate(hash)
		if id != -1 {
			q.Send(msg.ReplyTo, types.LocateInfo{Addr: os.Getenv("LISTEN_ADDRESS"), Id: id})
		}
	}
}

// CollectObjects 缓存预热，将以存储的文件名加入缓存
func CollectObjects() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for _, path := range files {
		file := strings.Split(filepath.Base(path), ".")
		if len(file) != 3 {
			panic(path)
		}
		hash := file[0]
		id, e := strconv.Atoi(file[1])
		if e != nil {
			panic(e)
		}
		objects[hash] = id
		//hash := filepath.Base(file)
		//objects[hash] = 1
	}
}
