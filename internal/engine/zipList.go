package engine

import (
	"encoding/binary"
	"sync"
)

type ZipList struct {
	buf []byte
	mu  sync.RWMutex
	//layout: [zlbytes][zltail][zllen][entries][END]
}

func NewZipList() *ZipList {
	zl := &ZipList{
		buf: make([]byte, 11),
	}

	binary.LittleEndian.PutUint32(
		zl.buf[0:4],
		11,
	)

	binary.LittleEndian.PutUint32(
		zl.buf[4:8],
		10,
	)

	binary.LittleEndian.PutUint16(
		zl.buf[8:10],
		0,
	)

	zl.buf[10] = 0xFF

	return zl
}

func (zl *ZipList) Length() uint16 {
	return binary.LittleEndian.Uint16(zl.buf[8:10])
}

func (zl *ZipList) GetTail() uint32 {
	return binary.LittleEndian.Uint32(zl.buf[4:8])
}

func (zl *ZipList) updateHeader() {
	// oldLen := binary.LittleEndian.Uint32(zl.buf[0:4])
	binary.LittleEndian.PutUint32(
		zl.buf[0:4],
		uint32(len(zl.buf)),
	)

	// binary.LittleEndian.PutUint32(
	// 	zl.buf[4:8],
	// 	uint32(oldLen-1),
	// )

	// count := zl.Length()

	// binary.LittleEndian.PutUint16(
	// 	zl.buf[8:10],
	// 	uint16(count+1),
	// )
}

func (zl *ZipList) PushBack(element string) {
	content := []byte(element)

	encoding := byte(len(content))

	var prevLen uint8 = 0

	if zl.Length() > 0 {
		prevLen = uint8(zl.buf[zl.GetTail()+1]) + 2
	}

	entry := []byte{
		prevLen,
		encoding,
	}

	entry = append(entry, content...)

	zl.buf = zl.buf[:len(zl.buf)-1]

	zl.buf = append(zl.buf, entry...)
	zl.buf = append(zl.buf, 0xFF)

	oldLen := binary.LittleEndian.Uint32(zl.buf[0:4])
	zl.updateHeader()
	binary.LittleEndian.PutUint32(
		zl.buf[4:8],
		uint32(oldLen-1),
	)
	count := zl.Length()

	binary.LittleEndian.PutUint16(
		zl.buf[8:10],
		uint16(count+1),
	)
}

// func (zl *ZipList) PushFront() {

// }

func (zl *ZipList) GetElements() (res []string) {
	var offset uint32
	offset = 10
	for {
		if zl.buf[offset] == 0xFF {
			break
		}
		encoding := uint8(zl.buf[offset+1])
		res = append(res, string(zl.buf[offset+2:offset+2+uint32(encoding)]))
	}
	return res
}

func (zl *ZipList) ReverseTravel() {

}

// func (z *ZipList) next(pos int) int {

// }
