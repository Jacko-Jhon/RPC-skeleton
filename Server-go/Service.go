package main

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type Service struct {
	Id             string
	Name           string
	Port           int
	NumOfRequest   int32
	MyProcess      int32
	MaxProcess     int32
	Timeout        int
	ServerSocket   *net.UDPConn
	RegisterSocket *net.UDPConn
	function       Function
	cache          sync.Map
}

type buffer struct {
	seq    []byte
	data   *[]byte
	isDone bool
	time   int64
}

var (
	h        = []byte("SERVER\r\nACK:")
	timeoutH = []byte("\r\nTIMEOUT:")
	status0  = []byte("\r\nSTATUS:0")
	status1  = []byte("\r\nSTATUS:1")
)

func (s *Service) Init() {
	splitAddr := strings.Split(GlobalServer.RegisterAddr, ":")
	port, err := strconv.Atoi(splitAddr[1])
	if err != nil {
		panic(err)
	}
	addr := &net.UDPAddr{IP: net.ParseIP(splitAddr[0]), Port: port}
	dail, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}
	s.RegisterSocket = dail
	url := GlobalServer.Ip + ":" + strconv.Itoa(s.Port)
	Addr, err := net.ResolveUDPAddr("udp", url)
	if err != nil {
		panic(err)
	}
	sk, err := net.ListenUDP("udp", Addr)
	if err != nil {
		panic(err)
	}
	s.ServerSocket = sk
	s.NumOfRequest = 0
	s.MyProcess = 0
	s.MaxProcess = s.function.MaxProcess
	s.Timeout = s.function.Timeout
	flag := GlobalServer.UpdateUrl(s.Id, GlobalServer.Ip, s.Port)
	if flag {
		fmt.Println("Service "+s.Name+" is listening port", s.Port)
	} else {
		panic("Service " + s.Name + " Start Failed")
	}
}

func (s *Service) Run() {
	go s.Heartbeat()
	defer func(ServerSocket *net.UDPConn) {
		err := ServerSocket.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(s.ServerSocket)
	for {
		buf := make([]byte, 4096)
		n, addr, err := s.ServerSocket.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		go s.Handle(n, buf, addr)
	}
}

func (s *Service) Handle(n int, buf []byte, addr *net.UDPAddr) {
	if string(buf[:17]) == "CLIENT\r\nREQ\r\nSEQ:" {
		if atomic.LoadInt32(&s.MyProcess) < s.MaxProcess {
			s.ReturnTimout(buf[17:n], addr)
		} else {
			data := append(h, buf[17:21]...)
			_, err := s.ServerSocket.WriteToUDP(data, addr)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else if string(buf[:17]) == "CLIENT\r\nRUN\r\nSEQ:" {
		data := append(h, buf[17:21]...)
		_, err := s.ServerSocket.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println(err)
		}
		ptr, ok := s.cache.Load(addr.String())
		if ok {
			p := atomic.LoadPointer(ptr.(*unsafe.Pointer))
			b := (*buffer)(p)
			if b.isDone && b.time+120 < time.Now().Unix() && cmp4(buf[17:21], b.seq) {
				_, err = s.ServerSocket.WriteToUDP(*b.data, addr)
				if err != nil {
					fmt.Println(err)
				}
			} else if !b.isDone {
				return
			} else {
				atomic.StorePointer(ptr.(*unsafe.Pointer), unsafe.Pointer(&buffer{
					isDone: false,
				}))
				s.HandleRequest(buf[17:n], addr, ptr.(*unsafe.Pointer))
			}
		} else {
			tmpPtr := unsafe.Pointer(&buffer{
				isDone: false,
			})
			s.cache.Store(addr.String(), &tmpPtr)
			s.HandleRequest(buf[17:n], addr, &tmpPtr)
		}
	}
}

func (s *Service) ReturnTimout(msg []byte, addr *net.UDPAddr) {
	T1 := time.Now().UnixMilli()
	T2 := decode64(msg[4:12])
	seq := msg[0:4]
	timeout := encode64(2*(T1-T2) + 60 + int64(s.Timeout))
	res := append(h, seq...)
	res = append(res, timeoutH...)
	res = append(res, timeout...)
	_, err := s.ServerSocket.WriteToUDP(res, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *Service) HandleRequest(msg []byte, addr *net.UDPAddr, cache *unsafe.Pointer) {
	for t := int32(50); atomic.LoadInt32(&s.MyProcess) >= s.MaxProcess; t *= 2 {
		// 二进制指数避让
		time.Sleep(time.Millisecond * time.Duration(rand.Int31n(t)))
	}
	atomic.AddInt32(&s.NumOfRequest, 1)
	atomic.AddInt32(&s.MyProcess, 1)
	seq := msg[0:4]
	ret, err := s.function.run(msg[4:])
	res := append(h, seq...)
	if err != nil {
		res = append(res, status0...)
	} else {
		res = append(res, status1...)
	}
	res = append(res, ret...)
	for i := 0; i < SendTimes; i++ {
		_, err = s.ServerSocket.WriteToUDP(res, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	atomic.AddInt32(&s.MyProcess, -1)
	b := buffer{
		seq:    seq,
		isDone: true,
		data:   &res,
		time:   time.Now().Unix(),
	}
	ptr := unsafe.Pointer(&b)
	atomic.StorePointer(cache, ptr)
}

func (s *Service) Heartbeat() {
	defer func(RegisterSocket *net.UDPConn) {
		err := RegisterSocket.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(s.RegisterSocket)
	for {
		seq := rand.Int31()
		if seq > 100 {
			seq -= 50
		}
		h := []byte("SERVER\r\nSEQ:")
		seqB := encode32(seq)
		op := []byte("\r\nOP:0")
		msg := ServiceMessage{
			Id: s.Id,
		}
		js := msg.ToJson()
		body := append(h, seqB...)
		body = append(body, op...)
		body = append(body, js...)
		_, err := s.RegisterSocket.Write(body)
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(time.Duration(SleepTime) * time.Second)
		ExPause := s.NumOfRequest / factor
		fmt.Println("Service", s.Name, "has ran", s.NumOfRequest, "requests in the past 25s")
		time.Sleep(time.Duration(ExPause) * time.Second)
	}
}
