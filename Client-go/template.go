package main

var util_tpl = `
package Client

import (
	"net"
	"time"
)

func decode64(code []byte) int64 {
	return int64(code[0])<<56 | int64(code[1])<<48 | int64(code[2])<<40 | int64(code[3])<<32 | int64(code[4])<<24 | int64(code[5])<<16 | int64(code[6])<<8 | int64(code[7])
}

func encode32(n int32) []byte {
	code := make([]byte, 4)
	code[0] = byte(n >> 24)
	code[1] = byte(n >> 16)
	code[2] = byte(n >> 8)
	code[3] = byte(n)
	return code
}

func getSeq() []byte {
	r := rd.Int31n(16)
	t := time.Now().UnixMicro()
	seq := int32(t)<<4 | r
	return encode32(seq)
}

func dailRegistry(name string) []byte {
	h := []byte("CLIENT\r\nSEQ:")
	seq := getSeq()
	op := []byte("\r\nOP:0")
	body := append(h, seq...)
	body = append(body, op...)
	body = append(body, []byte("\r\nNAME:"+name)...)
	return rgSendAndRecv(body, seq, 120)
}

func cmp4(a []byte, b []byte) bool {
	if a[0] == b[0] && a[1] == b[1] && a[2] == b[2] && a[3] == b[3] {
		return true
	}
	return false
}

func srSendAndRecv(msg []byte, seq []byte, timeout int, conn *net.PacketConn, addr *net.UDPAddr) []byte {
	res := make([]byte, buffSize)
	tt := 0
	isResend := true
	for {
		if isResend {
			isResend = false
			_, err := (*conn).WriteTo(msg, addr)
			if printError(err) {
				return nil
			}
		}
		err := (*conn).SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(timeout)))
		n, _, err := (*conn).ReadFrom(res)
		if err != nil {
			tt += timeout
			if tt >= timeOUT {
				printError(err)
				return nil
			}
			isResend = true
			timeout = timeout * 3 / 2
		} else if string(res[0:12]) == "SERVER\r\nACK:" && cmp4(res[12:16], seq) {
			return res[16:n]
		}
	}
}

func rgSendAndRecv(msg []byte, seq []byte, timeout int) []byte {
	var res = make([]byte, buffSize)
	tt := 0
	for {
		_, err := rSocket.Write(msg)
		if printError(err) {
			return nil
		}
		err = rSocket.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(timeout)))
		n, _, err := rSocket.ReadFrom(res)
		if err != nil {
			tt += timeout
			if tt >= timeOUT {
				printError(err)
				return nil
			}
			timeout = timeout * 3 / 2
		} else if string(res[0:12]) == "REGIST\r\nACK:" && cmp4(res[12:16], seq) {
			return res[16:n]
		}
	}
}
`

var error_tpl = `
package Client

import (
	"fmt"
	"log"
)

func fatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func printError(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}
`

var struct_tpl = "package Client\ntype serviceUrls struct {\n" +
	"\tStatus  bool     `json:\"status\"`\n" +
	"\tInfo    string   `json:\"info\"`\n" +
	"\tName    string   `json:\"name\"`\n" +
	"\tIps     []string `json:\"ips\"`\n" +
	"\tPorts   []int    `json:\"ports\"`\n" +
	"\tFactors []int32  `json:\"factors\"`\n" +
	"\tArgs    []string `json:\"args\"`\n" +
	"\tRet     []string `json:\"ret\"`\n}\n" +
	`
type urls struct {
	s       int32
	ips     []string
	ports   []int
	factors []int32
}


`

var client_tpl = `
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
	serviceName   = []string{$case$}
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


`

var call_tpl = `
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
	$case$
	default:
		fatalError(fmt.Errorf("service %s is not found", ServiceName))
	}
}`
