package main

import (
	"strconv"
)

type ServiceProvider struct {
	Registry
}

//func (sp ServiceProvider) SetSPLiveTime(T int64) {
//	sp.liveTime = T
//}
//
//func (sp ServiceProvider) GetSPLiveTime() int64 {
//	return sp.liveTime
//}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{
		Registry: *GlobalRegistry,
	}
}

// getService 获取在线服务器列表
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

// RequestService 请求服务
func (sp ServiceProvider) RequestService(name string) MessageToClient {
	_, ok := GlobalRegistry.nameToId[name]
	if !ok {
		return MessageToClient{status: false, info: "service not found"}
	} else {
		List := sp.getService(name)
		if len(List) == 0 {
			return MessageToClient{status: false, info: "server offline"}
		}
		urls := make([]string, 0, len(List))
		statusL := make([]int, 0, len(List))
		for i := 0; i < len(List); i++ {
			urls = append(urls, List[i].Ip+":"+strconv.Itoa(List[i].Port))
			statusL = append(statusL, List[i].Status)
		}
		return MessageToClient{status: true, info: "service found", serverName: name, urlList: urls, statusList: statusL}
	}
}

// GetServiceList 获取所有服务列表
func (sp ServiceProvider) GetServiceList() MessageToClient {
	List := make([]string, 0, len(GlobalRegistry.nameToId))
	for name := range GlobalRegistry.nameToId {
		List = append(List, name)
	}
	return MessageToClient{status: true, info: "service list found", serverList: List}
}

// FetchServices 根据client id获取服务
//func (sp ServiceProvider) FetchServices(id string) MessageToClient {
//	List, ok := GlobalRegistry.clientsWithService[id]
//	if !ok {
//		return MessageToClient{status: false, id: id, info: "client not found"}
//	}
//	SList := make(map[string][]ServiceInfo)
//	for _, name := range List {
//		SList[name] = sp.getService(name)
//	}
//	GlobalRegistry.clients[id] = time.Now().Unix()
//	return MessageToClient{status: true, id: id, info: "fetch success", fetch: SList}
//}

// CheckClient 检查客户端是否离线
//func (sp ServiceProvider) CheckClient() {
//	T := time.Now().Unix()
//	for k, v := range GlobalRegistry.clients {
//		if T > v+sp.liveTime {
//			delete(GlobalRegistry.clients, k)
//			delete(GlobalRegistry.clientsWithService, k)
//		}
//	}
//}
