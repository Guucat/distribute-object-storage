package objects

import (
	"distribute-object-system/api-server/heartbeat"
	"distribute-object-system/api-server/locate"
	"distribute-object-system/common/es"
	"distribute-object-system/common/rs"
	"distribute-object-system/common/utils"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func post(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if locate.Exist(url.PathEscape(hash)) {
		e = es.AddVersion(name, hash, size)
		if e != nil {
			log.Println("fail to add version info on es:", e)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	}
	ds := heartbeat.ChooseRandomDataServers(rs.AllShards, nil)
	if len(ds) != rs.AllShards {
		log.Println("cannot find enough dataServer")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	stream, e := rs.NewResumablePutStream(ds, name, url.PathEscape(hash), size)
	if e != nil {
		log.Println("fail to get ResumablePutStream:", e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("location", "/temp/"+url.PathEscape(stream.ToToken()))
	w.WriteHeader(http.StatusCreated)
}
