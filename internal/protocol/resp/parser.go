package resp

import (
	"fmt"
	"io"
	"strconv"
)

func (r *Resp) ParseRESP() (Value, error) {
	fmt.Println("parsing request to value...")
	output := Value{}
	firstByte, err := r.reader.ReadByte()
	if err != nil {
		return output, err
	}
	// fmt.Printf("first byte: %d (%q)\n", firstByte, firstByte)

	switch firstByte {
	case STRING:
		return r.readString()
	case ERROR:
		return r.readError()
	case INTEGER:
		return r.readInterger()
	case BULK:
		return r.readBulk()
	case ARRAY:
		return r.readArray()
	default:
		return output, fmt.Errorf("resp: unsupported value type %q", firstByte)
	}
}

func (r *Resp) ReadLine() ([]byte, int, error) {
	var size int
	line := make([]byte, 0)

	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}

		size++
		line = append(line, b)

		if len(line) >= 2 &&
			line[len(line)-2] == '\r' &&
			line[len(line)-1] == '\n' {

			line = line[:len(line)-2]
			break
		}
	}
	return line, size, nil
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
	v := Value{Typ: "interger"}
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
		return v, nil
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
