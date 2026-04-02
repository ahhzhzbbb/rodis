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
	Typ   string
	Str   string
	Er    string
	In    int
	Bulk  string
	Array []Value
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

func NewError(msg string) Value {
	return Value{
		Typ: "error",
		Str: msg,
	}
}

func NewString(msg string) Value {
	return Value{
		Typ: "string",
		Str: msg,
	}
}

func NewNum(msg int) Value {
	return Value{
		Typ: "num",
		In:  msg,
	}
}

func NewBulk(msg string) Value {
	return Value{
		Typ:  "bulk",
		Bulk: msg,
	}
}

// func NewArray(values []Value) {
// 	return []
// }
