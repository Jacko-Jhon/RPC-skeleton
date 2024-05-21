package main

import (
	"flag"
	"fmt"
)

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

var SleepTime = 25

func main() {
	flag.StringVar(&GlobalServer.RegisterAddr, "R", "127.0.0.1:8888", "Set register address")
	flag.StringVar(&GlobalServer.Ip, "l", "0.0.0.0", "Set listen address")
	flag.IntVar(&GlobalServer.Port, "p", 8080,
		"Set listen port. \nNote: every service has its own port,\nfor example, the port of first service is 8080, \nand the port of second service is 8081")
	flag.Parse()
	GlobalServer.Init()
	GlobalServer.Start()
	defer func(s Server) {
		err := s.RegisterSocket.Close()
		if err != nil {
			fmt.Println(err)
		}
		for _, sv := range s.ServiceList {
			err := sv.RegisterSocket.Close()
			if err != nil {
				fmt.Println(err)
			}
			err = sv.ServerSocket.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
	}(GlobalServer)
}
