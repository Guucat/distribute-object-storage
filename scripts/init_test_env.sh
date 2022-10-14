#!/bin/bash
# 创建测试的工作目录和网络环境

for i in `seq 1 6`
do
    mkdir -p /home/tsy/tmp/$i/objects
    mkdir -p /home/tsy/tmp/$i/temp
done


# eth0 是此linux系统的网络接口名，使用ifconfig查看
# Ubuntu20.04 LTS 内核支持接口别名 使用ifconfig命令设置别名
sudo ifconfig eth0:1 10.29.1.1/16
sudo ifconfig eth0:2 10.29.1.2/16
sudo ifconfig eth0:3 10.29.1.3/16
sudo ifconfig eth0:4 10.29.1.4/16
sudo ifconfig eth0:5 10.29.1.5/16
sudo ifconfig eth0:6 10.29.1.6/16
sudo ifconfig eth0:7 10.29.2.1/16
sudo ifconfig eth0:8 10.29.2.2/16

#rabbitMQ
#python3 rabbitmqadmin declare exchange name=apiServers type=fanout
#python3 rabbitmqadmin declare exchange name=dataServers type=fanout