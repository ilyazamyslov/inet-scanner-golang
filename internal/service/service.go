package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dean2021/go-nmap"
)

func scanByIp(ip string, wg *sync.WaitGroup) ([]byte, error) {
	defer wg.Done()
	var object Host
	n := nmap.New()

	args := []string{"-O"}
	n.SetArgs(args...)
	n.SetHosts(ip)
	object.Ip = ip

	err := n.Run()
	if err != nil {
		return nil, err
	}
	result, err := n.Parse()
	if err != nil {
		return nil, err
	}

	var (
		osName     string
		osAccuracy = 0
	)
	for _, host := range result.Hosts {
		if host.Status.State == "up" {
			for _, osMatch := range host.Os.OsMatches {
				tempOsAccuracy, _ := strconv.Atoi(osMatch.Accuracy)
				if tempOsAccuracy >= osAccuracy {
					osName = osMatch.Name
					osAccuracy = tempOsAccuracy
				}
			}
			object.Os = osName
			for _, port := range host.Ports {
				object.Ports = append(object.Ports, Service{port.PortId, port.Service.Name})
			}
		}
	}
	object.Timestamp = time.Now().Unix()
	jsonObject, err := json.Marshal(object)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(object)
	return jsonObject, nil
}

func strIp2Int(ip string) (uint32, error) {
	sliceStrIP := strings.Split(ip, ".")
	sliceIntIP := make([]uint8, len(sliceStrIP))
	for i, s := range sliceStrIP {
		value, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		sliceIntIP[i] = uint8(value)
	}
	var intIP uint32 = (uint32(sliceIntIP[0]) << 24) | (uint32(sliceIntIP[1]) << 16) | (uint32(sliceIntIP[2]) << 8) | uint32(sliceIntIP[3])
	return intIP, nil
}

func intIp2Str(ip uint32) string {
	sliceStrIP := make([]string, 4)
	sliceStrIP[0] = strconv.Itoa(int(ip >> 24))
	sliceStrIP[1] = strconv.Itoa(int(ip << 8 >> 24))
	sliceStrIP[2] = strconv.Itoa(int(ip << 16 >> 24))
	sliceStrIP[3] = strconv.Itoa(int(ip << 24 >> 24))
	return sliceStrIP[0] + "." + sliceStrIP[1] + "." + sliceStrIP[2] + "." + sliceStrIP[3]
}

func mask(lenMask int) (mask uint32, err error) {
	if lenMask > 32 {
		return 0, fmt.Errorf("invalid lenMask")
	}
	for i := 0; i < (32 - lenMask); i++ {
		mask |= 1 << i
	}
	mask = mask ^ uint32(0xffffffff)
	return mask, nil
}

func scanNetwork(ip string, lenMask int) ([]byte, error) {
	intIp, err := strIp2Int(ip)
	if err != nil {
		return nil, err
	}
	mask, err := mask(lenMask)
	if err != nil {
		return nil, err
	}
	network := intIp & mask
	currentHost := network + 1
	listHost := []string{}
	for (currentHost & mask) == network {
		listHost = append(listHost, intIp2Str(currentHost))
		currentHost += 1
	}
	var wg sync.WaitGroup
	wg.Add(len(listHost))
	for _, host := range listHost {
		go scanByIp(host, &wg)
	}
	wg.Wait()
	return nil, nil
}
