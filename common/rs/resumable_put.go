package rs

import (
	"distribute-object-system/common/objectstream"
	"distribute-object-system/common/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type resumableToken struct {
	Name    string
	Size    int64
	Hash    string
	Servers []string
	Uuids   []string
}

type ResumablePutStream struct {
	*PutStream
	*resumableToken
}

// NewResumablePutStream 初始化并返回断点上传流
func NewResumablePutStream(dataServers []string, name, hash string, size int64) (*ResumablePutStream, error) {
	putStream, e := NewPutStream(dataServers, hash, size)
	if e != nil {
		return nil, e
	}
	uuids := make([]string, AllShards)
	for i := range uuids {
		uuids[i] = putStream.writers[i].(*objectstream.TempPutStream).Uuid
	}
	token := &resumableToken{name, size, hash, dataServers, uuids}
	return &ResumablePutStream{putStream, token}, nil
}

func NewResumablePutStreamByToken(token string) (*ResumablePutStream, error) {
	b, e := base64.StdEncoding.DecodeString(token)
	if e != nil {
		return nil, e
	}

	var t resumableToken
	e = json.Unmarshal(b, &t)
	if e != nil {
		return nil, e
	}

	writers := make([]io.Writer, AllShards)
	for i := range writers {
		writers[i] = &objectstream.TempPutStream{Server: t.Servers[i], Uuid: t.Uuids[i]}
	}
	enc := NewEncoder(writers)
	return &ResumablePutStream{&PutStream{enc}, &t}, nil
}

func (s *ResumablePutStream) ToToken() string {
	b, _ := json.Marshal(s)
	return base64.StdEncoding.EncodeToString(b)
}

func (s *ResumablePutStream) CurrentSize() int64 {
	r, e := http.Head(fmt.Sprintf("http://%s/temp/%s", s.Servers[0], s.Uuids[0]))
	if e != nil {
		log.Println(e)
		return -1
	}
	if r.StatusCode != http.StatusOK {
		log.Println(r.StatusCode)
		return -1
	}
	size := utils.GetSizeFromHeader(r.Header) * DataShards
	// ?
	if size > s.Size {
		size = s.Size
	}
	return size
}
