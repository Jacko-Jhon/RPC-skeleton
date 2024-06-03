package RUDP

func (b *buffer) add([]byte) bool {
	return true
}

func newBuffer(buf []byte) *buffer {
	m := decode32(buf[0:4])
	idx := decode32(buf[4:8])
	flag := make([]bool, m)
	data := make([]byte, m*1024)
	flag[idx] = true
	i := (idx - 1) * 1024
	for j := 8; j < len(buf); j++ {
		data[i] = buf[j]
		i++
	}
	return &buffer{
		len:      len(buf) - 8,
		unsolved: int(m - 1),
		flag:     flag,
		data:     data,
	}
}
