// Package es 封装了以http访问ES的各种API操作 version: 7.17.6
// 对象名 + 对象版本号 是对象元数据的唯一标识
// ES_SERVER: 环境变量，ES服务器地址
// metadata: 索引
// objects: 类型
package es

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Metadata struct {
	Name    string
	Version int
	Size    int64
	Hash    string
}

type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total int
		Hits  []hit
	}
}

// 根据id(对象名拼接对象版本号)获取对象元数据
func getMetadata(name string, versionId int) (meta Metadata, e error) {
	// GET <index>/_doc/<_id>	完整文档属性，包含文档源属性
	// GET <index>/_source/<_id>	只包含mapping属性
	url := fmt.Sprintf("http://%s/metadata/_source/%s_%d",
		os.Getenv("ES_SERVER"), name, versionId)
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to get %s_%d: %d", name, versionId, r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(result, &meta)
	return
}

// SearchLatestVersion 根据对象名搜索对象元数据，降序排列版本号并只返回第一条数据
// 若搜索为空，返回初始化结构体
func SearchLatestVersion(name string) (meta Metadata, e error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=name:%s&size=1&sort=version:desc",
		os.Getenv("ES_SERVER"), url.PathEscape(name))
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search metadata: %d", r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)

	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
	}
	return
}

// GetMetaData 根据对象名和对象版本号获取一条对象元数据，可能为空
func GetMetaData(name string, version int) (meta Metadata, e error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}

// PutMetadata 向ES服务上传一条新的文档，文档id由元数据的name和version拼接而成
func PutMetadata(name string, version int, size int64, hash string) error {
	doc := fmt.Sprintf(`{"name":"%s","version":%d,"size":%d,"hash":"%s"}`,
		name, version, size, hash)
	client := http.Client{}
	// PUT /<index>/_create/<_id>
	url := fmt.Sprintf("http://%s/metadata/_create/%s_%d",
		os.Getenv("ES_SERVER"), name, version)
	request, _ := http.NewRequest("PUT", url, strings.NewReader(doc))
	request.Header.Add("Content-Type", "application/json")
	r, e := client.Do(request)
	if e != nil {
		return e
	}
	// 由于使用了_created参数，如果同时有多个客户端上传同一个元数据，会发生冲突，且只有第一个文档被创建成功。
	// 之后的PUT请求， ES会返回409 Conflict， 此时函数让版本号+1后递归调用自身继续上传。
	if r.StatusCode == http.StatusConflict {
		return PutMetadata(name, version+1, size, hash)
	}
	if r.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("fail to put metadata: %d %s", r.StatusCode, result)
	}
	return nil
}

func AddVersion(name string, hash string, size int64) error {
	version, e := SearchLatestVersion(name)
	if e != nil {
		return e
	}
	return PutMetadata(name, version.Version+1, size, hash)
}

// SearchAllVersions 搜索指定对象name的所有版本的元数据，支持分页。
// 若name=""则搜索所有对象的所有版本元数据
func SearchAllVersions(name string, from, size int) ([]Metadata, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?sort=name,version&from=%d&size=%d",
		os.Getenv("ES_SERVER"), from, size)
	if name != "" {
		url += "&q=name:" + name
	}
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	metas := make([]Metadata, 0)
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	for _, j := range sr.Hits.Hits {
		metas = append(metas, j.Source)
	}
	return metas, nil
}
