package core

import (
	"encoding/binary"
	"fmt"
	"math"
)

func ReadUint8(buffer []byte, offset *int) uint8 {
	ret := buffer[*offset]
	*offset++
	return ret
}
func ReadUInt16(buffer []byte, offset *int) uint16 {
	ret := binary.LittleEndian.Uint16(buffer[*offset:])
	*offset += 2
	return ret
}
func ReadUInt32(buffer []byte, offset *int) uint32 {
	ret := binary.LittleEndian.Uint32(buffer[*offset:])
	*offset += 4
	return ret
}
func ReadUInt64(buffer []byte, offset *int) uint64 {
	ret := binary.LittleEndian.Uint64(buffer[*offset:])
	*offset += 8
	return ret
}

func ReadInt32(buffer []byte, offset *int) int32 {
	return int32(ReadUInt32(buffer, offset))
}

func ReadFloat32(buffer []byte, offset *int) float32 {
	bits := ReadUInt32(buffer, offset)
	return math.Float32frombits(bits)
}

func ReadString(bytes []byte, offset *int) string {
	len := 0
	tmp := bytes[*offset:]
	for _, v := range tmp {
		*offset++
		if v == byte(0) {
			break
		}
		len++
	}
	return string(tmp[:len])
}

func SearchHeader(bytes []byte, header []byte, end int) int {
	step := len(header)

	for i := 0; i < end-step; i++ {
		for CompareByte(bytes[i:i+step], header) {
			return i + step
		}
	}
	return -1

}
func CompareByte(bytes1 []byte, bytes2 []byte) bool {
	len1 := len(bytes1)
	len2 := len(bytes2)
	if len1 != len2 {
		return false
	}
	for i := 0; i < len1; i++ {

		if bytes1[i] != bytes2[i] {
			return false
		}

	}
	return true

}

func PrintHex(bytes []byte) {
	r := ""
	for i, b := range bytes {
		r = r + fmt.Sprintf("%x ", b)
		if i%16 == 0 && i > 0 {
			fmt.Println(r)
			r = ""
		}
	}

}
