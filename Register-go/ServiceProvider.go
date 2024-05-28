package main

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
func (sp *ServiceProvider) getService(name string) []ServiceInfo {
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
func (sp *ServiceProvider) RequestService(name string) *ServiceUrls {
	_, ok := GlobalRegistry.nameToId[name]
	if !ok {
		return &ServiceUrls{Status: false, Info: "service not found"}
	} else {
		List := sp.getService(name)
		if len(List) == 0 {
			return &ServiceUrls{Status: false, Info: "server offline"}
		}
		ips := make([]string, 0, len(List))
		ports := make([]int, 0, len(List))
		statusL := make([]int32, 0, len(List))
		for _, it := range List {
			ips = append(ips, it.Ip)
			ports = append(ports, it.Port)
			statusL = append(statusL, it.Status)
		}
		return &ServiceUrls{
			Status:  true,
			Info:    "service found",
			Name:    name,
			Ips:     ips,
			Ports:   ports,
			Factors: statusL,
			Args:    List[0].Args,
			Ret:     List[0].Ret,
		}
	}
}

// GetServiceList 获取所有服务列表
func (sp *ServiceProvider) GetServiceList() ServiceList {
	List := make([]string, 0, len(GlobalRegistry.nameToId))
	for name := range GlobalRegistry.nameToId {
		List = append(List, name)
	}
	if len(List) == 0 {
		return ServiceList{Status: false, Info: "service list not found"}
	}
	return ServiceList{Status: true, Info: "service list found", List: List}
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
