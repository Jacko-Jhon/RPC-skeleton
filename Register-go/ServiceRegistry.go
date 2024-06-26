package main

import (
	"sync/atomic"
	"time"
)

// removeElement 从字符串切片中移除指定的元素
// 参数：
//
//	slice：要操作的字符串切片
//	id：要移除的元素的标识符
//
// 说明：该函数会遍历切片，找到第一个与 id 相等的元素并将其移除。如果找到多个相等元素，只移除第一个。
func removeElement(slice *[]string, id string) int {
	idx := -1
	for i := 0; i < len(*slice); i++ {
		if (*slice)[i] == id {
			idx = i
		}
	}
	if idx != -1 {
		*slice = append((*slice)[:idx], (*slice)[idx+1:]...)
	}
	return len(*slice)
}

type ServiceRegistry struct {
	liveTime int64
}

func (sr *ServiceRegistry) SetSRLiveTime(t int64) {
	sr.liveTime = t
}

func (sr *ServiceRegistry) GetSRLiveTime() int64 {
	return sr.liveTime
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
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
func (sr *ServiceRegistry) Register(name, ip string, port int, args []string, ret []string) MessageToServer {
	GlobalRegistry.lock.RLock()
	// 查询服务是否存在
	_, ok := GlobalRegistry.nameToId[name]
	GlobalRegistry.lock.RUnlock()
	if ok {
		return MessageToServer{Status: false, Id: "", Info: "name already exist"}
	} else {
		GlobalRegistry.lock.Lock()
		// 再次确认服务是否存在
		_, ok := GlobalRegistry.nameToId[name]
		if ok {
			GlobalRegistry.lock.Unlock()
			return MessageToServer{Status: false, Id: "", Info: "name already exist"}
		}
		// 生成服务ID
		id := IdGenerate()
		// 存储服务信息
		GlobalRegistry.services.Store(id, NewServiceInfo(id, name, ip, port, args, ret))
		atomic.AddInt32(&GlobalRegistry.count, 1)
		GlobalRegistry.nameToId[name] = &[]string{id}
		GlobalRegistry.lock.Unlock()
		return MessageToServer{Status: true, Id: id, Info: "registry success"}
	}
}

// UpdateUrl 根据服务ID更新服务的URL。
func (sr *ServiceRegistry) UpdateUrl(id, ip string, port int) MessageToServer {
	sv, ok := GlobalRegistry.services.Load(id)
	if ok {
		sv.(*ServiceInfo).UpdateUrl(ip, port)
		return MessageToServer{Status: true, Id: id, Info: "update url success"}
	} else {
		return MessageToServer{Status: false, Id: id, Info: "id not exist"}
	}
}

// UpdateArgs 根据服务ID更新服务的参数。
//func (sr *ServiceRegistry) UpdateArgs(id string, args []string) MessageToServer {
//	_, ok := GlobalRegistry.services[id]
//	if ok {
//		service := GlobalRegistry.services[id]
//		service.UpdateArgs(args)
//		return MessageToServer{Status: true, Id: id, Info: "update args success"}
//	} else {
//		return MessageToServer{Status: false, Id: id, Info: "id not exist"}
//	}
//}

// Heartbeat 根据服务ID更新服务的心跳时间。
func (sr *ServiceRegistry) Heartbeat(id string) {
	sv, ok := GlobalRegistry.services.Load(id)
	if ok {
		sv.(*ServiceInfo).HeartBeat()
	}
}

// Unregister 根据服务ID注销服务。
func (sr *ServiceRegistry) Unregister(id string) MessageToServer {
	// 查找服务是否存在
	sv, ok := GlobalRegistry.services.Load(id)
	if ok {
		name := sv.(*ServiceInfo).Name
		GlobalRegistry.lock.Lock()
		// 获取同名服务列表
		list := GlobalRegistry.nameToId[name]
		// 删除列表中的该服务ID，如果只有该ID，删除整个列表
		if removeElement(list, id) == 0 {
			delete(GlobalRegistry.nameToId, name)
		}
		GlobalRegistry.lock.Unlock()
		// 删除服务
		GlobalRegistry.services.Delete(id)
		return MessageToServer{Status: true, Id: id, Info: "unregister success"}
	} else {
		return MessageToServer{Status: false, Id: id, Info: "id not exist"}
	}
}

// UnregisterAll 根据服务名称注销所有同名服务。
func (sr *ServiceRegistry) UnregisterAll(id, name string) MessageToServer {
	GlobalRegistry.lock.RLock()
	ids, ok := GlobalRegistry.nameToId[name]
	if ok {
		flag := true
		for _, it := range *ids {
			if it == id {
				flag = false
				break
			}
		}
		if flag {
			GlobalRegistry.lock.RUnlock()
			return MessageToServer{Status: false, Id: "", Info: "id not exist"}
		}
		GlobalRegistry.lock.RUnlock()
		GlobalRegistry.lock.Lock()
		ids, ok = GlobalRegistry.nameToId[name]
		if !ok {
			GlobalRegistry.lock.Unlock()
			return MessageToServer{Status: false, Id: "", Info: "name not exist"}
		}
		for _, it := range *ids {
			GlobalRegistry.services.Delete(it)
		}
		delete(GlobalRegistry.nameToId, name)
		GlobalRegistry.lock.Unlock()
		return MessageToServer{Status: true, Id: "", Info: "unregister all success"}
	} else {
		GlobalRegistry.lock.RUnlock()
		return MessageToServer{Status: false, Id: "", Info: "name not exist"}
	}
}

// CheckHealth 检查服务心跳，更新服务状态。
func (sr *ServiceRegistry) CheckHealth() {
	// 获取当前时间戳
	T := time.Now().Unix()
	// 遍历服务列表，更新服务状态
	GlobalRegistry.services.Range(func(key, value interface{}) bool {
		service := value.(*ServiceInfo)
		if service.IsAlive(T, sr.liveTime) {
			s := service.Status
			if s > 2 && s < 100 {
				atomic.AddInt32(&service.Status, -2)
			} else if s >= 100 {
				atomic.StoreInt32(&service.Status, 100)
			}
		} else {
			atomic.StoreInt32(&service.Status, 0)
		}
		return true
	})
}

// RegisterByName 根据服务名称注册服务。（需要用已有的id进行验证）
func (sr *ServiceRegistry) RegisterByName(id, name, ip string, port int) MessageToServer {
	// 验证id和name是否存在
	sv, ok := GlobalRegistry.services.Load(id)
	if ok && sv.(*ServiceInfo).Name == name {
		// 生成新的id
		nid := IdGenerate()
		GlobalRegistry.lock.Lock()
		// 获取同名服务列表
		nameList, ok := GlobalRegistry.nameToId[name]
		if !ok {
			// 如果同名服务列表不存在，则创建新的id列表 （健壮性考虑）
			GlobalRegistry.nameToId[name] = &[]string{id}
		} else {
			// 如果同名服务列表存在，则添加新的id
			*nameList = append(*nameList, id)
		}
		GlobalRegistry.lock.Unlock()
		// 存储服务信息
		GlobalRegistry.services.Store(nid, NewServiceInfo(id, name, ip, port, sv.(*ServiceInfo).Args, sv.(*ServiceInfo).Ret))
		return MessageToServer{Status: true, Id: nid, Info: "registry success"}
	} else {
		GlobalRegistry.lock.RUnlock()
		return MessageToServer{Status: false, Id: "", Info: "id or name not exist"}
	}
}
