package objects

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		put(w, r)
		return
	}
	if m == http.MethodGet {
		get(w, r)
		return
	}
	if m == http.MethodDelete {
		del(w, r)
		return
	}
	// 处理断点上传
	if m == http.MethodPost {
		post(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
