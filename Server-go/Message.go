package main

import "encoding/json"

type MessageToServer struct {
	Status bool   `json:"status"`
	Id     string `json:"id"`
	Info   string `json:"info"`
}

type ServiceMessage struct {
	Id   string   `json:"id"`
	Name string   `json:"name"`
	Ip   string   `json:"ip"`
	Port int      `json:"port"`
	Args []string `json:"args"`
	Ret  []string `json:"ret"`
}

func (sm ServiceMessage) ToJson() []byte {
	js, err := json.Marshal(sm)
	if err != nil {
		panic(err)
	}
	return js
}
