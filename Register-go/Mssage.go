package main

import (
	"encoding/json"
)

type MessageToServer struct {
	Status bool   `json:"status"`
	Id     string `json:"id"`
	Info   string `json:"info"`
}

func (mts MessageToServer) ToJson() []byte {
	js, err := json.Marshal(mts)
	if err != nil {
		panic(err)
	}
	return js
}

type MessageToClient struct {
	status     bool
	info       string
	serverList []string
	serverName string
	urlList    []string
	statusList []int
	args       []string
	ret        []string
}

func (mtc MessageToClient) ToJson() []byte {
	js, err := json.Marshal(mtc)
	if err != nil {
		panic(err)
	}
	return js
}

type ServiceMessage struct {
	Id   string   `json:"id"`
	Name string   `json:"name"`
	Ip   string   `json:"ip"`
	Port int      `json:"port"`
	Args []string `json:"args"`
	Ret  []string `json:"ret"`
}

type ClientMessage struct {
	name string
}
