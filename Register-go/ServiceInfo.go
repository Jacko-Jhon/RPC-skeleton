package main

import (
	"sync/atomic"
	"time"
)

type ServiceInfo struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Ip        string   `json:"ip"`
	Port      int      `json:"port"`
	Status    int32    `json:"status"`
	Heartbeat int64    `json:"heartbeat"`
	Args      []string `json:"args"`
	Ret       []string `json:"ret"`
}

func NewServiceInfo(id, name, ip string, port int, args []string, ret []string) *ServiceInfo {
	Args := make([]string, len(args), cap(args))
	copy(Args, args)
	return &ServiceInfo{
		Id:        id,
		Name:      name,
		Ip:        ip,
		Port:      port,
		Status:    0,
		Heartbeat: 0,
		Args:      Args,
		Ret:       ret,
	}
}

func (si *ServiceInfo) UpdateUrl(ip string, port int) {
	si.Ip = ip
	si.Port = port
}

func (si *ServiceInfo) HeartBeat() {
	atomic.StoreInt64(&si.Heartbeat, time.Now().Unix())
	atomic.AddInt32(&si.Status, 1)
}

func (si *ServiceInfo) IsAlive(T, liveTime int64) bool {
	return T-si.Heartbeat < liveTime
}

func (si *ServiceInfo) UpdateArgs(args []string) {
	si.Args = make([]string, len(args), cap(args))
	copy(si.Args, args)
}
