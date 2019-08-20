package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func RestrictIP(addr string, wetherIsSubnet bool) error{
	ip := strings.Split(addr, "/")[0]
	err0 := net.ParseIP(ip)
	if err0 == nil {
		return errors.New("ip error")
	}
	gatewayIP, broadcastIP, netCIDR, err := IPRange(addr)
	if err != nil {
		return err
	}
	if !wetherIsSubnet && ip == gatewayIP || ip == broadcastIP || ip == netCIDR {
		return errors.New("the ip has been exclude")
	}
	if wetherIsSubnet && ip != netCIDR {
		return errors.New("ip is not correct subnet")
	}
	return nil
}

//fill zero
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

//suffix of broadcastIP for IPv4
func completeDefIPv4(n int) (def string) {
	middle := make([]string, 0)
	for i:=0; i<n; i++ {
		middle = append(middle, "255")
	}
	temp := strings.Join(middle, ".")
	return temp
}

//calculate range of the IP net
//get gatewayIP and broadcastIP
func IPRange(addr string) (gatewayIP, broadcastIP string, netSeg string, err error){
	if !strings.Contains(addr, "/") {
		return "", "", "", errors.New("please input cidr")
	}
	var add string
	_, netCIDR, _ := net.ParseCIDR(addr)
	netNum, _ := strconv.Atoi(strings.Split(netCIDR.String(), "/")[1])
	if addr != "" && strings.Contains(addr, ".") {
		netAddr := strings.Split(netCIDR.String(), "/")[0]
		temp := make([]string, 0)
		GatewayIP := strings.Split(netAddr, ".")
		BroadcastIP := strings.Split(netAddr, ".")
		m := netNum / 8
		n := netNum % 8
		if n != 0 {
			add = calculateNetSeg(addr, n)
			//GatewayIP[m] = add
			if GatewayIP[m] != add {
				return "", "", "", errors.New("ip is not host ip! ")
			}
		}
		x, _ := strconv.Atoi(GatewayIP[3])
		x++
		GatewayIP[3] = strconv.Itoa(x)
		gatewayIP := strings.Join(GatewayIP, ".")
		for i:=0; i<m; i++ {
			temp = append(temp, BroadcastIP[i])
		}
		t := 4 - m
		middle := completeDefIPv4(t)
		temp = append(temp, middle)
		broadcastIP := strings.Join(temp, ".")
		return gatewayIP, broadcastIP, netCIDR.String(), nil
	}
	if addr != "" && strings.Contains(addr, ":") {
		ipv6, _ := transformIPv6(strings.Split(netCIDR.String(), "/")[0])
		temp := make([]string, 0)
		n := strings.Split(ipv6, ":")
		p := netNum / 16
		q := netNum % 16
		if q != 0 {
			add = calculateNetSeg(addr, q)
			//n[p] = add
			if n[p] != add {
				return "", "", "", errors.New("ip is not host ip! ")
			}
		}
		x, _ := strconv.ParseInt(n[7], 16, 32)
		x++
		h := fmt.Sprintf("%x", x)
		n[7] = h
		gatewayIP := strings.Join(n, ":")

		for i:=0; i<p; i++ {
			temp = append(temp, n[i])
		}
		t := 8 - p
		middle := completeDef(t, "f")
		temp = append(temp, middle)
		broadcastIP := strings.Join(temp, ":")
		return gatewayIP, broadcastIP, netCIDR.String(), nil
	}
	return "", "", "", errors.New("ip error")
}

//show integral IPv6
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

//if flag = 0, then return string to fill IPv6
//if flag = f, then return suffix of broadcastIP
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

//merge the prefix and suffix of IPv6
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

func calculateNetSeg(addr string, x int) (netseg string) {
	var t int
	temp := make([]string, 0)
	if strings.Contains(addr, ".") {
		t = 8
	}else if strings.Contains(addr, ":") {
		t = 16
	}
	for i:=0; i<t; i++ {
		temp = append(temp, "0")
	}
	for i:=0; i<x; i++ {
		temp[i] = "1"
	}
	seg := strings.Join(temp, "")
	segInt64, _ := strconv.ParseInt(seg, 2, 64)
	if t == 8 {
		segDeci := strconv.FormatInt(segInt64, 10)
		return segDeci
	}
	if t == 16 {
		segDeci := strconv.FormatInt(segInt64, 16)
		return segDeci
	}
	return ""
}

func main() {
	addr := "00ff::12/125"
	err := RestrictIP(addr, false)
	if err != nil {
		fmt.Println(err)
	}else {
		fmt.Println("you have a correct ip")
	}
}