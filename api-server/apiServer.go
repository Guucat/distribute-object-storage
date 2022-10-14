package main

import (
	"distribute-object-system/api-server/heartbeat"
	"distribute-object-system/api-server/locate"
	"distribute-object-system/api-server/objects"
	"distribute-object-system/api-server/temp"
	"distribute-object-system/api-server/versions"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	// 性能分析

	go heartbeat.ListenHeartBeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	// 断点上传系列
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatalln(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
