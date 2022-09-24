package objects

import (
	"distribute-object-system/api-server/heartbeat"
	"distribute-object-system/common/objectstream"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	c, e := storeObject(r.Body, object)
	if e != nil {
		log.Println(e)
	}
	w.WriteHeader(c)
}

func storeObject(r io.Reader, object string) (int, error) {
	stream, e := putStream(object)
	if e != nil {
		return http.StatusServiceUnavailable, e
	}
	io.Copy(stream, r)
	e = stream.Close()
	if e != nil {
		return http.StatusInternalServerError, e
	}
	return http.StatusOK, nil
}

func putStream(object string) (*objectstream.PutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}
	return objectstream.NewPutStream(server, object), nil
}
