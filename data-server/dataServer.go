package main

import (
	"distribute-object-system/data-server/heartbeat"
	"distribute-object-system/data-server/locate"
	"distribute-object-system/data-server/objects"
	"distribute-object-system/data-server/temp"
	"log"
	"net/http"
	"os"
)

func main() {
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
