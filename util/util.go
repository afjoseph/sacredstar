package util

func Uint32ToByteSlice(u uint32) []byte {
	return []byte{
		byte(u >> 24),
		byte(u >> 16),
		byte(u >> 8),
		byte(u),
	}
}

func ByteSliceToUint32(b []byte) uint32 {
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}
