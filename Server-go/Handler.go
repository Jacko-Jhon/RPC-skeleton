package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	SleepTime          = 25
	factor       int32 = 1000
	ConfigPath         = "./config.json"
	QuitSignal         = ""
	SendTimes          = 1
	isStart      bool
	isUnregister bool
	isList       bool
)

func encode64(n int64) []byte {
	code := make([]byte, 8)
	code[0] = byte(n >> 56)
	code[1] = byte(n >> 48)
	code[2] = byte(n >> 40)
	code[3] = byte(n >> 32)
	code[4] = byte(n >> 24)
	code[5] = byte(n >> 16)
	code[6] = byte(n >> 8)
	code[7] = byte(n)
	return code
}

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

func decode32(code []byte) int32 {
	return int32(code[0])<<24 | int32(code[1])<<16 | int32(code[2])<<8 | int32(code[3])
}

func cmp4(a []byte, b []byte) bool {
	if a[0] == b[0] && a[1] == b[1] && a[2] == b[2] && a[3] == b[3] {
		return true
	}
	return false
}

func CheckMyFunction() {
	fns := make(map[string]int, 0)
	for i, fn := range MyFunctions {
		_, ok := GlobalServer.ServiceList[fn.Name]
		if !ok {
			fns[fn.Name] = i
		}
	}
	if len(fns) > 0 {
		fmt.Println("You have functions that are not registered, do you want to register them?")
		fmt.Println("\tyou can input the name of them to register, like, 'add sub mul div'")
		fmt.Println("\tif you want to register all of them, just input '-all'")
		fmt.Println("\telse if you don't want to register, just input '-no'")
		var i = 1
		for name := range fns {
			fmt.Println(strconv.Itoa(i)+",", name)
			i++
		}
		var input string
		rd := bufio.NewReader(os.Stdin)
		bytes, _, _ := rd.ReadLine()
		input = string(bytes)
		if input == "-all" {
			for _, i := range fns {
				port := i + GlobalServer.Port
				GlobalServer.Register(MyFunctions[i], port)
			}
		} else if input == "-no" {
			return
		} else {
			fn := strings.Split(input, " ")
			for _, name := range fn {
				i, ok := fns[name]
				if ok {
					port := i + GlobalServer.Port
					GlobalServer.Register(MyFunctions[i], port)
				}
			}
		}
		GlobalServer.Dump(ConfigPath)
	}
}

func Unregister() {
	fmt.Println("Please input the id of the service you want to unregister")
	fmt.Println("\tyou can input the id of them to unregister, like, 'name:ID1 name:ID2 name:ID3'")
	fmt.Println("\tif you want to unregister all of them, just input '-all'")
	fmt.Println("\telse if you don't want to unregister, just input '-no'")
	var input string
	rd := bufio.NewReader(os.Stdin)
	bytes, _, _ := rd.ReadLine()
	input = string(bytes)
	switch input {
	case "-all":
		fmt.Println("Do you want to unregister all services ?")
		fmt.Print("Y/N  >>> ")
		rd := bufio.NewReader(os.Stdin)
		bytes, _, _ := rd.ReadLine()
		isSure := string(bytes)
		if isSure == "Y" || isSure == "y" {
			for _, service := range GlobalServer.ServiceList {
				GlobalServer.UnRegister(service.Id, service.Name)
				delete(GlobalServer.ServiceList, service.Name)
			}
		}
		GlobalServer.Dump(ConfigPath)
	case "-no":
		return
	default:
		idAndNames := strings.Split(input, " ")
		for _, idn := range idAndNames {
			sp := strings.Split(idn, ":")
			if len(sp) == 2 {
				if GlobalServer.ServiceList[sp[0]].Id == sp[1] {
					GlobalServer.UnRegister(sp[1], sp[0])
					delete(GlobalServer.ServiceList, sp[0])
				} else {
					fmt.Println("Cannot find service", idn)
				}
			} else {
				fmt.Println("Invalid input", idn)
			}
		}
		GlobalServer.Dump(ConfigPath)
	}
}

func ListServices() {
	var i = 1
	fmt.Println("List of services:")
	for name, service := range GlobalServer.ServiceList {
		fmt.Println(strconv.Itoa(i)+",", name+":"+service.Id)
		i++
	}
}

func main() {
	flag.StringVar(&GlobalServer.RegisterAddr, "R", "127.0.0.1:8888", "Set register address")
	flag.StringVar(&GlobalServer.Ip, "l", "0.0.0.0", "Set listen address")
	flag.IntVar(&GlobalServer.Port, "p", 0,
		"Set listen port. \nNote: every service has its own port,\nfor example, the port of first service is 8080, \nand the port of second service is 8081")
	flag.BoolVar(&isStart, "start", false, "Start the server")
	flag.BoolVar(&isUnregister, "unregister", false, "Unregister the services")
	flag.BoolVar(&isList, "list", false, "List all services")
	mode := flag.String("mode", "LAN", "Set mode, such as: LAN, WAN, Boost")
	flag.Parse()
	if isStart && GlobalServer.Port == 0 {
		fmt.Println("You must set port to start")
		return
	}
	GlobalServer.Port = 10086
	GlobalServer.Load(ConfigPath)
	GlobalServer.Init()
	defer func(s Server) {
		if QuitSignal == "Quit" {
			GlobalServer.Dump(ConfigPath)
		}
		err := s.RegisterSocket.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(GlobalServer)
	fmt.Println("Server is init")
	if isList {
		ListServices()
	}
	if isUnregister && isStart {
		fmt.Println("You can't unregister and start at the same time")
		return
	} else if isUnregister && !isStart {
		Unregister()
		return
	} else if !isUnregister && isStart {
		fmt.Println("Start the server...")
		fmt.Println("Registry Address :", GlobalServer.RegisterAddr)
		fmt.Println("Listen Address   :", GlobalServer.Ip, ":", GlobalServer.Port)
		switch *mode {
		case "LAN":
			SendTimes = 1
		case "WAN":
			SendTimes = 2
		case "Boost":
			SendTimes = 3
		}
	} else {
		return
	}
	CheckMyFunction()
	GlobalServer.Start()
	if len(GlobalServer.ServiceList) > 0 {
		fmt.Println("Server started")
	} else {
		return
	}
	defer fmt.Println("Server is closed")
	for {
		_, err := fmt.Scanf("%s\n", &QuitSignal)
		if err != nil {
			fmt.Println(err)
		}
		if QuitSignal == "Quit" {
			return
		}
	}
}
