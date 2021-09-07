package model

type Service struct {
	PortNum int
	Name    string
}

type Host struct {
	Ip        string
	Os        string
	Timestamp int64
	Ports     []Service
}
