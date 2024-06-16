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
	return &ServiceProvider{}
}

// getService 获取在线服务器列表
func (sp *ServiceProvider) getService(name string) []ServiceInfo {
	GlobalRegistry.lock.RLock()
	// 获取相应服务名称的服务id
	idList, ok := GlobalRegistry.nameToId[name]
	if !ok {
		GlobalRegistry.lock.RUnlock()
		// 不存在服务返回 nil
		return nil
	} else {
		List := make([]ServiceInfo, 0, len(*idList))
		for _, id := range *idList {
			sv, _ := GlobalRegistry.services.Load(id)
			// 如果该服务不在线则不加入列表
			if sv.(*ServiceInfo).Status == 0 {
				continue
			}
			List = append(List, *sv.(*ServiceInfo))
		}
		GlobalRegistry.lock.RUnlock()
		// 返回列表
		return List
	}
}

// RequestService 请求服务
func (sp *ServiceProvider) RequestService(name string) *ServiceUrls {
	GlobalRegistry.lock.RLock()
	_, ok := GlobalRegistry.nameToId[name]
	GlobalRegistry.lock.RUnlock()
	if !ok {
		return &ServiceUrls{Status: false, Info: "service not found"}
	} else {
		// 获取在线服务器列表
		List := sp.getService(name)
		if len(List) == 0 {
			return &ServiceUrls{Status: false, Info: "server offline"}
		}
		// 填写 ServiceUrls 信息
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
	GlobalRegistry.lock.RLock()
	List := make([]string, 0, len(GlobalRegistry.nameToId))
	// 获取所有服务名称
	for name := range GlobalRegistry.nameToId {
		List = append(List, name)
	}
	GlobalRegistry.lock.RUnlock()
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
