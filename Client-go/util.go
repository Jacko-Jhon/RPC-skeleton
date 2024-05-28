package main

import (
	"fmt"
	"math/rand"
	"time"
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

func DailRegistry(opCode string, Name ...string) []byte {
	seq := rand.Int31()
	if seq > 100 {
		seq -= 50
	}
	h := []byte("CLIENT\r\nSEQ:")
	seqB := encode32(seq)
	op := []byte("\r\nOP:" + opCode)
	body := append(h, seqB...)
	body = append(body, op...)
	res := make([]byte, 1024)
	if opCode == "0" && Name[0] != "" {
		body = append(body, []byte("\r\nNAME:"+Name[0])...)
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
			if ack == seq+1 && string(res[0:12]) == "REGIST\r\nACK:" {
				return res[16:n]
			}
		}
	}
	return nil
}
