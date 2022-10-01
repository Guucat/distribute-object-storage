#!/bin/bash
#设置测试运行的环境变量并运行程序

export RABBITMQ_SERVER=amqp://test:test@localhost:5672
export ES_SERVER=localhost:9200

LISTEN_ADDRESS=10.29.1.1:12345 STORAGE_ROOT=/home/tsy/tmp/1 go run $1/data-server/dataServer.go &
LISTEN_ADDRESS=10.29.1.2:12345 STORAGE_ROOT=/home/tsy/tmp/2 go run $1/data-server/dataServer.go &
LISTEN_ADDRESS=10.29.1.3:12345 STORAGE_ROOT=/home/tsy/tmp/3 go run $1/data-server/dataServer.go &
LISTEN_ADDRESS=10.29.1.4:12345 STORAGE_ROOT=/home/tsy/tmp/4 go run $1/data-server/dataServer.go &
LISTEN_ADDRESS=10.29.1.5:12345 STORAGE_ROOT=/home/tsy/tmp/5 go run $1/data-server/dataServer.go &
LISTEN_ADDRESS=10.29.1.6:12345 STORAGE_ROOT=/home/tsy/tmp/6 go run $1/data-server/dataServer.go &

LISTEN_ADDRESS=10.29.2.1:12345 go run $1/api-server/apiServer.go &
LISTEN_ADDRESS=10.29.2.2:12345 go run $1/api-server/apiServer.go &
