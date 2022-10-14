package rs

import (
	"distribute-object-system/common/objectstream"
	"io"
)

// ResumableGetStream 封装断点下载操作
type ResumableGetStream struct {
	*decoder
}

func NewResumableGetStream(dataServers []string, uuids []string, size int64) (*ResumableGetStream, error) {
	readers := make([]io.Reader, AllShards)
	var e error
	for i := 0; i < AllShards; i++ {
		readers[i], e = objectstream.NewTempGetStream(dataServers[i], uuids[i])
		if e != nil {
			return nil, e
		}
	}
	writers := make([]io.Writer, AllShards)
	dec := NewDecoder(readers, writers, size)
	return &ResumableGetStream{dec}, nil
}
