#!/bin/bash
#这里可替换为你自己的执行程序，其他代码无需更改
CMS_NAME=drone-patrol-backend

# 创建日志目录
mkdir -p ./logs

# 检查可执行文件是否存在
if [ ! -f "./drone-patrol-backend" ]; then
    echo "Error: drone-patrol-backend executable not found!"
    echo "Please build the project first with: go build -o drone-patrol-backend"
    exit 1
fi

chmod 777 drone-patrol-backend
#使用说明，用来提示输入参数
usage() {
 echo "Usage: sh rundev.sh [start|stop|restart|status]"
 exit 1
}

#检查程序是否在运行
is_exist(){
 # shellcheck disable=SC2006
 # shellcheck disable=SC2009
 pid=`ps -ef|grep drone-patrol-backend|grep -v grep|awk '{print $2}' `
 #如果不存在返回1，存在返回0
 if [ -z "${pid}" ]; then
 return 1
 else
 return 0
 fi
}

#启动方法
start(){
 is_exist
 if [ $? -eq "0" ]; then
 echo "${CMS_NAME} is already running. pid=${pid} ."
 else

 nohup ./$CMS_NAME > ./logs/backend.log 2>&1 &
 echo "${CMS_NAME} start success"
 fi
}

#停止方法
stop(){
 is_exist
 if [ $? -eq "0" ]; then
 kill -9 $pid
 else
 echo "${CMS_NAME} is not running"
 fi
}

#输出运行状态
status(){
 is_exist
 if [ $? -eq "0" ]; then
 echo "${CMS_NAME} is running. Pid is ${pid}"
 else
 echo "${CMS_NAME} is NOT running."
 fi
}

#重启
restart(){
 stop
 start
}

#根据输入参数，选择执行对应方法，不输入则执行使用说明
case "$1" in
 "start")
 start
 ;;
 "stop")
 stop
 ;;
 "status")
 status
 ;;
 "restart")
 restart
 ;;
 *)
 usage
 ;;
esac