package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// get all ip from given range
func mainFunction(start, end string, count int) ([]string, error) {
	allAddr := make([]string, count)
	splitStartAddr := strings.Split(start, ".")
	a, err := strconv.Atoi(splitStartAddr[0])
	if err != nil {
		return nil, err
	}
	b, err := strconv.Atoi(splitStartAddr[1])
	if err != nil {
		return nil, err
	}
	c, err := strconv.Atoi(splitStartAddr[2])
	if err != nil {
		return nil, err
	}
	d, err := strconv.Atoi(splitStartAddr[3])
	if err != nil {
		return nil, err
	}

	for i:=0; i<count; i++ {
		if d > 255 {
			d = 0
			c++
		}
		if c > 255 {
			c = 0
			b++
		}
		if b > 255 {
			b = 0
			a++
		}
		if a > 255 {
			return nil, errors.New("ip address is not exist")
		}

		temp := make([]string, 4)
		temp[0] = strconv.Itoa(a)
		temp[1] = strconv.Itoa(b)
		temp[2] = strconv.Itoa(c)
		temp[3] = strconv.Itoa(d)

		allAddr[i] = strings.Join(temp, ".")
		d++
	}
	count--
	fmt.Println("test: ")
	fmt.Println(allAddr)
	if allAddr[count] != end {
		return nil, errors.New("fail")
	}

	return allAddr, nil
}

func ipMaskToInt(netmask string) (int, error) {
	ipSplitArr := strings.Split(netmask, ".")
	if len(ipSplitArr) != 4 {
		return 0, fmt.Errorf("netmask:%v is not valid, pattern should like: 255.255.255.0", netmask)
	}
	ipv4MaskArr := make([]byte, 4)
	for i, value := range ipSplitArr {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("ipMaskToInt call strconv.Atoi error:[%v] string value is: [%s]", err, value)
		}
		if intValue > 255 {
			return 0, fmt.Errorf("netmask cannot greater than 255, current value is: [%s]", value)
		}
		ipv4MaskArr[i] = byte(intValue)
	}

	ones, _ := net.IPv4Mask(ipv4MaskArr[0], ipv4MaskArr[1], ipv4MaskArr[2], ipv4MaskArr[3]).Size()
	return ones, nil
}

func main() {
	start := "10.30.10.254"
	end := "10.30.11.2"
	count := 5
	mask := "255.255.223.0"
	n, err := ipMaskToInt(mask)
	if err != nil {
		fmt.Println("parse mask error: ", err)
	}
	fmt.Println(n)
	allip, err := mainFunction(start, end, count)
	if err != nil {
		fmt.Println("get allip faild: ", err)
	}
	fmt.Println(allip)

	test := make([]int, 4)
	test[0] = 0
	test[1] = 1
	test[2] = 2
	test[3] = 3
	for _, t := range test {
		fmt.Println(t)
	}
}
