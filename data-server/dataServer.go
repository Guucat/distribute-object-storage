package main

import (
	"distribute-object-system/data-server/heartbeat"
	"distribute-object-system/data-server/locate"
	"distribute-object-system/data-server/objects"
	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
