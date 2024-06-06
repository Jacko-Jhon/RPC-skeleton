package RUDP

import (
	"math/rand"
	"net"
)

func (r *RDP) Listen(flag int32, IP string, Port int) {
	r.flag = flag
	r.listenAddr = &net.UDPAddr{
		IP:   net.ParseIP(IP),
		Port: Port,
	}
	udp, err := net.ListenUDP("udp", r.listenAddr)
	if err != nil {
		panic(err)
	} else {
		r.socket = udp
	}
}

func (r *RDP) ReadFrom() (*net.UDPAddr, *Packet) {
	buf := make([]byte, 0, 1040)
	for {
		n, addr, err := r.socket.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}
		if buf[3] == 0 {
			l := decode32(buf[8:12])
			if l == 1 {
				flag := decodeFlag(buf[0:3])
				return addr, &Packet{
					Flag: flag,
					Data: buf[16:n],
				}
			}
			b := []byte(addr.String())
			var sum int32 = 0
			for i := 0; i < len(b); i++ {
				sum += int32(b[i])
			}
			seq := decode32(buf[4:8])
			var f = (int64(sum) << 31) | int64(seq)
			if b, ok := r.bufferCache[f]; ok {
				r.sendACK(addr, seq, &b)
				if b.add(buf[8:n]) {
					flag := decodeFlag(buf[0:3])
					return addr, &Packet{
						Flag: flag,
						Data: b.data[:b.len],
					}
				}
			} else {
				r.bufferCache[f] = *newBuffer(buf[8:n])
			}
		} else if buf[3] == 2 {
			// ack
		}
	}
}

func (r *RDP) SendTo(addr *net.UDPAddr, packet *Packet) {
	i := 0
	h := make([]byte, 16)
	f := encodeFlag(packet.Flag)
	for i = 0; i <= 2; i++ {
		h[i] = f[i]
	}
	h[3] = 0
	seq := rand.Int31()
	b := encode32(seq)
	for i = 4; i <= 7; i++ {
		h[i] = b[i-4]
	}
	n := len(packet.Data) / 1024
	if len(packet.Data)%1024 != 0 {
		n++
	}
	b = encode32(int32(n))
	for i = 8; i <= 11; i++ {
		h[i] = b[i-8]
	}
	if n < 3 {
		for i = 1; i < n; i++ {
			h[15] = uint8(i)
			l := i*1024 - len(packet.Data)
			if l < 0 {

			} else {
				l = i * 1024
			}
			//_, err := r.socket.WriteToUDP(bytes, addr)
			//if err != nil {
			//	return
			//}
		}
	}
}

func (r *RDP) sendACK(addr *net.UDPAddr, seq int32, b *buffer) {

}

func (r *RDP) singleSend() {

}
