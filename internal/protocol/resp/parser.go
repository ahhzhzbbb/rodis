package resp

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func (r *Resp) ParseRESP() (Value, error) {
	output := Value{}
	firstByte, err := r.reader.ReadByte()
	if err != nil {
		return output, err
	}

	switch firstByte {
	case STRING:
		return r.readString()
	case ERROR:
		return r.readError()
	case INTERGER:
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

		size += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' && line[len(line)-1] == '\n' {
			break
		}
	}
	return line, size, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.ReadLine()

	if err != nil {

		return 0, 0, err
	}
	temp := strings.TrimSuffix(string(line), "\r\n")
	i64, err := strconv.ParseInt(temp, 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return int(i64), n, nil
}

func (r *Resp) readInterger() (Value, error) {
	v := Value{typ: "interger"}
	num, _, err := r.readInteger()
	if err != nil {
		return v, err
	}
	v.num = num
	return v, nil
}

func (r *Resp) readString() (Value, error) {
	v := Value{typ: "string"}
	line, _, err := r.ReadLine()
	if err != nil {
		return v, err
	}
	v.str = strings.TrimSuffix(string(line), "\r\n")
	return v, nil
}

func (r *Resp) readError() (Value, error) {
	v := Value{typ: "error"}
	line, _, err := r.ReadLine()
	if err != nil {
		return v, err
	}
	v.er = strings.TrimSuffix(string(line), "\r\n")
	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"
	length, _, err := r.readInteger()
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

	v.bulk = string(bulk)

	return v, nil
}

func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}
	if length == -1 {
		return v, nil
	}
	if length < -1 {
		return v, fmt.Errorf("resp: invalid array length %d", length)
	}

	v.array = make([]Value, length)
	for i := range v.array {
		v.array[i], err = r.ParseRESP()
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
