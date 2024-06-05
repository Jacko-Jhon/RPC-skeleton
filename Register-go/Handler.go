package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"time"
)

var isQuit = false

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
func NewHandler(_addr string) *Handler {
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

func (hd *Handler) SendResponse(res []byte, addr *net.UDPAddr, seq int) {
	prevH := []byte("REGIST\r\nACK:")
	ack := encode(seq)
	h := append(prevH, ack...)
	msg := append(h, res...)
	_, err := hd.socket.WriteToUDP(msg, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// SRegister Handle the register message
func (hd *Handler) SRegister(name, ip string, port int, args []string, ret []string, addr *net.UDPAddr, seq int) {
	if ip == "0.0.0.0" {
		ip = addr.IP.String()
	}
	msgJs := hd.ServiceRegistry.Register(name, ip, port, args, ret).ToJson()
	hd.SendResponse(msgJs, addr, seq)
}

// SHeartbeat Maybe we don't need to respond to the heartbeat message
func (hd *Handler) SHeartbeat(id string) {
	hd.ServiceRegistry.Heartbeat(id)
}

func (hd *Handler) SRegisterByName(id, name, ip string, port int, addr *net.UDPAddr, seq int) {
	if ip == "0.0.0.0" {
		ip = addr.IP.String()
	}
	msgJs := hd.ServiceRegistry.RegisterByName(id, name, ip, port).ToJson()
	hd.SendResponse(msgJs, addr, seq)
}

func (hd *Handler) SUnregister(id string, addr *net.UDPAddr, seq int) {
	msgJs := hd.ServiceRegistry.Unregister(id).ToJson()
	hd.SendResponse(msgJs, addr, seq)
}

func (hd *Handler) SUnregisterAll(id, name string, addr *net.UDPAddr, seq int) {
	msgJs := hd.ServiceRegistry.UnregisterAll(id, name).ToJson()
	hd.SendResponse(msgJs, addr, seq)
}

func (hd *Handler) SUpdateUrl(id, ip string, port int, addr *net.UDPAddr, seq int) {
	if ip == "0.0.0.0" {
		ip = addr.IP.String()
	}
	msgJs := hd.ServiceRegistry.UpdateUrl(id, ip, port).ToJson()
	hd.SendResponse(msgJs, addr, seq)
}

func (hd *Handler) HandleServer(msg []byte, addr *net.UDPAddr, seq int) {
	opCode := int(msg[21] - 48)
	var sm *ServiceMessage = nil
	err := json.Unmarshal(msg[22:], &sm)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch opCode {
	case 0:
		hd.SHeartbeat(sm.Id)
	case 1:
		hd.SRegister(sm.Name, sm.Ip, sm.Port, sm.Args, sm.Ret, addr, seq)
	case 2:
		hd.SRegisterByName(sm.Id, sm.Name, sm.Ip, sm.Port, addr, seq)
	case 3:
		hd.SUnregister(sm.Id, addr, seq)
	case 4:
		hd.SUnregisterAll(sm.Id, sm.Name, addr, seq)
	case 5:
		hd.SUpdateUrl(sm.Id, sm.Ip, sm.Port, addr, seq)
	}
}

func (hd *Handler) CRequestService(name string, addr *net.UDPAddr, seq int) {
	msgJs := hd.ServiceProvider.RequestService(name).ToJson()
	hd.SendResponse(msgJs, addr, seq)
}

// CGetServiceList 处理客户端请求服务列表
// （由于UDP限制，所以支持的服务数量有限，后续可以通过限制服务名称长度以及更换成TCP解决）
func (hd *Handler) CGetServiceList(addr *net.UDPAddr, seq int) {
	msgJs := hd.ServiceProvider.GetServiceList().ToJson()
	hd.SendResponse(msgJs, addr, seq)
}

func (hd *Handler) HandleClient(msg []byte, addr *net.UDPAddr, seq int) {
	opCode := int(msg[21] - 48)
	switch opCode {
	case 0:
		var name string
		if string(msg[24:29]) == "NAME:" {
			name = string(msg[29:])
		}
		hd.CRequestService(name, addr, seq)
	case 1:
		hd.CGetServiceList(addr, seq)
	}
}

func (hd *Handler) HealthChecker() {
	dt := hd.GetSRLiveTime()
	for {
		hd.ServiceRegistry.CheckHealth()
		time.Sleep(time.Duration(dt) * time.Second)
	}
}

func (hd *Handler) run() {
	go hd.HealthChecker()
	for {
		msg := make([]byte, 1024)
		n, addr, err := hd.socket.ReadFromUDP(msg)
		if err != nil {
			if isQuit {
				return
			}
			fmt.Println(err)
			continue
		}
		seqByte := msg[12:16]
		seq := decode(seqByte)
		if n < 22 {
			continue
		}
		if string(msg[:6]) == "SERVER" {
			go hd.HandleServer(msg[0:n], addr, seq)
		} else if string(msg[:6]) == "CLIENT" {
			go hd.HandleClient(msg[0:n], addr, seq)
		}
	}
}

func main() {
	var DPath string
	var LPath string
	var ListenAddr string
	var slt int64
	flag.StringVar(&LPath, "load", "./Services.json", "To load json from file")
	flag.StringVar(&DPath, "dump", "./Services.json", "Dump json to file later")
	flag.StringVar(&ListenAddr, "l", "0.0.0.0:8888", "Set listen address")
	flag.Int64Var(&slt, "slt", 120, "Set live time for server")
	//解析命令行参数
	flag.Parse()
	GlobalRegistry.load(LPath)
	hd := NewHandler(ListenAddr)
	defer func(socket *net.UDPConn) {
		err := socket.Close()
		if err != nil {
			fmt.Println(err)
		}
		GlobalRegistry.dump(DPath)
	}(hd.socket)
	hd.SetSRLiveTime(slt)
	go hd.run()
	defer fmt.Println("Registry is closed")
	for {
		var input string
		_, err := fmt.Scanf("%s\n", &input)
		if err != nil {
			fmt.Println(err)
		}
		if input == "Quit" {
			isQuit = true
			return
		}
	}
}
