package Register_go

import "time"

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

func (sr ServiceRegistry) SetLiveTime(t int64) {
	sr.liveTime = t
}

func (sr ServiceRegistry) GetLiveTime() int64 {
	return sr.liveTime
}

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

func (sr ServiceRegistry) UpdateUrl(id, ip string, port int) MessageToServer {
	_, ok := GlobalRegistry.services[id]
	if ok {
		GlobalRegistry.services[id].UpdateUrl(ip, port)
		return MessageToServer{status: true, id: id, info: "update url success"}
	} else {
		return MessageToServer{status: false, id: id, info: "id not exist"}
	}
}

func (sr ServiceRegistry) UpdateArgs(id string, args []string) MessageToServer {
	_, ok := GlobalRegistry.services[id]
	if ok {
		GlobalRegistry.services[id].UpdateArgs(args)
		return MessageToServer{status: true, id: id, info: "update args success"}
	} else {
		return MessageToServer{status: false, id: id, info: "id not exist"}
	}
}

func (sr ServiceRegistry) HeartBeat(id string) MessageToServer {
	_, ok := GlobalRegistry.services[id]
	if ok {
		GlobalRegistry.services[id].HeartBeat()
		return MessageToServer{status: true, id: id, info: "heartbeat success"}
	} else {
		return MessageToServer{status: false, id: id, info: "id not exist"}
	}
}

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

func (sr ServiceRegistry) UnregisterAll(name string) MessageToServer {
	_, ok := GlobalRegistry.nameToId[name]
	if ok {
		for _, id := range GlobalRegistry.nameToId[name] {
			delete(GlobalRegistry.services, id)
		}
		delete(GlobalRegistry.nameToId, name)
		return MessageToServer{status: true, id: "", info: "unregister all success"}
	} else {
		return MessageToServer{status: false, id: "", info: "name not exist"}
	}
}

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

func (sr ServiceRegistry) RegisterByName(name, ip string, port int, args []string) MessageToServer {
	_, ok := GlobalRegistry.nameToId[name]
	if ok {
		id := IdGenerate()
		GlobalRegistry.services[id] = *NewServiceInfo(id, name, ip, port, args)
		GlobalRegistry.nameToId[name] = append(GlobalRegistry.nameToId[name], id)
		return MessageToServer{status: true, id: id, info: "registry success"}
	} else {
		return MessageToServer{status: false, id: "", info: "name not exist"}
	}
}
