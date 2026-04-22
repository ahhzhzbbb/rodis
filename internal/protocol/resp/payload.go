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

type Payload struct {
	Typ    string
	Str    string
	Er     string
	In     int
	Bulk   string
	Array  []Payload
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

func NewError(msg string) Payload {
	return Payload{
		Typ: "error",
		Er:  msg,
	}
}

func NewString(msg string) Payload {
	return Payload{
		Typ: "string",
		Str: msg,
	}
}

func NewInteger(msg int) Payload {
	return Payload{
		Typ: "integer",
		In:  msg,
	}
}

func NewBulk(msg string) Payload {
	return Payload{
		Typ:  "bulk",
		Bulk: msg,
	}
}

func NewNullBulk() Payload {
	return Payload{
		Typ: "null",
	}
}

func NewArray(msg []Payload) Payload {
	return Payload{
		Typ:   "array",
		Array: msg,
	}
}

func (r *Resp) HasBufferedData() bool {
	return r.reader.Buffered() > 0
}

// func NewArray(Payloads []Payload) {
// 	return []
// }
