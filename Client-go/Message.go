package main

type ServiceList struct {
	Status bool     `json:"status"`
	Info   string   `json:"info"`
	List   []string `json:"list"`
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

type ServiceBrief struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
	Ret  []string `json:"ret"`
}
