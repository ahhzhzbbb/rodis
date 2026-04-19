package resp

import (
	"bufio"
	"fmt"
	"io"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	Typ    string
	Str    string
	Er     string
	In     int
	Bulk   string
	Array  []Value
	Inline string
}

type Resp struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewResp(rw io.ReadWriter) *Resp {
	return &Resp{reader: bufio.NewReader(rw), writer: bufio.NewWriterSize(rw, 64*1024)}
}

func (r *Resp) Writer(bytes []byte) error {
	_, err := r.writer.Write(bytes)
	if err != nil {
		fmt.Println("failed to write")
		return err
	}
	return nil
}

func NewError(msg string) Value {
	return Value{
		Typ: "error",
		Er:  msg,
	}
}

func NewString(msg string) Value {
	return Value{
		Typ: "string",
		Str: msg,
	}
}

func NewInteger(msg int) Value {
	return Value{
		Typ: "integer",
		In:  msg,
	}
}

func NewBulk(msg string) Value {
	return Value{
		Typ:  "bulk",
		Bulk: msg,
	}
}

func NewNullBulk() Value {
	return Value{
		Typ: "null",
	}
}

func NewArray(msg []Value) Value {
	return Value{
		Typ:   "array",
		Array: msg,
	}
}

func (r *Resp) HasBufferedData() bool {
	return r.reader.Buffered() > 0
}

// func NewArray(values []Value) {
// 	return []
// }
