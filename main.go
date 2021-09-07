package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/basho/riak-go-client"
	"github.com/dean2021/go-nmap"
)

type service struct {
	PortNum int
	Name    string
}

type host struct {
	Ip        string
	Os        string
	Timestamp int64
	Ports     []service
}

func scanByIp(ip string, wg *sync.WaitGroup) ([]byte, error) {
	defer wg.Done()
	var object host
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
				object.Ports = append(object.Ports, service{port.PortId, port.Service.Name})
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

const defaultIp = "192.168.1.69"

/*func insert() (err error) {

	return
}*/

/*
   Code samples from:
   http://docs.basho.com/riak/latest/dev/using/2i/
   make sure the 'indexes' bucket-type is created using the leveldb backend
*/

func main() {
	nodeOpts := &riak.NodeOptions{
		RemoteAddress: "127.0.0.1:8087",
	}

	var node *riak.Node
	var err error
	if node, err = riak.NewNode(nodeOpts); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	nodes := []*riak.Node{node}
	opts := &riak.ClusterOptions{
		Nodes: nodes,
	}

	cluster, err := riak.NewCluster(opts)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer func() {
		if err := cluster.Stop(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}()

	if err := cluster.Start(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	obj := &riak.Object{
		ContentType:     "text/plain",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Value:           []byte("this is a value in Riak"),
	}

	cmd, err := riak.NewStoreValueCommandBuilder().
		WithBucket("testBucketName").
		WithContent(obj).
		Build()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := cluster.Execute(cmd); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	svc := cmd.(*riak.StoreValueCommand)
	rsp := svc.Response
	fmt.Println(rsp.GeneratedKey)
}
