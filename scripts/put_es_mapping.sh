#!/bin/bash
# ES6.0之后，ES对content-type的检测更为严格
# 在请求时需要增加content-type类型
# 当前版本为 es 7.17:
# mapping 不再支持类型, string数据类型 分为 keyword 和 text 两种
# keyword: not_analyzed 用于关键词搜索
# text: analyzed        用于全文检索
# es2 迁移 es7:
#      "type": "string",     "type": "text", /  "type": "string",        "type": "keyword",
 #     "index": "analyzed"    "index": true  /  index": "not_analyzed"    "index": true

#curl localhost:9200/metadata -XPUT -d'{"mappings":{"objects":{"properties":{"name":{"type":"string","index":"not_analyzed"},"version":{"type":"integer"},"size":{"type":"integer"},"hash":{"type":"string"}}}}}'

curl localhost:9200/metadata -H 'Content-Type: application/json' -XPUT -d'{"mappings":{"properties":{"name":{"type":"keyword","index":true},"version":{"type":"integer"},"size":{"type":"integer"},"hash":{"type":"text","index":true}}}}'
