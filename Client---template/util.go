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
