package resp

import (
	"fmt"
	"io"
	"strconv"
)

func (r *Resp) ParseRESP() (Payload, error) {
	output := Payload{}
	firstByte, err := r.reader.ReadByte()
	if err != nil {
		return output, err
	}
	// fmt.Printf("first byte: %d (%q)\n", firstByte, firstByte)

	switch firstByte {
	case STRING:
		// fmt.Printf("%c\n", STRING)
		return r.readString()
	case ERROR:
		// fmt.Printf("%c\n", ERROR)
		return r.readError()
	case INTEGER:
		// fmt.Printf("%c\n", INTEGER)
		return r.readInterger()
	case BULK:
		// fmt.Printf("%c\n", BULK)
		return r.readBulk()
	case ARRAY:
		// fmt.Printf("%c\n", ARRAY)
		return r.readArray()
	default:
		return r.readInline(firstByte)
	}
}

func (r *Resp) ReadLine() ([]byte, int, error) {
	line, err := r.reader.ReadSlice('\n')
	if err != nil {
		return nil, 0, err
	}

	if len(line) >= 2 && line[len(line)-2] == '\r' {
		line = line[:len(line)-2]
	}

	return line, len(line), nil
}

func (r *Resp) readNum() (x int, n int, err error) {
	num, n, err := r.ReadLine()
	if err != nil {
		return 0, 0, err
	}

	i64, err := strconv.ParseInt(string(num), 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return int(i64), n, nil
}

func (r *Resp) readInterger() (Payload, error) {
	v := Payload{Typ: "integer"}
	num, _, err := r.readNum()
	if err != nil {
		return v, err
	}
	v.In = num
	return v, nil
}

func (r *Resp) readString() (Payload, error) {
	v := Payload{Typ: "string"}
	line, _, err := r.ReadLine()
	if err != nil {
		return v, err
	}
	v.Str = string(line)
	return v, nil
}

func (r *Resp) readError() (Payload, error) {
	v := Payload{Typ: "error"}
	line, _, err := r.ReadLine()
	if err != nil {
		return v, err
	}
	v.Er = string(line)
	return v, nil
}

func (r *Resp) readBulk() (Payload, error) {
	v := Payload{}
	v.Typ = "bulk"
	length, _, err := r.readNum()
	if err != nil {
		return v, err
	}

	if length == -1 {
		return NewNullBulk(), nil
	}
	if length < -1 {
		return v, fmt.Errorf("resp: invalid bulk length %d", length)
	}

	bulk := make([]byte, length)
	if _, err := io.ReadFull(r.reader, bulk); err != nil {
		return v, err
	}
	if err := r.readCRLF(); err != nil {
		return v, err
	}

	v.Bulk = string(bulk)

	return v, nil
}

func (r *Resp) readArray() (Payload, error) {
	v := Payload{}
	v.Typ = "array"

	length, _, err := r.readNum()
	if err != nil {
		return v, err
	}
	if length == -1 {
		return v, nil
	}
	if length < -1 {
		return v, fmt.Errorf("resp: invalid array length %d", length)
	}

	v.Array = make([]Payload, length)
	for i := range v.Array {
		v.Array[i], err = r.ParseRESP()
		if err != nil {
			return v, err
		}
	}
	return v, nil
}

func (r *Resp) readCRLF() error {
	crlf := make([]byte, 2)
	if _, err := io.ReadFull(r.reader, crlf); err != nil {
		return err
	}
	if crlf[0] != '\r' || crlf[1] != '\n' {
		return fmt.Errorf("resp: expected CRLF after bulk string")
	}
	return nil
}

func (r *Resp) readInline(firstByte byte) (Payload, error) {
	var line []byte
	line = append(line, firstByte)

	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return Payload{}, err
		}

		if b == '\n' {
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			break
		}
		line = append(line, b)
	}

	fields := parseInlineCommand(string(line))

	v := Payload{Typ: "array"}
	v.Array = make([]Payload, len(fields))
	for i, field := range fields {
		v.Array[i] = Payload{Typ: "bulk", Bulk: field}
	}

	return v, nil
}

func parseInlineCommand(line string) []string {
	var fields []string
	var current []byte

	for i := 0; i < len(line); i++ {
		ch := line[i]
		if ch == ' ' || ch == '\t' {
			if len(current) > 0 {
				fields = append(fields, string(current))
				current = nil
			}
		} else {
			current = append(current, ch)
		}
	}
	if len(current) > 0 {
		fields = append(fields, string(current))
	}

	return fields
}

func (r *Resp) WriteBytes(bytes []byte) error {
	_, err := r.writer.Write(bytes)
	return err
}

func (r *Resp) FlushWriter() error {
	return r.writer.Flush()
}
