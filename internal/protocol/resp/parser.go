package resp

import (
	"fmt"
	"io"
	"strconv"
)

func (r *Resp) ParseRESP() (Value, error) {
	output := Value{}
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

func (r *Resp) readInterger() (Value, error) {
	v := Value{Typ: "integer"}
	num, _, err := r.readNum()
	if err != nil {
		return v, err
	}
	v.In = num
	return v, nil
}

func (r *Resp) readString() (Value, error) {
	v := Value{Typ: "string"}
	line, _, err := r.ReadLine()
	if err != nil {
		return v, err
	}
	v.Str = string(line)
	return v, nil
}

func (r *Resp) readError() (Value, error) {
	v := Value{Typ: "error"}
	line, _, err := r.ReadLine()
	if err != nil {
		return v, err
	}
	v.Er = string(line)
	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}
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

func (r *Resp) readArray() (Value, error) {
	v := Value{}
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

	v.Array = make([]Value, length)
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

func (r *Resp) readInline(firstByte byte) (Value, error) {
	// Read the entire inline command line
	// firstByte is the first character of the command
	var line []byte
	line = append(line, firstByte)

	// Read until \r\n
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return Value{}, err
		}

		if b == '\n' {
			// Remove trailing \r if present
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			break
		}
		line = append(line, b)
	}

	// Parse inline command: split by spaces into command and args
	fields := parseInlineCommand(string(line))

	// Convert to RESP array format
	v := Value{Typ: "array"}
	v.Array = make([]Value, len(fields))
	for i, field := range fields {
		v.Array[i] = Value{Typ: "bulk", Bulk: field}
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

func (r *Resp) WriteBytes() {
	r.writer.Flush()
}
