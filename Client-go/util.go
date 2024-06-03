package main

import (
	"fmt"
	"math/rand"
	"time"
)

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

func DailRegistry(opCode string, Name ...string) []byte {
	seq := rand.Int31()
	h := []byte("CLIENT\r\nSEQ:")
	seqB := encode32(seq)
	op := []byte("\r\nOP:" + opCode)
	body := append(h, seqB...)
	body = append(body, op...)
	res := make([]byte, 1024)
	if opCode == "0" && Name[0] != "" {
		body = append(body, []byte("\r\nNAME:"+Name[0])...)
	} else if opCode == "1" && Name[0] != "" {
		panic("Args Error")
	} else if opCode == "0" && Name[0] == "" {
		panic("Args Error")
	}
	var err error
	for i := 0; i < 3; i++ {
		_, err = RegisterSocket.Write(body)
		if err != nil {
			panic(err)
		}
		err = RegisterSocket.SetReadDeadline(time.Now().Add(time.Millisecond * 460))
		if err != nil {
			fmt.Println(err)
			continue
		}
		n, _, err := RegisterSocket.ReadFrom(res)
		if err == nil {
			ackByte := res[12:16]
			ack := decode32(ackByte)
			if ack == seq && string(res[0:12]) == "REGIST\r\nACK:" {
				return res[16:n]
			}
		}
	}
	return nil
}
