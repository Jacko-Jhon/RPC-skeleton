package RUDP

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

func encodeFlag(n int32) []byte {
	code := make([]byte, 3)
	code[0] = byte(n >> 16)
	code[1] = byte(n >> 8)
	code[2] = byte(n)
	return code
}

func decodeFlag(code []byte) int32 {
	return int32(code[0])<<16 | int32(code[1])<<8 | int32(code[2])
}
