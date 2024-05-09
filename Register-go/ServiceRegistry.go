package Register_go

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

type ServiceRegister struct {
	Register
	liveTime int64
}

func (sr ServiceRegister) SetLiveTime(t int64) {
	sr.liveTime = t
}

func (sr ServiceRegister) GetLiveTime() int64 {
	return sr.liveTime
}

func (sr ServiceRegister) Register(id, name, ip string, port int, args []string) MessageToServer {
	_, ok := GlobalRegister.nameToId[name]
	if ok {
		return MessageToServer{status: false, id: id, info: "name already exist"}
	} else {
		GlobalRegister.services[id] = *NewServiceInfo(id, name, ip, port, args)
		GlobalRegister.nameToId[name] = append(GlobalRegister.nameToId[name], id)
		return MessageToServer{status: true, id: id, info: "registry success"}
	}
}

func (sr ServiceRegister) UpdateUrl(id, ip string, port int) MessageToServer {
	_, ok := GlobalRegister.services[id]
	if ok {
		GlobalRegister.services[id].UpdateUrl(ip, port)
		return MessageToServer{status: true, id: id, info: "update url success"}
	} else {
		return MessageToServer{status: false, id: id, info: "id not exist"}
	}
}

func (sr ServiceRegister) UpdateArgs(id string, args []string) MessageToServer {
	_, ok := GlobalRegister.services[id]
	if ok {
		GlobalRegister.services[id].UpdateArgs(args)
		return MessageToServer{status: true, id: id, info: "update args success"}
	} else {
		return MessageToServer{status: false, id: id, info: "id not exist"}
	}
}

func (sr ServiceRegister) HeartBeat(id string) MessageToServer {
	_, ok := GlobalRegister.services[id]
	if ok {
		GlobalRegister.services[id].HeartBeat()
		return MessageToServer{status: true, id: id, info: "heartbeat success"}
	} else {
		return MessageToServer{status: false, id: id, info: "id not exist"}
	}
}

func (sr ServiceRegister) UnRegister(id string) MessageToServer {
	_, ok := GlobalRegister.services[id]
	if ok {
		name := GlobalRegister.services[id].Name
		removeElement(GlobalRegister.nameToId[name], id)
		delete(GlobalRegister.services, id)
		return MessageToServer{status: true, id: id, info: "unregister success"}
	} else {
		return MessageToServer{status: false, id: id, info: "id not exist"}
	}
}
