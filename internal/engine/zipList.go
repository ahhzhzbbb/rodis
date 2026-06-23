package engine

import (
	"encoding/binary"
)

const HEADER_SIZE = 10

type ZipList struct {
	buf []byte
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

func (zl *ZipList) GetHeaderSize() int {
	return HEADER_SIZE
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
	content := element
	encoding := len(content)

	var prevLen byte = 0

	if zl.Length() > 0 {
		prevLen = uint8(zl.buf[zl.GetTail()+1]) + 2
	}

	zl.buf = zl.buf[:len(zl.buf)-1]

	zl.buf = append(zl.buf, prevLen)
	zl.buf = append(zl.buf, byte(encoding))
	zl.buf = append(zl.buf, content...)
	zl.buf = append(zl.buf, 0xFF)

	oldLen := binary.LittleEndian.Uint32(zl.buf[0:4])
	newTail := oldLen - 1
	newLen := zl.Length() + 1

	zl.UpdateHeader(newTail, newLen)
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
		// fmt.Printf("Offset: %d\n", offset)
		if zl.buf[offset] == 0xFF {
			break
		}
		encoding := uint8(zl.buf[offset+1])
		res = append(res, string(zl.buf[offset+2:offset+2+uint32(encoding)]))
		offset = offset + 2 + uint32(encoding)
	}
	return res
}

func (zl *ZipList) GetIndexOfElement(element string) (int, bool) {
	var offset uint32
	offset = 10
	index := 0
	for {
		if zl.buf[offset] == 0xFF {
			break
		}
		encoding := uint8(zl.buf[offset+1])
		if string(zl.buf[offset+2:offset+2+uint32(encoding)]) == element {
			return index, true
		}
		offset = offset + 2 + uint32(encoding)
		index++
	}
	return -1, false
}

func (zl *ZipList) SplitList(index int) *ZipList {
	if zl.Length() <= 1 {
		return nil
	}
	if index < 1 || index > int(zl.Length())-1 {
		return nil
	}
	var offset uint32
	var newTail uint32
	offset = HEADER_SIZE
	for range index {
		encoding := uint8(zl.buf[offset+1])
		newTail = offset
		offset = offset + 2 + uint32(encoding)
	}

	entries := zl.buf[offset : zl.GetBytes()-1]
	newZl := &ZipList{
		buf: make([]byte, HEADER_SIZE+len(entries)+1),
	}
	copy(newZl.buf[HEADER_SIZE:], entries)
	newZl.buf[len(newZl.buf)-1] = 0xFF
	newZl.UpdateHeader(uint32(len(newZl.buf)), uint16(len(entries)/2))
	zl.buf = zl.buf[:offset]
	zl.buf = append(zl.buf, 0xFF)
	zl.UpdateHeader(newTail, uint16(index))
	return newZl
}

func (zl *ZipList) Insert(index int, element string) bool {
	if index < 0 || index > int(zl.Length()) {
		return false
	}
	if index == 0 {
		zl.PushFront(element)
		return true
	}
	if index == int(zl.Length()) {
		zl.PushBack(element)
		return true
	}
	var offset uint32
	offset = HEADER_SIZE
	for range index {
		encodingEntry := uint8(zl.buf[offset+1])
		offset = offset + 2 + uint32(encodingEntry)
	}

	content := []byte(element)
	var prevLen uint8 = 0
	encoding := uint8(len(content))

	entry := []byte{
		prevLen,
		encoding,
	}

	entry = append(entry, content...)

	oldEntries := make([]byte, zl.GetBytes()-offset)
	copy(oldEntries, zl.buf[offset:])

	zl.buf = zl.buf[:offset]
	zl.buf = append(zl.buf, entry...)

	zl.buf = append(zl.buf, oldEntries...)
	newTail := zl.GetTail() + uint32(2+encoding)
	zl.UpdateHeader(newTail, zl.Length()+1)
	return true
}
