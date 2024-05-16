package main

import (
	"encoding/json"
)

type MessageToServer struct {
	status bool
	id     string
	info   string
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
}

func (mtc MessageToClient) ToJson() []byte {
	js, err := json.Marshal(mtc)
	if err != nil {
		panic(err)
	}
	return js
}

type ServiceMessage struct {
	id   string
	name string
	ip   string
	port int
	args []string
}

type ClientMessage struct {
	name string
}
