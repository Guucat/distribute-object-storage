package temp

import (
	"distribute-object-system/data-server/locate"
	"log"
	"net/http"
	"os"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempInfo, e := readFromFile(uuid)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	dataFile := infoFile + ".dat"
	f, e := os.Open(dataFile)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	os.Remove(infoFile)
	if actual > tempInfo.Size {
		os.Remove(dataFile)
		log.Println("actual size mismatch, expect", tempInfo.Size, "actual", actual)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	commitTempObject(dataFile, tempInfo)
}

func commitTempObject(datFile string, temp *tempInfo) {
	os.Rename(datFile, os.Getenv("STORAGE_ROOT")+"/objects/"+temp.Name)
	locate.Add(temp.Name)
}
