#!/bin/bash
#前端服务启动脚本
FRONTEND_NAME="无人机监控前端服务"

#使用说明，用来提示输入参数
usage() {
 echo "Usage: sh run-frontend.sh [start|stop|restart|status]"
 exit 1
}

#检查程序是否在运行
is_exist(){
 # shellcheck disable=SC2006
 # shellcheck disable=SC2009
 pid=`ps -ef|grep "http-server.*dist"|grep -v grep|awk '{print $2}' `
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
 echo "${FRONTEND_NAME} is already running. pid=${pid} ."
 else
 # 确保logs目录存在
 mkdir -p ./logs
 # 启动http-server服务静态文件
 nohup npx http-server dist -p 3000 -a 0.0.0.0 > ./logs/frontend.log 2>&1 &
 echo "${FRONTEND_NAME} start success"
 echo "前端服务已启动，访问地址: http://0.0.0.0:3000"
 fi
}

#停止方法
stop(){
 is_exist
 if [ $? -eq "0" ]; then
 kill -9 $pid
 echo "${FRONTEND_NAME} stopped"
 else
 echo "${FRONTEND_NAME} is not running"
 fi
}

#输出运行状态
status(){
 is_exist
 if [ $? -eq "0" ]; then
 echo "${FRONTEND_NAME} is running. Pid is ${pid}"
 else
 echo "${FRONTEND_NAME} is NOT running."
 fi
}

#重启
restart(){
 stop
 sleep 2
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
