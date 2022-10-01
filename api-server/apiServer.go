package main

import (
	"distribute-object-system/api-server/heartbeat"
	"distribute-object-system/api-server/locate"
	"distribute-object-system/api-server/objects"
	"distribute-object-system/api-server/versions"
	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.ListenHeartBeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatalln(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
