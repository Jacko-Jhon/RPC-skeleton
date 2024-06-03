package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FunctionInfo struct {
	Port       int      `json:"port"`
	Name       string   `json:"name"`
	Id         string   `json:"id"`
	Args       []string `json:"args"`
	Ret        []string `json:"ret"`
	Timeout    int      `json:"timeout"`
	MaxProcess int32    `json:"MaxProcess"`
}

type Server struct {
	Ip              string
	Port            int
	RegisterAddr    string
	RegisterSocket  *net.UDPConn
	ServiceList     map[string]*Service
	LoadedFunctions map[string]FunctionInfo
}

var GlobalServer = Server{
	Ip:              "127.0.0.1",
	Port:            8080,
	RegisterAddr:    "127.0.0.1:8888",
	ServiceList:     make(map[string]*Service),
	LoadedFunctions: make(map[string]FunctionInfo),
}

func (s *Server) Start() {
	w := sync.WaitGroup{}
	w.Add(len(s.ServiceList))
	for _, sv := range s.ServiceList {
		sv.Init()
		w.Done()
		go sv.Run()
	}
	w.Wait()
}

func (s *Server) DailRegistry(msg ServiceMessage, opCode string) MessageToServer {
	seq := rand.Int31()
	h := []byte("SERVER\r\nSEQ:")
	seqB := encode32(seq)
	op := []byte("\r\nOP:" + opCode)
	js := msg.ToJson()
	body := append(h, seqB...)
	body = append(body, op...)
	body = append(body, js...)
	res := make([]byte, 1024)
	var err error
	for i := 0; i < 3; i++ {
		_, err = s.RegisterSocket.Write(body)
		if err != nil {
			panic(err)
		}
		err = s.RegisterSocket.SetReadDeadline(time.Now().Add(time.Millisecond * 460))
		if err != nil {
			fmt.Println(err)
			continue
		}
		n, _, err := s.RegisterSocket.ReadFrom(res)
		if err == nil {
			ackByte := res[12:16]
			ack := decode32(ackByte)
			if ack == seq && string(res[0:12]) == "REGIST\r\nACK:" {
				var resMsg MessageToServer
				err1 := json.Unmarshal(res[16:n], &resMsg)
				if err1 != nil {
					panic(err1)
				}
				return resMsg
			}
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	return MessageToServer{
		Info:   "timeout",
		Status: false,
	}
}

func (s *Server) Init() {
	splitAddr := strings.Split(s.RegisterAddr, ":")
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
	s.CheckLoadFunctions()
}

func (s *Server) CheckLoadFunctions() {
	for i, fn := range MyFunctions {
		port := s.Port + i
		_, ok := s.LoadedFunctions[fn.Name]
		if ok {
			s.ServiceList[fn.Name] = &Service{
				Id:   s.LoadedFunctions[fn.Name].Id,
				Port: port,
				Name: fn.Name,
				function: Function{
					Name:       fn.Name,
					Args:       fn.Args,
					Ret:        fn.Ret,
					Id:         s.LoadedFunctions[fn.Name].Id,
					Timeout:    s.LoadedFunctions[fn.Name].Timeout,
					MaxProcess: s.LoadedFunctions[fn.Name].MaxProcess,
					run:        fn.run,
				},
			}
		}
	}
}

func (s *Server) Register(f Function, port int) {
	opCode := "1"
	if f.Id != "" {
		opCode = "2"
	}
	msg := ServiceMessage{
		Id:   f.Id,
		Name: f.Name,
		Ip:   s.Ip,
		Port: port,
		Args: f.Args,
		Ret:  f.Ret,
	}
	res := s.DailRegistry(msg, opCode)
	if res.Status {
		fmt.Println(f.Name, "register success")
		f.Id = res.Id
		s.ServiceList[f.Name] = &Service{
			function:       f,
			Id:             res.Id,
			Port:           port,
			Name:           f.Name,
			ServerSocket:   nil,
			RegisterSocket: nil,
		}
	} else {
		fmt.Println(f.Name, "register failed: ", res.Info)
	}
}

func (s *Server) UnRegister(id string, name string) {
	opCode := "3"
	msg := ServiceMessage{
		Id: id,
	}
	res := s.DailRegistry(msg, opCode)
	if res.Status {
		delete(s.ServiceList, name)
		fmt.Println("unregister success")
	} else {
		fmt.Println("unregister failed: ", res.Info)
	}
}

func (s *Server) UpdateUrl(id string, ip string, port int) bool {
	opCode := "5"
	msg := ServiceMessage{
		Id:   id,
		Ip:   ip,
		Port: port,
	}
	res := s.DailRegistry(msg, opCode)
	if res.Status {
		return true
	} else {
		fmt.Println("update url failed: ", res.Info)
		return false
	}
}

func (s *Server) Dump(path string) {
	fp := filepath.Clean(path)
	f, err1 := os.Create(fp)
	if err1 != nil {
		log.Fatal(err1)
	}
	functions := make([]FunctionInfo, 0, len(s.ServiceList))
	for _, service := range s.ServiceList {
		if service.function.run == nil { // 函数为空则不会保存到本地
			continue
		}
		functions = append(functions, FunctionInfo{
			Name:       service.Name,
			Id:         service.Id,
			Port:       service.Port,
			Args:       service.function.Args,
			Ret:        service.function.Ret,
			MaxProcess: service.function.MaxProcess,
			Timeout:    service.function.Timeout,
		})
	}
	jsonData, err := json.Marshal(functions)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("dump to", path)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
}

func (s *Server) Load(path string) {
	fp := filepath.Clean(path)
	f, err := os.Open(fp)
	if err != nil {
		return
	}
	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
	var functions []FunctionInfo
	err = json.Unmarshal(data, &functions)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range functions {
		if f.Id == "" {
			continue
		}
		s.LoadedFunctions[f.Name] = FunctionInfo{
			Name:       f.Name,
			Id:         f.Id,
			Port:       f.Port,
			Args:       f.Args,
			Ret:        f.Ret,
			MaxProcess: f.MaxProcess,
			Timeout:    f.Timeout,
		}
	}
	fmt.Println("load from", path)
}
