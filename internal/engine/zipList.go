package engine

import (
	"encoding/binary"
	"sync"
)

const HEADER_SIZE = 10

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

func (zl *ZipList) GetBytes() uint32 {
	return binary.LittleEndian.Uint32(zl.buf[0:4])
}

func (zl *ZipList) UpdateHeader(newTail uint32, newLen uint16) {
	binary.LittleEndian.PutUint32(
		zl.buf[0:4],
		uint32(len(zl.buf)),
	)
	binary.LittleEndian.PutUint32(
		zl.buf[4:8],
		newTail,
	)
	binary.LittleEndian.PutUint16(
		zl.buf[8:10],
		newLen,
	)
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
	binary.LittleEndian.PutUint32(
		zl.buf[0:4],
		uint32(len(zl.buf)),
	)
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

func (zl *ZipList) PopBack() string {
	if zl.Length() == 0 {
		return ""
	}
	tailList := zl.GetTail()
	prevLen := uint8(zl.buf[tailList])
	encoding := uint8(zl.buf[tailList+1])
	value := string(zl.buf[tailList+2 : tailList+2+uint32(encoding)])
	zl.buf = zl.buf[:tailList]
	zl.buf = append(zl.buf, 0xFF)
	tailList = tailList - uint32(prevLen)
	zl.UpdateHeader(tailList, zl.Length()-1)
	return value
}

func (zl *ZipList) PushFront(element string) {
	content := []byte(element)
	var prevLen uint8 = 0
	encoding := uint8(len(content))

	entry := []byte{
		prevLen,
		encoding,
	}

	entry = append(entry, content...)

	if zl.Length() == 0 {
		zl.buf = zl.buf[:HEADER_SIZE]
		zl.buf = append(zl.buf, entry...)
		zl.buf = append(zl.buf, 0xFF)
		zl.UpdateHeader(uint32(HEADER_SIZE), 1)
		return
	}

	zl.buf[HEADER_SIZE] = byte(encoding + 2)
	newTail := zl.GetTail() + uint32(encoding) + 2
	newLen := zl.Length() + 1

	oldEntries := make([]byte, zl.GetBytes()-HEADER_SIZE)
	copy(oldEntries, zl.buf[HEADER_SIZE:])

	zl.buf = zl.buf[:HEADER_SIZE]
	zl.buf = append(zl.buf, entry...)
	zl.buf = append(zl.buf, oldEntries...)

	zl.UpdateHeader(newTail, newLen)
}

func (zl *ZipList) PopFront() string {
	if zl.Length() == 0 {
		return ""
	}
	if zl.Length() == 1 {
		value := string(zl.buf[HEADER_SIZE+2 : zl.GetBytes()-1])
		zl.buf = zl.buf[:HEADER_SIZE]
		zl.buf = append(zl.buf, 0xFF)
		zl.UpdateHeader(uint32(HEADER_SIZE), 0)
		return value
	}
	popedEntryEncoding := uint8(zl.buf[HEADER_SIZE+1])
	value := string(zl.buf[HEADER_SIZE+2 : HEADER_SIZE+2+popedEntryEncoding])
	zl.buf[HEADER_SIZE+uint32(popedEntryEncoding)+2] = byte(0)

	oldEntries := make([]byte, zl.GetBytes()-HEADER_SIZE-uint32(popedEntryEncoding)-2)
	copy(oldEntries, zl.buf[HEADER_SIZE+uint32(popedEntryEncoding)+2:])

	zl.buf = zl.buf[:HEADER_SIZE]
	zl.buf = append(zl.buf, oldEntries...)

	newTail := zl.GetTail() - uint32(popedEntryEncoding) - 2
	newLen := zl.Length() - 1
	zl.UpdateHeader(newTail, newLen)

	return value
}

func (zl *ZipList) GetElements() (res []string) {
	var offset uint32
	offset = 10
	for {
		if zl.buf[offset] == 0xFF {
			break
		}
		encoding := uint8(zl.buf[offset+1])
		res = append(res, string(zl.buf[offset+2:offset+2+uint32(encoding)]))
		offset = offset + 2 + uint32(encoding)
	}
	return res
}
