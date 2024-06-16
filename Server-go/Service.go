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
	cache          sync.Map // 缓存结果
}

type buffer struct { // 缓存结构体
	seq    []byte  // 请求序列号
	data   *[]byte // 结果
	isDone bool    // 是否完成
	time   int64   // 缓存时间
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

// Run starts the service
func (s *Service) Run() {
	// 启动心跳
	go s.Heartbeat()
	defer func(ServerSocket *net.UDPConn) {
		err := ServerSocket.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(s.ServerSocket)
	// 启动服务
	for {
		buf := make([]byte, 4096)
		// 接收请求
		n, addr, err := s.ServerSocket.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// 处理请求
		go s.Handle(n, buf, addr)
	}
}

// Handle handles the request and returns the result
func (s *Service) Handle(n int, buf []byte, addr *net.UDPAddr) {
	// 解析请求类型
	if string(buf[:17]) == "CLIENT\r\nREQ\r\nSEQ:" {
		if atomic.LoadInt32(&s.MyProcess) < s.MaxProcess {
			// 如果当前并发数没有达到最大值，则返回超时时间
			s.ReturnTimout(buf[17:n], addr)
		} else {
			// 如果当前并发数已经达到最大值，则返回空ACK表明服务繁忙
			data := append(h, buf[17:21]...)
			_, err := s.ServerSocket.WriteToUDP(data, addr)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else if string(buf[:17]) == "CLIENT\r\nRUN\r\nSEQ:" {
		data := append(h, buf[17:21]...)
		// 返回ACK表明服务端已接受到请求
		_, err := s.ServerSocket.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println(err)
		}
		// 查看缓存是否有结果
		ptr, ok := s.cache.Load(addr.String())
		if ok {
			p := atomic.LoadPointer(ptr.(*unsafe.Pointer))
			b := (*buffer)(p)
			if b.isDone && b.time+120 < time.Now().Unix() && cmp4(buf[17:21], b.seq) {
				// 缓存结果匹配，直接返回结果
				_, err = s.ServerSocket.WriteToUDP(*b.data, addr)
				if err != nil {
					fmt.Println(err)
				}
			} else if !b.isDone {
				// 该请求正在处理中
				return
			} else {
				// 缓存结果过期或不匹配，重新处理请求
				atomic.StorePointer(ptr.(*unsafe.Pointer), unsafe.Pointer(&buffer{
					isDone: false,
				}))
				s.HandleRequest(buf[17:n], addr, ptr.(*unsafe.Pointer))
			}
		} else {
			// 缓存中没有结果， 处理请求
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
	// 解析客户端发送请求时的时间戳
	T2 := decode64(msg[4:12])
	seq := msg[0:4]
	// 计算超时时间
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

// HandleRequest handles the request and returns the result
func (s *Service) HandleRequest(msg []byte, addr *net.UDPAddr, cache *unsafe.Pointer) {
	i := 0
	for t := int32(20); atomic.LoadInt32(&s.MyProcess) >= s.MaxProcess && i < 5; t *= 2 {
		// 二进制指数避让, 超过5次就直接进入处理
		time.Sleep(time.Millisecond * time.Duration(rand.Int31n(t)))
		i++
	}
	// 原子操作, 标记并发数以实现最大并发数限制
	atomic.AddInt32(&s.NumOfRequest, 1)
	atomic.AddInt32(&s.MyProcess, 1)
	seq := msg[0:4]
	// 调用服务函数并返回结果
	ret, err := s.function.run(msg[4:])
	res := append(h, seq...)
	if err != nil {
		res = append(res, status0...)
	} else {
		res = append(res, status1...)
	}
	// 发送结果， 消息格式：h + seq + status + result = “SERVER\r\nACK:xxxx\r\nSTATUS:x” + result
	res = append(res, ret...)
	for i := 0; i < SendTimes; i++ {
		// 重发机制， 局域网内发1次， 外网发2次，Boost模式下发3次 （由服务启动时模式决定，默认发1次）
		_, err = s.ServerSocket.WriteToUDP(res, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	// 原子操作, 标记并发数以实现最大并发数限制
	atomic.AddInt32(&s.MyProcess, -1)
	// 缓存结果， 用于处理丢包情况
	b := buffer{
		seq:    seq,
		isDone: true,
		data:   &res,
		time:   time.Now().Unix(),
	}
	ptr := unsafe.Pointer(&b)
	atomic.StorePointer(cache, ptr)
}

// Heartbeat sends heartbeat messages to the register server
func (s *Service) Heartbeat() {
	defer func(RegisterSocket *net.UDPConn) {
		err := RegisterSocket.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(s.RegisterSocket)
	for {
		// 编写心跳包
		h := []byte("SERVER\r\nSEQ:")
		seqB := encode32(rand.Int31())
		op := []byte("\r\nOP:0")
		msg := ServiceMessage{
			Id: s.Id,
		}
		js := msg.ToJson()
		body := append(h, seqB...)
		body = append(body, op...)
		body = append(body, js...)
		// 发送心跳包
		_, err := s.RegisterSocket.Write(body)
		if err != nil {
			fmt.Println(err)
			return
		}
		// 暂停一段时间
		time.Sleep(time.Duration(SleepTime) * time.Second)
		// 计算额外延迟心跳时间 计算公式：延迟心跳时间 = 一段时间内并发请求数 * 设置的函数超时时间 / 指定因子
		ExPause := s.NumOfRequest * int32(s.Timeout) / factor
		// 打印并发请求数，用于查看日志以及方便测试
		fmt.Println("Service", s.Name, "has ran", s.NumOfRequest, "requests in the past 25s")
		// 暂停额外延迟心跳时间
		time.Sleep(time.Duration(ExPause) * time.Second)
	}
}
