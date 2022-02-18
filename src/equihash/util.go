package equihash

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"unsafe"
)

func fromBytes(in []byte) []uint32 {
	r := bytes.NewReader(in)
	out := make([]uint32, len(in)/4)
	binary.Read(r, binary.LittleEndian, &out)
	return out
}

func toBytes(in []uint32) []byte {
	ptr := unsafe.Pointer(&in)
	hdr := *(*reflect.SliceHeader)(ptr)

	hdr.Len *= 4
	hdr.Cap *= 4

	return *(*[]byte)(unsafe.Pointer(&hdr))
}
