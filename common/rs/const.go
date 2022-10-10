// Package rs 封装了对 reed solo-man纠错码包的使用
package rs

const (
	DataShards    = 4
	ParityShards  = 2
	AllShards     = DataShards + ParityShards
	BlockPreShard = 8000
	BlockSize     = BlockPreShard * DataShards
)
