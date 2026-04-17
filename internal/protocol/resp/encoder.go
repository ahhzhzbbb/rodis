package resp

import "strconv"

func (r *Resp) Marshal(v Value) []byte {
	switch v.Typ {
	case "array":
		return r.marshalArray(v)
	case "bulk":
		return r.marshalBulk(v)
	case "string":
		return r.marshalString(v)
	case "null":
		return r.marshallNullBulk()
	case "error":
		return r.marshalError(v)
	case "integer":
		return r.marshalInteger(v)
	default:
		return nil
	}
}

func (r *Resp) marshalArray(v Value) []byte {
	// r.writer.WriteByte(ARRAY)
	// r.writer.WriteString(strconv.Itoa(len(v.Array)))
	// r.writer.WriteString("\r\n")
	// for _, element := range v.Array {
	// 	r.Marshal(element)
	// }
	var buf []byte

	buf = append(buf, ARRAY)
	buf = append(buf, strconv.Itoa(len(v.Array))...)
	buf = append(buf, '\r', '\n')
	for _, element := range v.Array {
		buf = append(buf, r.Marshal(element)...)
	}

	return buf
}

func (r *Resp) marshalBulk(v Value) []byte {
	var buf []byte

	buf = append(buf, BULK)
	buf = append(buf, strconv.Itoa(len(v.Bulk))...)
	buf = append(buf, '\r', '\n')
	buf = append(buf, v.Bulk...)
	buf = append(buf, '\r', '\n')

	return buf
}

func (r *Resp) marshalString(v Value) []byte {
	// r.writer.WriteByte(STRING)
	// r.writer.WriteString(v.Str)
	// r.writer.WriteString("\r\n")
	var buf []byte

	buf = append(buf, STRING)
	buf = append(buf, []byte(v.Str)...)
	buf = append(buf, '\r', '\n')

	return buf
}

func (r *Resp) marshallNullBulk() []byte {
	return []byte("$-1\r\n")
}

func (r *Resp) marshalError(v Value) []byte {
	// r.writer.WriteByte(ERROR)
	// r.writer.WriteString(v.Er)
	// r.writer.WriteString("\r\n")
	var buf []byte

	buf = append(buf, ERROR)
	buf = append(buf, []byte(v.Er)...)
	buf = append(buf, '\r', '\n')

	return buf
}

func (r *Resp) marshalInteger(v Value) []byte {
	// r.writer.WriteByte(INTEGER)
	// r.writer.WriteString(strconv.Itoa(v.In))
	// r.writer.WriteString("\r\n")

	var buf []byte

	buf = append(buf, INTEGER)
	buf = append(buf, []byte(strconv.Itoa(v.In))...)
	buf = append(buf, '\r', '\n')

	return buf
}
