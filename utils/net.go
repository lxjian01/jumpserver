package utils

import (
	"errors"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	l         sync.Mutex
	localhost *string
)

func GetLocalHostIP() (*string, error) {
	if localhost == nil {
		l.Lock()
		defer l.Unlock()

		if localhost == nil {
			alh := os.Getenv(LocalHostEnv)
			if alh != "" {
				localhost = &alh
			} else {
				//lh, err := GetLocalIPOffline()
				//if err != nil {
				//	return nil, err
				//}
				lh := GetLocalIp4()
				localhost = &lh
			}
		}
	}

	return localhost, nil
}

func GetLocalIPOnline() (string, error) {
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	return strings.Split(conn.LocalAddr().String(), ":")[0], nil

}

func GetLocalIPOffline() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("not found")
}

func Ips() (map[string]string, error) {
	ips := make(map[string]string)
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range interfaces {
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			return nil, err
		}
		addresses, err := byName.Addrs()
		for _, v := range addresses {
			if ipnet, ok := v.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ips[byName.Name] = ipnet.IP.String()
				}
			}
		}
	}
	return ips, nil
}

func GetLocalIp4() (ip string) {
	mapIf4Ip, err := Ips()
	if err != nil || len(mapIf4Ip) == 0 {
		return ""
	}
	if len(mapIf4Ip) == 1 {
		for _, ip := range mapIf4Ip {
			return ip
		}
	}

	if ip, ok := mapIf4Ip["cni0"]; ok && !strings.HasPrefix(ip, "172.") {
		return ip
	}

	if ip, ok := mapIf4Ip["bond0"]; ok && !strings.HasPrefix(ip, "172.") {
		return ip
	}

	if ip, ok := mapIf4Ip["bond4"]; ok && !strings.HasPrefix(ip, "172.") {
		return ip
	}

	if ip, ok := mapIf4Ip["eth0"]; ok && !strings.HasPrefix(ip, "172.") {
		return ip
	}

	if ip, ok := mapIf4Ip["eth1"]; ok && !strings.HasPrefix(ip, "172.") {
		return ip
	}
	for _, ip := range mapIf4Ip {
		return ip
	}
	return ""
}

//func GetLocalIp4() (ip string) {
//	interfaces, err := net.Interfaces()
//	if err != nil {
//		return
//	}
//
//	for _, face := range interfaces {
//		if strings.Contains(face.Name, "lo") {
//			continue
//		}
//		addrs, err := face.Addrs()
//		if err != nil {
//			return
//		}
//
//		for _, addr := range addrs {
//			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
//				if ipnet.IP.To4() != nil {
//					currIp := ipnet.IP.String()
//					if !strings.Contains(currIp, ":") && currIp != "127.0.0.1" && !isIntranetIpv4(currIp) {
//						ip = currIp
//					}
//				}
//			}
//		}
//	}
//
//	return
//}

func isIntranetIpv4(ip string) bool {
	//if strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "169.254.") {
	if strings.HasPrefix(ip, "172.") || strings.HasPrefix(ip, "169.254.") {
		return true
	}
	return false
}
