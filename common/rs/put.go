package rs

import (
	"distribute-object-system/common/objectstream"
	"fmt"
	"io"
)

type PutStream struct {
	*encoder
}

func NewPutStream(dataServers []string, hash string, size int64) (*PutStream, error) {
	if len(dataServers) != AllShards {
		return nil, fmt.Errorf("dataServers number mismatch")
	}
	// 计算单个分片的数据量
	perShard := (size + DataShards - 1) / DataShards
	writers := make([]io.Writer, AllShards)
	var e error
	for i := range writers {
		writers[i], e = objectstream.NewTempPutStream(dataServers[i], fmt.Sprintf("%s.%d", hash, i), perShard)
		if e != nil {
			return nil, e
		}
	}
	enc := NewEncoder(writers)
	return &PutStream{enc}, nil
}

func (s *PutStream) Commit(success bool) {
	s.Flush()
	for i := range s.writers {
		s.writers[i].(*objectstream.TempPutStream).Commit(success)
	}
}
