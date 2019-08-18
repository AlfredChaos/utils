package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func RestrictIP(addr string, subnet bool) error{
	gatewayIP, broadcastIP, _, err := IPRange(addr)
	if err != nil {
		return err
	}
	_, netCIDR, _ := net.ParseCIDR(addr)


	return nil
}

func completeIPv6(addr string) (ip string) {
	ipv6 := strings.Split(addr, ":")
	aloneBit := make([]string, 2)
	for i:=0; i<len(ipv6); i++ {
		if len(ipv6[i]) == 1 && ipv6[i] != "" {
			aloneBit[0] = "000"
			aloneBit[1] = ipv6[i]
			ipv6[i] = strings.Join(aloneBit, "")
		} else if len(ipv6[i]) == 2 && !strings.Contains(ipv6[i], "/") {
			aloneBit[0] = "00"
			aloneBit[1] = ipv6[i]
			ipv6[i] = strings.Join(aloneBit, "")
		} else if len(ipv6[i]) == 3 && !strings.Contains(ipv6[i], "/") {
			aloneBit[0] = "0"
			aloneBit[1] = ipv6[i]
			ipv6[i] = strings.Join(aloneBit, "")
		}
	}
	address := strings.Join(ipv6, ":")
	return address
}

func IPRange(addr string) (gatewayIP, broadcastIP string, err error){
	if addr != "" && strings.Contains(addr, ".") {
		var i int
		_, netCIDR, _ := net.ParseCIDR(addr)
		netAddr := strings.Split(netCIDR.String(), "/")[0]
		netNum, _ := strconv.Atoi(strings.Split(netCIDR.String(), "/")[1])
		GatewayIP := strings.Split(netAddr, ".")
		BroadcastIP := strings.Split(netAddr, ".")
		temp, _ := strconv.Atoi(GatewayIP[3])
		GatewayIP[3] = strconv.Itoa(temp + 1)
		if netNum >= 24 && netNum < 32 {
			i = 3
			BroadcastIP[i] = "255"
			return strings.Join(GatewayIP, "."), strings.Join(BroadcastIP, "."), nil
		}else if netNum >= 16 && netNum < 24 {
			i = 2
			BroadcastIP[2] = "255"
			BroadcastIP[3] = "255"
			return strings.Join(GatewayIP, "."), strings.Join(BroadcastIP, "."), nil
		}else if netNum >= 8 && netNum <16 {
			i = 1
			BroadcastIP[1] = "255"
			BroadcastIP[2] = "255"
			BroadcastIP[3] = "255"
			return 	strings.Join(GatewayIP, "."), strings.Join(BroadcastIP, "."), nil
		}else if netNum > 0 && netNum < 8 {
			i = 0
			BroadcastIP[0] = "255"
			BroadcastIP[1] = "255"
			BroadcastIP[2] = "255"
			BroadcastIP[3] = "255"
			return 	strings.Join(GatewayIP, "."), strings.Join(BroadcastIP, "."), nil
		}else {
			return "", "", errors.New("ip error")
		}

	}
	if addr != "" && strings.Contains(addr, ":") {
		_, netCIDR, _ := net.ParseCIDR(addr)
		ipv6, _ := transformIPv6(strings.Split(netCIDR.String(), "/")[0])
		temp := make([]string, 0)
		var broadcastIP string
		n := strings.Split(ipv6, ":")
		x, _ := strconv.ParseInt(n[7], 16, 32)
		x++
		h := fmt.Sprintf("%x", x)
		n[7] = h
		gatewayIP := strings.Join(n, ":")
		netSeg, _ := strconv.Atoi(strings.Split(addr, "/")[1])
		p := netSeg / 16
		//q := netSeg % 16
		for i:=0; i<p; i++ {
			temp = append(temp, n[i])
		}
		t := 8 - p
		middle := completeDef(t, "f")
		temp = append(temp, middle)
		broadcastIP = strings.Join(temp, ":")
		return gatewayIP, broadcastIP, nil
	}
	return "", "", errors.New("ip error")
}

func transformIPv6(addr string) (ip string, err error) {
	var address string
	if strings.Contains(addr, "::") && addr != "" {
		ipv6 := strings.Split(addr, "::")
		if len(ipv6) > 2 {
			return "", errors.New("ipv6 error")
		}else if len(ipv6) == 2 {
			m := len(strings.Split(ipv6[0], ":"))
			n := len(strings.Split(ipv6[1], ":"))
			if ipv6[0] == "" && ipv6[1] != "" {
				m = 0
			}else if ipv6[0] != "" && ipv6[1] == "" {
				n = 0
			}
			x := 8 - m - n
			temp := completeDef(x, "0")
			address = merge(addr, temp)
		}
	}
	return address, nil
}

func completeDef(x int, flag string) (def string) {
	middle := make([]string, 0)
	for i:=0; i<x; i++ {
		if flag == "0" {
			middle = append(middle, "0000")
		}else if flag == "f" {
			middle = append(middle, "ffff")
		}
	}
	temp := strings.Join(middle, ":")
	return temp
}

func merge(addr, add string) (ip string){
	var address string
	temp := make([]string, 0)
	ipv6 := strings.Split(addr, "::")
	if ipv6[0] == "" && ipv6[1] != ""{
		temp = append(temp, add)
		temp = append(temp, ipv6[1])
		address = strings.Join(temp, ":")
	}else if ipv6[1] == "" && ipv6[0] != ""{
		temp = append(temp, ipv6[0])
		temp = append(temp, add)
		address = strings.Join(temp, ":")
	}else if ipv6[1] != "" && ipv6[0] != "" {
		temp = append(temp, ipv6[0])
		temp = append(temp, add)
		temp = append(temp, ipv6[1])
		address = strings.Join(temp, ":")
	}
	return address
}

func main() {
	addr := "ff:67::0020/66"
}