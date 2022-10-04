package objects

import (
	"distribute-object-system/common/utils"
	"distribute-object-system/data-server/locate"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
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

func getFile(hash string) string {
	file := os.Getenv("STORAGE_ROOT") + "/objects/" + hash
	f, _ := os.Open(file)
	d := url.PathEscape(utils.CalculateHash(f))
	f.Close()
	if d != hash {
		log.Println("objects hash mismatch, remove", file)
		locate.Del(hash)
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
