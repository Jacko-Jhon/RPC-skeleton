package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"time"
)

func encode(n int) []byte {
	code := make([]byte, 4)
	code[0] = byte(n >> 24)
	code[1] = byte(n >> 16)
	code[2] = byte(n >> 8)
	code[3] = byte(n)
	return code
}

func decode(code []byte) int {
	return int(code[0])<<24 | int(code[1])<<16 | int(code[2])<<8 | int(code[3])
}

type Handler struct {
	ServiceRegistry
	ServiceProvider
	addr   string
	socket *net.UDPConn //We may use SO_REUSEPORT later
}

// NewHandler @addr: ip:port default = "0.0.0.0:8888"
func NewHandler(addr ...string) *Handler {
	_addr := "0.0.0.0:8888"
	if len(addr) > 0 {
		_addr = addr[0]
	}
	Addr, err := net.ResolveUDPAddr("udp", _addr)
	if err != nil {
		panic(err)
	}
	sk, err := net.ListenUDP("udp", Addr)
	if err != nil {
		panic(err)
	}
	return &Handler{
		ServiceRegistry: *NewServiceRegistry(),
		ServiceProvider: *NewServiceProvider(),
		addr:            _addr,
		socket:          sk,
	}
}

func (hd Handler) SendResponse(res []byte, addr *net.UDPAddr, seq int) {
	prevH := []byte("REGIST\r\nACK:")
	ack := encode(seq + 1)
	h := append(prevH, ack...)
	msg := append(h, res...)
	for i := 0; i < 3; i++ {
		_, err := hd.socket.WriteToUDP(msg, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// SRegister Handle the register message
func (hd Handler) SRegister(name, ip string, port int, args []string, addr *net.UDPAddr, seq int) {
	msgJs := hd.ServiceRegistry.Register(name, ip, port, args).ToJson()
	hd.SendResponse(msgJs, addr, seq)
}

// SHeartbeat Maybe we don't need to respond to the heartbeat message
func (hd Handler) SHeartbeat(id string) {
	hd.ServiceRegistry.Heartbeat(id)
}

func (hd Handler) SRegisterByName(id, name, ip string, port int, addr *net.UDPAddr, seq int) {
	msgJs := hd.ServiceRegistry.RegisterByName(id, name, ip, port).ToJson()
	hd.SendResponse(msgJs, addr, seq)
}

func (hd Handler) SUnregister(id string, addr *net.UDPAddr, seq int) {
	msgJs := hd.ServiceRegistry.Unregister(id).ToJson()
	hd.SendResponse(msgJs, addr, seq)
}

func (hd Handler) SUnregisterAll(id, name string, addr *net.UDPAddr, seq int) {
	msgJs := hd.ServiceRegistry.UnregisterAll(id, name).ToJson()
	hd.SendResponse(msgJs, addr, seq)
}

func (hd Handler) SUpdateUrl(id, ip string, port int, addr *net.UDPAddr, seq int) {
	msgJs := hd.ServiceRegistry.UpdateUrl(id, ip, port).ToJson()
	hd.SendResponse(msgJs, addr, seq)
}

func (hd Handler) HandleServer(msg []byte, addr *net.UDPAddr) {
	seqByte := msg[12:16]
	seq := decode(seqByte)
	opCode := int(msg[21] - 48)
	var sm ServiceMessage
	err := json.Unmarshal(msg[22:], &sm)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch opCode {
	case 0:
		hd.Heartbeat(sm.id)
	case 1:
		hd.SRegister(sm.name, sm.ip, sm.port, sm.args, addr, seq)
	case 2:
		hd.SRegisterByName(sm.id, sm.name, sm.ip, sm.port, addr, seq)
	case 3:
		hd.SUnregister(sm.id, addr, seq)
	case 4:
		hd.SUnregisterAll(sm.id, sm.name, addr, seq)
	case 5:
		hd.SUpdateUrl(sm.id, sm.ip, sm.port, addr, seq)
	}
}

func (hd Handler) CRequestService(name string, addr *net.UDPAddr, seq int) {
	msgJs := hd.ServiceProvider.RequestService(name).ToJson()
	hd.SendResponse(msgJs, addr, seq)
}

// CGetServiceList 处理客户端请求服务列表
// （由于UDP限制，所以支持的服务数量有限，后续可以通过限制服务名称长度以及更换成TCP解决）
func (hd Handler) CGetServiceList(addr *net.UDPAddr, seq int) {
	msgJs := hd.ServiceProvider.GetServiceList().ToJson()
	hd.SendResponse(msgJs, addr, seq)
}

func (hd Handler) HandleClient(msg []byte, addr *net.UDPAddr) {
	seqByte := msg[12:16]
	seq := decode(seqByte)
	opCode := int(msg[21] - 48)
	var cm ClientMessage
	err := json.Unmarshal(msg[22:], &cm)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch opCode {
	case 0:
		hd.CRequestService(cm.name, addr, seq)
	case 1:
		hd.CGetServiceList(addr, seq)
	}
}

func (hd Handler) HealthChecker() {
	dt := hd.GetSRLiveTime()
	for {
		hd.ServiceRegistry.CheckHealth()
		time.Sleep(time.Duration(dt) * time.Second)
	}
}

func (hd Handler) run() {
	go hd.HealthChecker()
	for {
		msg := make([]byte, 1024)
		n, addr, err := hd.socket.ReadFromUDP(msg)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n < 22 {
			continue
		}
		if string(msg[:6]) == "SERVER" {
			go hd.HandleServer(msg, addr)
		} else if string(msg[:6]) == "CLIENT" {
			go hd.HandleClient(msg, addr)
		}
	}
}

func main() {
	var DPath string
	var LPath string
	var slt int64
	flag.StringVar(&LPath, "l", "", "To load json from file (C:/service.json)")
	flag.StringVar(&DPath, "d", "", "Dump json to file (C:/service.json) later")
	flag.Int64Var(&slt, "slt", 60, "Set live time for server")
	//解析命令行参数
	flag.Parse()
	hd := NewHandler()
	if LPath != "" {
		GlobalRegistry.load(LPath)
	}
	if DPath != "" {
		defer GlobalRegistry.dump(DPath)
	}
	hd.SetSRLiveTime(slt)
	hd.run()
	defer func(socket *net.UDPConn) {
		err := socket.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(hd.socket)
}
