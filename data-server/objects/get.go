package objects

import (
	"distribute-object-system/common/utils"
	"distribute-object-system/data-server/locate"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	file := getFile(strings.Split(r.URL.EscapedPath(), "/")[2])
	if file == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	sendFile(w, file)
}

func getFile(fileHash string) string {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + fileHash + ".*")
	// 正常情况下应该只存在一个对象的分片数据
	if len(files) != 1 {
		return ""
	}

	file := files[0]
	f, _ := os.Open(file)
	defer f.Close()
	curShardHash := url.PathEscape(utils.CalculateHash(f))
	shardHash := strings.Split(file, ".")[2]

	if curShardHash != shardHash {
		log.Println("objects hash mismatch, remove", file)
		locate.Del(shardHash)
		os.Remove(file)
		return ""
	}
	return file
}

func sendFile(w io.Writer, file string) {
	f, _ := os.Open(file)
	defer f.Close()
	io.Copy(w, f)
}
