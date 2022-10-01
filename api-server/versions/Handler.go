package versions

import (
	"distribute-object-system/common/es"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	from := 0
	size := 1000
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	for {
		metas, e := es.SearchAllVersions(name, from, size)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, m := range metas {
			b, _ := json.Marshal(m)
			w.Write(b)
			w.Write([]byte("\n"))
		}
		// 服务器没有多余的数据，则结束循环
		if len(metas) != size {
			return
		}
		from += size
	}
}
