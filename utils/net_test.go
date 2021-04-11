package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestGetLocalIPOnline(t *testing.T) {
	ip, err := GetLocalIPOnline()
	if err != nil {
		t.Error(err)
	}

	t.Log("GetLocalIPOnline: ", ip)
}

func TestGetLocalIPOffline(t *testing.T) {
	ip, err := GetLocalIPOffline()
	if err != nil {
		t.Error(err)
	}

	t.Log("GetLocalIPOffline: ", ip)
}

func TestIps(t *testing.T) {
	Ips()
	t.Log("ips: ")
}

func TestGetLocalIp4(t *testing.T) {
	fmt.Println(GetLocalIp4())
}

func TestGetLocalHostIP(t *testing.T) {
	ip, err := GetLocalHostIP()
	if err != nil {
		t.Error(err)
	}
	t.Log("before set env:", *ip)
}

func TestGetLocalHostIPBySetEnv(t *testing.T) {
	c := "192.168.8.9"
	os.Setenv(LocalHostEnv, c)

	ip, err := GetLocalHostIP()
	if err != nil {
		t.Error(err)
	}
	if *ip != c {
		t.Error("GetLocalHostIP expect localhost ip is", c, "here!")
	}

	t.Log("after set env:", *ip)
}
