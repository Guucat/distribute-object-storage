package rs

import (
	"distribute-object-system/common/objectstream"
	"fmt"
	"io"
)

type GetStream struct {
	*decoder
}

func NewGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*GetStream, error) {
	if len(locateInfo)+len(dataServers) != AllShards {
		return nil, fmt.Errorf("dataServers number mismatch")
	}

	readers := make([]io.Reader, AllShards)
	// 更新损坏的数据节点地址
	for i := 0; i < AllShards; i++ {
		server := locateInfo[i]
		if server == "" {
			locateInfo[i] = dataServers[0]
			dataServers = dataServers[1:]
			continue
		}
		//
		reader, e := objectstream.NewGetStream(server, fmt.Sprintf("%s.%d", hash, i))
		if e == nil {
			readers[i] = reader
		}
	}

	writers := make([]io.Writer, AllShards)
	// 设置单个分片数据容量 ?
	perShard := (size + DataShards - 1) / DataShards
	var e error
	// 获取写入流，以将缺失的数据分片重新写入dataServer节点
	for i := range readers {
		if readers[i] == nil {
			writers[i], e = objectstream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", hash, i), perShard)
			if e != nil {
				return nil, e
			}
		}
	}

	dec := NewDecoder(readers, writers, size)
	return &GetStream{dec}, nil
}

func (s *GetStream) Close() {
	for i := range s.writers {
		if s.writers[i] != nil {
			s.writers[i].(*objectstream.TempPutStream).Commit(true)
		}
	}
}

// Seek 从指定偏移offset处返回数据流
func (s *GetStream) Seek(offset int64, whence int) (int64, error) {
	if whence != io.SeekCurrent {
		panic("only support SeekCurrent")
	}
	if offset < 0 {
		panic("only support SeekCurrent")
	}
	for offset != 0 {
		length := int64(BlockSize)
		if offset < length {
			length = offset
		}
		buf := make([]byte, length)
		io.ReadFull(s, buf)
		offset -= length
	}
	return offset, nil
}
