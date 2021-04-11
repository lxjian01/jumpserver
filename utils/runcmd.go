package utils

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func runInWindows(cmd string) ([]byte, error) {
	result, err := exec.Command("cmd", "/c", cmd).Output()
	if err != nil {
		return nil, err
	}
	return result, err
}

func RunCommand(cmd string) (string, error) {
	var out []byte
	var err error
	if runtime.GOOS == "windows" {
		out,err = runInWindows(cmd)
	} else {
		out,err = runInLinux(cmd)
	}
	return BytesToStr(out),err
}

func runInLinux(cmd string) ([]byte, error) {
	fmt.Println("Running Linux cmd:" + cmd)
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return nil, err
	}
	return result, err
}

//根据进程名判断进程是否运行
func CheckProRunning(serverName string) (bool, error) {
	a := `ps ux | awk '/` + serverName + `/ && !/awk/ {print $2}'`
	pid, err := RunCommand(a)
	if err != nil {
		return false, err
	}
	pid = strings.Replace(pid,"\n","",-1)
	pid = strings.Replace(pid," ","",-1)
	fmt.Println(pid)
	return pid != "", nil
}
//根据进程名称获取进程ID
func GetPid(serverName string) (string, error) {
	a := `ps ux | awk '/` + serverName + `/ && !/awk/ {print $2}'`
	pid, err := RunCommand(a)
	return pid , err
}