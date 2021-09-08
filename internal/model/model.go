package model

type Error struct {
	Error string `json:"error"`
}

type Service struct {
	PortNum int    `json:"port"`
	Name    string `json:"name"`
}

type Host struct {
	Ip        string    `json:"ip"`
	Os        string    `json:"os"`
	Timestamp int64     `json:"timestamp"`
	Ports     []Service `json:"service"`
}
