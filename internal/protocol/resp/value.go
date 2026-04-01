package resp

import (
	"bufio"
	"fmt"
	"io"
)

const (
	STRING   = '+'
	ERROR    = '-'
	INTERGER = ':'
	BULK     = '$'
	ARRAY    = '*'
)

type Value struct {
	typ   string
	str   string
	er    string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewResp(rw io.ReadWriter) *Resp {
	return &Resp{reader: bufio.NewReader(rw), writer: bufio.NewWriter(rw)}
}

func (r *Resp) Writer(bytes []byte) error {
	_, err := r.writer.Write(bytes)
	if err != nil {
		fmt.Println("failed to write")
		return err
	}
	return r.writer.Flush()
}
