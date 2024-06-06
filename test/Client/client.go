
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
	serviceName   = []string{"sub", "unregister", "DelayTest", "GetUrl", "add", "div", "mul", "MergeSort", "QSort", "register"}
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
			recv, flag := srSendAndRecv(msg, seq, t1, &conn, addr)
			if flag == WriteError || flag == ReadError {
				return false
			} else if flag < 10 {
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
	ptr := atomic.LoadPointer(myServices[name])
	service := *(*urls)ptr
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
				ok, timeout := dailAndWait(addr, &conn, 60); 
				if ok {
					return timeout, addr
				} else if timeout == -1 {
					service.factors[i] = 0
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
	res, n := srSendAndRecv(body, seq, timeout, conn, addr)
	if n == WriteError || n == ReadError {
		return false, -1
	} else if n == 18 && string(res[:10]) == "\r\nTIMEOUT:" {
		return true, decode64(res[10:n])
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
		case "sub":
		args := ArgsType.(*SubArgs)
		buff, err := json.Marshal(args)
		fatalError(err)
		call(&buff, ServiceName)
		var res SubArgs
		err = json.Unmarshal(buff, &res)
		fatalError(err)
		ArgsType = res
	case "unregister":
		args := ArgsType.(*UnregisterArgs)
		buff, err := json.Marshal(args)
		fatalError(err)
		call(&buff, ServiceName)
		var res UnregisterArgs
		err = json.Unmarshal(buff, &res)
		fatalError(err)
		ArgsType = res
	case "DelayTest":
		args := ArgsType.(*DelayTestArgs)
		buff, err := json.Marshal(args)
		fatalError(err)
		call(&buff, ServiceName)
		var res DelayTestArgs
		err = json.Unmarshal(buff, &res)
		fatalError(err)
		ArgsType = res
	case "GetUrl":
		args := ArgsType.(*GetUrlArgs)
		buff, err := json.Marshal(args)
		fatalError(err)
		call(&buff, ServiceName)
		var res GetUrlArgs
		err = json.Unmarshal(buff, &res)
		fatalError(err)
		ArgsType = res
	case "add":
		args := ArgsType.(*AddArgs)
		buff, err := json.Marshal(args)
		fatalError(err)
		call(&buff, ServiceName)
		var res AddArgs
		err = json.Unmarshal(buff, &res)
		fatalError(err)
		ArgsType = res
	case "div":
		args := ArgsType.(*DivArgs)
		buff, err := json.Marshal(args)
		fatalError(err)
		call(&buff, ServiceName)
		var res DivArgs
		err = json.Unmarshal(buff, &res)
		fatalError(err)
		ArgsType = res
	case "mul":
		args := ArgsType.(*MulArgs)
		buff, err := json.Marshal(args)
		fatalError(err)
		call(&buff, ServiceName)
		var res MulArgs
		err = json.Unmarshal(buff, &res)
		fatalError(err)
		ArgsType = res
	case "MergeSort":
		args := ArgsType.(*MergeSortArgs)
		buff, err := json.Marshal(args)
		fatalError(err)
		call(&buff, ServiceName)
		var res MergeSortArgs
		err = json.Unmarshal(buff, &res)
		fatalError(err)
		ArgsType = res
	case "QSort":
		args := ArgsType.(*QSortArgs)
		buff, err := json.Marshal(args)
		fatalError(err)
		call(&buff, ServiceName)
		var res QSortArgs
		err = json.Unmarshal(buff, &res)
		fatalError(err)
		ArgsType = res
	case "register":
		args := ArgsType.(*RegisterArgs)
		buff, err := json.Marshal(args)
		fatalError(err)
		call(&buff, ServiceName)
		var res RegisterArgs
		err = json.Unmarshal(buff, &res)
		fatalError(err)
		ArgsType = res
	default:
		fatalError(fmt.Errorf("service %s is not found", ServiceName))
	}
}