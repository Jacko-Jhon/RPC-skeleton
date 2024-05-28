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

type ServiceList struct {
	Status bool     `json:"status"`
	Info   string   `json:"info"`
	List   []string `json:"list"`
}

func (sl ServiceList) ToJson() []byte {
	js, err := json.Marshal(sl)
	if err != nil {
		panic(err)
	}
	return js
}

type ServiceUrls struct {
	Status  bool     `json:"status"`
	Info    string   `json:"info"`
	Name    string   `json:"name"`
	Ips     []string `json:"ips"`
	Ports   []int    `json:"ports"`
	Factors []int32  `json:"factors"`
	Args    []string `json:"args"`
	Ret     []string `json:"ret"`
}

func (su ServiceUrls) ToJson() []byte {
	js, err := json.Marshal(su)
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
