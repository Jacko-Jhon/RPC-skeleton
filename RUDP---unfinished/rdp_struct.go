package RUDP

import "net"

type RDP struct {
	flag        int32
	listenAddr  *net.UDPAddr
	socket      *net.UDPConn
	bufferCache map[int64]buffer
	ack         map[int32]bool
}

type RDPDialer struct {
	dialAddr *net.UDPAddr
	socket   *net.UDPConn
}

type buffer struct {
	len      int     // 字节总长度，每次读完就自增
	unsolved int     // 未收到的块的数量
	unACK    int     // 未发送的Ack量
	rec      []int32 // 最近接收的块序号
	flag     []bool  // 块序号的标记
	data     []byte  // 数据
}

type Packet struct {
	Flag int32
	Data []byte
}
