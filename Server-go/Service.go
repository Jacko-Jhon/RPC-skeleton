package main

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
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
}

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
		buf := make([]byte, 1024)
		n, addr, err := s.ServerSocket.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if string(buf[:17]) == "CLIENT\r\nREQ\r\nSEQ:" && s.MyProcess < s.MaxProcess {
			go s.ReturnTimout(buf[17:n], addr)
		} else if string(buf[:17]) == "CLIENT\r\nRUN\r\nSEQ:" {
			go s.HandleRequest(buf[17:n], addr)
		}
	}
}

func (s *Service) ReturnTimout(msg []byte, addr *net.UDPAddr) {
	T1 := time.Now().Unix() * 1000
	T2 := decode64(msg[4:12])
	seqByte := msg[0:4]
	seq := decode32(seqByte)
	h := []byte("SERVER\r\nACK:")
	ack := encode32(seq + 1)
	timeoutH := []byte("\r\nTIMEOUT:")
	timeout := encode64(2*(T1-T2) + 60 + int64(s.Timeout))
	res := append(h, ack...)
	res = append(res, timeoutH...)
	res = append(res, timeout...)
	_, err := s.ServerSocket.WriteToUDP(res, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *Service) HandleRequest(msg []byte, addr *net.UDPAddr) {
	atomic.AddInt32(&s.NumOfRequest, 1)
	atomic.AddInt32(&s.MyProcess, 1)
	seqByte := msg[0:4]
	seq := decode32(seqByte)
	h := []byte("SERVER\r\nACK:")
	ack := encode32(seq + 1)
	var status []byte
	ret, err := s.function.run(msg[4:])
	if err != nil {
		status = []byte("\r\nSTATUS:0")
	} else {
		status = []byte("\r\nSTATUS:1")
	}
	res := append(h, ack...)
	res = append(res, status...)
	res = append(res, ret...)
	for i := 0; i < SendTimes; i++ {
		_, err = s.ServerSocket.WriteToUDP(res, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	atomic.AddInt32(&s.MyProcess, -1)
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
		time.Sleep(time.Duration(ExPause) * time.Second)
	}
}
