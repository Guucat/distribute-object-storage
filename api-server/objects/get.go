package objects

import (
	"distribute-object-system/api-server/heartbeat"
	"distribute-object-system/api-server/locate"
	"distribute-object-system/common/es"
	"distribute-object-system/common/rs"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	versionId := r.URL.Query()["version"]
	version := 0
	var e error
	if len(versionId) != 0 {
		version, e = strconv.Atoi(versionId[0])
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	meta, e := es.GetMetaData(name, version)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if meta.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	hash := url.PathEscape(meta.Hash)
	stream, e := getStream(hash, meta.Size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err := io.Copy(w, stream)
	if err != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	stream.Close()
}

func getStream(hash string, size int64) (*rs.GetStream, error) {
	locateInfo := locate.Locate(hash)
	// 数据片缺失超过阈值，损坏数据不可修复
	if len(locateInfo) < rs.DataShards {
		return nil, fmt.Errorf("object %s locate fail, result %v", hash, locateInfo)
	}
	dataServers := make([]string, 0)
	// 数据片缺失数在可修复范围内，使用Reed Solomon纠错码修复数据
	if len(locateInfo) != rs.AllShards {
		dataServers = heartbeat.ChooseRandomDataServers(rs.AllShards-len(locateInfo), locateInfo)
	}
	return rs.NewGetStream(locateInfo, dataServers, hash, size)
	//return objectstream.NewGetStream(server, object)
}
