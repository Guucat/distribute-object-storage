package objects

import (
	"distribute-object-system/api-server/heartbeat"
	"distribute-object-system/api-server/locate"
	"distribute-object-system/common/es"
	"distribute-object-system/common/objectstream"
	"distribute-object-system/common/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing objects hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 调用dataServer临时存储服务
	size := utils.GetSizeFromHeader(r.Header)
	c, e := storeObject(r.Body, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(c)
		return
	}
	if c != http.StatusOK {
		w.WriteHeader(c)
		return
	}

	// 版本信息存入es
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	e = es.AddVersion(name, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	// 判断数据是否已经存储过
	if locate.Exist(url.PathEscape(hash)) {
		return http.StatusOK, nil
	}

	stream, e := putStream(url.PathEscape(hash), size)
	if e != nil {
		return http.StatusServiceUnavailable, e
	}

	//io.Copy(stream, r)
	//e = stream.Close()
	reader := io.TeeReader(r, stream)
	d := utils.CalculateHash(reader)
	if d != hash {
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("objects hash missmatch, calculate=%s, requested=%s", d, hash)
	}
	stream.Commit(true)
	return http.StatusOK, nil
}

func putStream(hash string, size int64) (*objectstream.TempPutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}
	return objectstream.NewTempPutStream(server, hash, size)
}
