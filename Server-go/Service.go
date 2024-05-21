package main

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	Id             string
	Name           string
	Port           int
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
}

func (s *Service) Run() {
	go s.Heartbeat()
	for {
		buf := make([]byte, 1024)
		n, addr, err := s.ServerSocket.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if string(buf[:12]) == "CLIENT\r\nSEQ:" {
			s.HandleRequest(buf[12:n], addr)
		}
	}
}

func (s *Service) HandleRequest(msg []byte, addr *net.UDPAddr) {
	seqByte := msg[0:4]
	seq := decode(seqByte)
	h := []byte("SERVER\r\nACK:")
	ack := encode(seq + 1)
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
	for i := 0; i < 3; i++ {
		_, err = s.RegisterSocket.WriteToUDP(res, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (s *Service) Heartbeat() {
	for {
		seq := int(rand.Int31())
		if seq > 100 {
			seq -= 50
		}
		h := []byte("SERVER\r\nSEQ:")
		seqB := encode(seq)
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
	}
}
