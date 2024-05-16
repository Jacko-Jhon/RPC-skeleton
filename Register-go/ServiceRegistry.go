package main

import "time"

// removeElement 从字符串切片中移除指定的元素
// 参数：
//
//	slice：要操作的字符串切片
//	id：要移除的元素的标识符
//
// 说明：该函数会遍历切片，找到第一个与 id 相等的元素并将其移除。如果找到多个相等元素，只移除第一个。
func removeElement(slice []string, id string) {
	idx := -1
	for i := 0; i < len(slice); i++ {
		if slice[i] == id {
			idx = i
		}
	}
	if idx != -1 {
		slice = append(slice[:idx], slice[idx+1:]...)
	}
}

type ServiceRegistry struct {
	Registry
	liveTime int64
}

func (sr ServiceRegistry) SetSRLiveTime(t int64) {
	sr.liveTime = t
}

func (sr ServiceRegistry) GetSRLiveTime() int64 {
	return sr.liveTime
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		Registry: *GlobalRegistry,
		liveTime: 60,
	}
}

// Register 在服务注册表中注册一个新的服务。
// 参数：
// - name: 服务的名称，必须唯一。
// - ip: 服务的IP地址。
// - port: 服务监听的端口号。
// - args: 传递给服务的额外参数。
// 返回值：
// - MessageToServer: 包含注册结果的消息，包括是否成功、服务ID和相关信息。
func (sr ServiceRegistry) Register(name, ip string, port int, args []string) MessageToServer {
	_, ok := GlobalRegistry.nameToId[name]
	if ok {
		return MessageToServer{status: false, id: "", info: "name already exist"}
	} else {
		id := IdGenerate()
		GlobalRegistry.services[id] = *NewServiceInfo(id, name, ip, port, args)
		GlobalRegistry.nameToId[name] = append(GlobalRegistry.nameToId[name], id)
		return MessageToServer{status: true, id: id, info: "registry success"}
	}
}

// UpdateUrl 根据服务ID更新服务的URL。
func (sr ServiceRegistry) UpdateUrl(id, ip string, port int) MessageToServer {
	_, ok := GlobalRegistry.services[id]
	if ok {
		GlobalRegistry.services[id].UpdateUrl(ip, port)
		return MessageToServer{status: true, id: id, info: "update url success"}
	} else {
		return MessageToServer{status: false, id: id, info: "id not exist"}
	}
}

// UpdateArgs 根据服务ID更新服务的参数。
func (sr ServiceRegistry) UpdateArgs(id string, args []string) MessageToServer {
	_, ok := GlobalRegistry.services[id]
	if ok {
		GlobalRegistry.services[id].UpdateArgs(args)
		return MessageToServer{status: true, id: id, info: "update args success"}
	} else {
		return MessageToServer{status: false, id: id, info: "id not exist"}
	}
}

// Heartbeat 根据服务ID更新服务的心跳时间。
func (sr ServiceRegistry) Heartbeat(id string) MessageToServer {
	_, ok := GlobalRegistry.services[id]
	if ok {
		GlobalRegistry.services[id].HeartBeat()
		return MessageToServer{status: true, id: id, info: "heartbeat success"}
	} else {
		return MessageToServer{status: false, id: id, info: "id not exist"}
	}
}

// Unregister 根据服务ID注销服务。
func (sr ServiceRegistry) Unregister(id string) MessageToServer {
	_, ok := GlobalRegistry.services[id]
	if ok {
		name := GlobalRegistry.services[id].Name
		removeElement(GlobalRegistry.nameToId[name], id)
		delete(GlobalRegistry.services, id)
		return MessageToServer{status: true, id: id, info: "unregister success"}
	} else {
		return MessageToServer{status: false, id: id, info: "id not exist"}
	}
}

// UnregisterAll 根据服务名称注销所有同名服务。
func (sr ServiceRegistry) UnregisterAll(id, name string) MessageToServer {
	ids, ok := GlobalRegistry.nameToId[name]
	if ok {
		flag := true
		for _, it := range ids {
			if it == id {
				flag = false
				break
			}
		}
		if flag {
			return MessageToServer{status: false, id: "", info: "id not exist"}
		}
		for _, id := range GlobalRegistry.nameToId[name] {
			delete(GlobalRegistry.services, id)
		}
		delete(GlobalRegistry.nameToId, name)
		return MessageToServer{status: true, id: "", info: "unregister all success"}
	} else {
		return MessageToServer{status: false, id: "", info: "name not exist"}
	}
}

// CheckHealth 检查服务心跳，更新服务状态。
func (sr ServiceRegistry) CheckHealth() {
	T := time.Now().Unix()
	for _, service := range GlobalRegistry.services {
		if service.IsAlive(T, sr.liveTime) {
			if service.Status > 1 {
				service.Status--
			}
		} else {
			service.Status = 0
		}
	}
}

// RegisterByName 根据服务名称注册服务。（需要用已有的id进行验证）
func (sr ServiceRegistry) RegisterByName(id, name, ip string, port int) MessageToServer {
	nameList, ok := GlobalRegistry.nameToId[name]
	flag := true
	if ok {
		for _, _id := range nameList {
			if id == _id {
				flag = false
				break
			}
		}
		if flag {
			return MessageToServer{status: false, id: id, info: "authentication failed"}
		}
		nid := IdGenerate()
		GlobalRegistry.services[nid] = *NewServiceInfo(id, name, ip, port, GlobalRegistry.services[id].Args)
		GlobalRegistry.nameToId[name] = append(GlobalRegistry.nameToId[name], id)
		return MessageToServer{status: true, id: nid, info: "registry success"}
	} else {
		return MessageToServer{status: false, id: "", info: "name not exist"}
	}
}
