#!/bin/bash
SERVER="jumpserver"
WORKSPACE=$(cd `dirname $0`; pwd)
PID_FILE="${WORKSPACE}/${SERVER}.pid"
PID=0
INTERVAL=2

function check_pid() {
  echo $PID_FILE
  if [ ! -f "$PID_FILE" ]; then
    touch $PID_FILE
    chmod 777 $PID_FILE
    echo 0 > $PID_FILE
  fi
  PID=`cat $PID_FILE`
}

function start()
{
  check_pid
  echo $WORKSPACE
	if [[ $PID -gt 0 ]];then
	  echo "$SERVER is running......"
	  exit 0
  else
    nohup ${WORKSPACE}/${SERVER} > /dev/null 2>&1 &
    echo $!
    echo $! > ${PID_FILE}
  fi
	echo "sleeping..." &&  sleep $INTERVAL
	check_pid
  if [[ $PID -gt 0 ]];then
	  echo "Start $SERVER success."
	else
	  echo "Start $SERVER failed."
	fi
}

function restart()
{
  echo "restart"
}

function status()
{
  check_pid
	if [[ $PID -gt 1 ]];then
		echo "$SERVER is running by pid $PID"
		return 0
	else
		echo "$SERVER is not running"
		return 1
	fi
}

function stop()
{
  check_pid
  if [[ $PID -ne 0 ]]; then
    echo "$SERVER will be stopped"
    kill $PID
    echo 0 > ${PID_FILE}
  else
    echo "$SERVER already stopped!!"
  fi
  echo "sleeping..." &&  sleep $INTERVAL
  check_pid
  echo $PID
  if [[ $PID -ne 0 ]]; then
    echo "$SERVER stop failed!"
    echo "$SERVER kill failed,please use  kill -9 $PID"
  else
    echo "$SERVER stop success!"
  fi
}

function version()
{
  ${WORKSPACE}/jumpserver version
}

case "$1" in
	'start')
	start
	;;
	'stop')
	stop
	;;
	'status')
	status
	;;
	'restart')
	stop && start
	;;
  'version')
  version
  ;;
	*)
	echo "usage: $0 {start|stop|restart|status|version}"
	exit 1
	;;
esac