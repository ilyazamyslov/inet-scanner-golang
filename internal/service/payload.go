package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dean2021/go-nmap"
	"github.com/ilyazamyslov/inet-scanner-golang/internal/model"
)

func strIp2Int(ip string) (uint32, error) {
	sliceStrIP := strings.Split(ip, ".")
	sliceIntIP := make([]uint32, len(sliceStrIP))
	for i, s := range sliceStrIP {
		value, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		sliceIntIP[i] = uint32(value)
	}
	var intIP uint32 = (sliceIntIP[0] << 24) | (sliceIntIP[1] << 16) | (sliceIntIP[2] << 8) | sliceIntIP[3]
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

func scanHost(ip string) (model.Host, error) {
	var object model.Host
	n := nmap.New()
	args := []string{"-O"}
	n.SetArgs(args...)
	n.SetHosts(ip)
	object.Ip = ip
	err := n.Run()
	if err != nil {
		fmt.Println(ip, err)
		return object, err
	}
	result, err := n.Parse()
	if err != nil {
		fmt.Println(ip, err)
		return object, err
	}
	var (
		osName     string
		hostName   string
		osAccuracy = 0
	)
	for _, host := range result.Hosts {
		if host.Status.State == "up" {
			for _, osMatch := range host.Os.OsMatches {
				tempOsAccuracy, _ := strconv.Atoi(osMatch.Accuracy)
				if tempOsAccuracy >= osAccuracy {
					osName += osMatch.Name + " "
					osAccuracy = tempOsAccuracy
				}
			}
			object.Os = osName
			for _, hostMatch := range host.Hostnames {
				hostName += hostMatch.Name
			}
			object.Name = hostName
			for _, port := range host.Ports {
				object.Ports = append(object.Ports, model.Service{
					PortNum: port.PortId,
					Name:    port.Service.Name,
				})
			}
		}
	}
	object.Timestamp = time.Now().Unix()
	return object, nil
}

func listHosts(ip string, lenMask int) ([]string, error) {
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
	return listHost, nil
}

func scanNetwork(ip []string) ([]model.Host, error) {
	hosts := make([]model.Host, len(ip))
	sem := make(chan struct{}, 64)
	for i, host := range ip {
		sem <- struct{}{}
		go func(host string, i int) {
			hostObject, err := scanHost(host)
			if err != nil {
				return
			}
			hosts[i] = hostObject
			<-sem
		}(host, i)
	}
	return hosts, nil
}
