package main

import (
	"distribute-object-system/standalone-version/objects"
	"log"
	"net/http"
	"os"
)

func init() {
	err1 := os.Setenv("STORAGE_ROOT", "/home/tsy/tem")
	if err1 != nil {
		//goland:noinspection ALL,Annotator
		panic("创建环境变量失败")
	}
	err2 := os.Setenv("LISTEN_ADDRESS", "127.0.0.1:8080")
	if err2 != nil {
		//goland:noinspection Annotator
		panic("创建环境变量失败")
	}
}

func main() {
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatalln(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
