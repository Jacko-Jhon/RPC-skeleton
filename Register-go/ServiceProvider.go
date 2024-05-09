package Register_go

import (
	"encoding/json"
	"time"
)

type ServiceProvider struct {
	Registry
	liveTime int64
}

func (sp ServiceProvider) SetLiveTime(T int64) {
	sp.liveTime = T
}

func (sp ServiceProvider) GetLiveTime() int64 {
	return sp.liveTime
}

func (sp ServiceProvider) getService(name string) []ServiceInfo {
	_, ok := GlobalRegistry.nameToId[name]
	if !ok {
		return nil
	} else {
		idList := GlobalRegistry.nameToId[name]
		List := make([]ServiceInfo, 0, len(idList))
		for _, id := range idList {
			if GlobalRegistry.services[id].Status == 0 {
				continue
			}
			List = append(List, GlobalRegistry.services[id])
		}
		return List
	}
}

func (sp ServiceProvider) FindService(name string) MessageToClient {
	_, ok := GlobalRegistry.nameToId[name]
	if !ok {
		return MessageToClient{status: false, info: "service not found", json: nil}
	} else {
		List := sp.getService(name)
		if len(List) == 0 {
			return MessageToClient{status: false, info: "service offline", json: nil}
		}
		serviceList, err := json.Marshal(List)
		if err != nil {
			return MessageToClient{status: false, info: "unknown error", json: nil}
		}
		return MessageToClient{status: true, info: "service found", json: serviceList}
	}
}

func (sp ServiceProvider) GetServiceList() MessageToClient {
	List := make([]string, 0, len(GlobalRegistry.nameToId))
	for name := range GlobalRegistry.nameToId {
		List = append(List, name)
	}
	serviceList, err := json.Marshal(List)
	if err != nil {
		return MessageToClient{status: false, info: "unknown error", json: nil}
	}
	return MessageToClient{status: true, info: "service list found", json: serviceList}
}

func (sp ServiceProvider) RequestService(id, name string) MessageToClient {
	_, ok := GlobalRegistry.nameToId[name]
	if !ok {
		return MessageToClient{status: false, info: "service not found", json: nil}
	}
	List := sp.getService(name)
	if len(List) == 0 {
		return MessageToClient{status: false, info: "service offline", json: nil}
	}
	serviceList, err := json.Marshal(List)
	if err != nil {
		return MessageToClient{status: false, info: "unknown error", json: nil}
	}
	if id == "" {
		id = IdGenerate()
		GlobalRegistry.clientsWithService[id] = make([]string, 0)
	}
	GlobalRegistry.clientsWithService[id] = append(GlobalRegistry.clientsWithService[id], name)
	return MessageToClient{status: true, id: id, info: "service requested", json: serviceList}
}

func (sp ServiceProvider) FetchServices(id string) MessageToClient {
	List, ok := GlobalRegistry.clientsWithService[id]
	if !ok {
		return MessageToClient{status: false, id: id, info: "client not found", json: nil}
	}
	SList := make(map[string][]ServiceInfo)
	for _, name := range List {
		SList[name] = sp.getService(name)
	}
	serviceList, err := json.Marshal(SList)
	if err != nil {
		return MessageToClient{status: false, id: id, info: "unknown error", json: nil}
	}
	GlobalRegistry.clients[id] = time.Now().Unix()
	return MessageToClient{status: true, id: id, info: "fetch success", json: serviceList}
}

func (sp ServiceProvider) CheckClient() {
	T := time.Now().Unix()
	for k, v := range GlobalRegistry.clients {
		if T > v+sp.liveTime {
			delete(GlobalRegistry.clients, k)
			delete(GlobalRegistry.clientsWithService, k)
		}
	}
}
