package equihash

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

var bytesOrder = binary.LittleEndian

const u32size = int(unsafe.Sizeof(uint32(0)))

func fromBytes(in []byte) []uint32 {
	r := bytes.NewReader(in)

	out := make([]uint32, len(in)/u32size)
	binary.Read(r, bytesOrder, &out)
	return out
}

func toBytes(in []uint32) []byte {
	b := make([]byte, len(in)*u32size)

	for i, u := range in {
		bytesOrder.PutUint32(b[u32size*i:], uint32(u))
	}

	return b
}
