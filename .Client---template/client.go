package Client

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	buffSize      = 4096
	serviceName   = []string{"add", "sub", "mul", "div"}
	myServices    map[string]*unsafe.Pointer
	rIp           string
	rPort         int
	rSocket       *net.UDPConn
	rd            = rand.New(rand.NewSource(time.Now().Unix()))
	reqH          = []byte("CLIENT\r\nREQ\r\nSEQ:")
	runH          = []byte("CLIENT\r\nRUN\r\nSEQ:")
	serviceSocket = make(map[string]net.PacketConn)
	timeOUT       = 5000
)

func Init(RegisterIp string, RegisterPort int) {
	rIp = RegisterIp
	rPort = RegisterPort
	var err error
	rSocket, err = net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(rIp),
		Port: rPort,
	})
	if err != nil {
		panic(err)
	}
	for _, s := range serviceName {
		ptr := unsafe.Pointer(nil)
		myServices[s] = &ptr
		socket, err := net.ListenPacket("udp", "")
		if err != nil {
			panic(err)
		}
		serviceSocket[s] = socket
	}
	fetch()
}

func fetch() {
	var Map = make(map[string][]byte)
	for it := range myServices {
		Map[it] = dailRegistry(it)
	}
	for it, v := range Map {
		if v == nil {
			continue
		}
		var msg serviceUrls
		err := json.Unmarshal(v, &msg)
		if err != nil {
			fmt.Println(err)
			continue
		}
		var s int32 = 0
		for _, f := range msg.Factors {
			s += f
		}
		url := urls{
			s:       s,
			ips:     msg.Ips,
			ports:   msg.Ports,
			factors: msg.Factors,
		}
		atomic.StorePointer(myServices[it], unsafe.Pointer(&url))
	}
}

// Call
/*
This function is used to call a service by name

@ServiceName is the name of the service

@ArgsType is the arguments of the service, you can get the result from it.

for example: Client.Call("add", &AddArgs{A: 1, B: 2})

Note: don't use this function to call a service asynchronously
*/
func Call(ServiceName string, ArgsType interface{}) {
	switch ServiceName {
	case "add":
		args := ArgsType.(*AddArgs)
		buff, err := json.Marshal(args)
		fatalError(err)
		call(&buff, ServiceName)
		var res AddArgs
		err = json.Unmarshal(buff, &res)
		fatalError(err)
		args = res
	case "sub":
		args := ArgsType.(*SubArgs)
		buff, err := json.Marshal(args)
		fatalError(err)
		call(&buff, ServiceName)
		var res SubArgs
		err = json.Unmarshal(buff, &res)
		fatalError(err)
		args = res
	default:
		fatalError(fmt.Errorf("service %s is not found", ServiceName))
	}

}

// CallByte
// This function is used to call a service by name
//
// @ServiceName is the name of the service
//
// @ArgsByte is the arguments of the service, you can get the result from it
//
// Note: don't use this function to call a service asynchronously
func CallByte(ServiceName string, ArgsByte *[]byte) bool {
	seq := getSeq()
	msg := append(runH, seq...)
	msg = append(msg, *ArgsByte...)
	TO, addr := req(ServiceName)
	if addr == nil {
		printError(fmt.Errorf("service %s is busy or offline", ServiceName))
		return false
	} else {
		conn := serviceSocket[ServiceName]
		t1 := 60
		for tt := 0; tt < timeOUT; {
			recv := srSendAndRecv(msg, seq, t1, &conn, addr)
			if len(recv) < 10 {
				err := conn.SetReadDeadline(time.Now().Add(time.Duration(TO) * time.Millisecond))
				fatalError(err)
				n, _, err := conn.ReadFrom(recv)
				if printError(err) || n < 26 {
					tt += t1 + int(TO)
					t1 = t1 * 2
					TO = TO * 3 / 2
					continue
				}
				if string(recv[0:12]) == "SERVER\r\nACK:" && cmp4(recv[12:16], seq) && string(recv[16:26]) == "\r\nSTATUS:0" {
					fatalError(fmt.Errorf("service %s failed to run, please check your args", ServiceName))
					return true
				} else if string(recv[0:12]) == "SERVER\r\nACK:" && cmp4(recv[12:16], seq) && string(recv[16:26]) == "\r\nSTATUS:1" {
					*ArgsByte = recv[26:n]
					return true
				}
			} else if string(recv[:10]) == "\r\nSTATUS:0" {
				fatalError(fmt.Errorf("service %s failed to run, please check your args", ServiceName))
				return true
			} else if string(recv[:10]) == "\r\nSTATUS:1" {
				*ArgsByte = recv[10:]
				return true
			}
		}
		printError(fmt.Errorf("service %s is busy or offline", ServiceName))
		return false
	}
}

func req(name string) (int64, *net.UDPAddr) {
	service := *(*urls)(*myServices[name])
	m := len(service.ips)
	if m == 0 {
		return 0, nil
	}
	var factors = make([]int32, m)
	copy(factors, service.factors)
	s := service.s
	for i := 0; i < m; i++ {
		r := rd.Int31n(s)
		for j := 0; j < m; j++ {
			if r < factors[j] {
				s -= factors[i]
				factors[i] = 0
				addr := &net.UDPAddr{
					IP:   net.ParseIP(service.ips[j]),
					Port: service.ports[j],
				}
				conn := serviceSocket[name]
				if ok, timeout := dailAndWait(addr, &conn, 120); ok {
					return timeout, addr
				}
				break
			}
			r -= factors[j]
		}
	}
	return 0, nil
}

func dailAndWait(addr *net.UDPAddr, conn *net.PacketConn, timeout int) (bool, int64) {
	seq := getSeq()
	body := append(reqH, seq...)
	res := srSendAndRecv(body, seq, timeout, conn, addr)
	if len(res) < 10 {
		return false, 0
	} else if string(res[:10]) == "\r\nTIMEOUT:" {
		return true, decode64(res[10:])
	} else {
		return false, 0
	}
}

func Close() {
	err := rSocket.Close()
	if err != nil {
		return
	}
	for _, v := range serviceSocket {
		err := v.Close()
		if err != nil {
			return
		}
	}
}

func call(buff *[]byte, ServiceName string) {
	for t := 1; t <= 3; t++ {
		ok := CallByte(ServiceName, buff)
		if ok {
			break
		} else {
			fmt.Printf("It will retry later [%d/3]\n", t)
		}
		time.Sleep(time.Duration(rd.Intn(1000)) * time.Millisecond)
	}
}
