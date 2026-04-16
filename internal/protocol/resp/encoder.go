package resp

import "strconv"

func (r *Resp) Marshal(v Value) {
	switch v.Typ {
	case "array":
		r.marshalArray(v)
	case "bulk":
		r.marshalBulk(v)
	case "string":
		r.marshalString(v)
	case "null":
		r.marshallNullBulk()
	case "error":
		r.marshalError(v)
	case "integer":
		r.marshalInteger(v)
	default:
	}
}

func (r *Resp) marshalArray(v Value) {
	r.writer.WriteByte(ARRAY)
	r.writer.WriteString(strconv.Itoa(len(v.Array)))
	r.writer.WriteString("\r\n")
	for _, element := range v.Array {
		r.Marshal(element)
	}
}

func (r *Resp) marshalBulk(v Value) {
	r.writer.WriteByte(BULK)
	r.writer.WriteString(strconv.Itoa(len(v.Bulk)))
	r.writer.WriteString("\r\n")
	r.writer.WriteString(v.Bulk)
	r.writer.WriteString("\r\n")
}

func (r *Resp) marshalString(v Value) {
	r.writer.WriteByte(STRING)
	r.writer.WriteString(v.Str)
	r.writer.WriteString("\r\n")
}

func (r *Resp) marshallNullBulk() {
	r.writer.WriteString("$-1\r\n")
}

func (r *Resp) marshalError(v Value) {
	r.writer.WriteByte(ERROR)
	r.writer.WriteString(v.Er)
	r.writer.WriteString("\r\n")
}

func (r *Resp) marshalInteger(v Value) {
	r.writer.WriteByte(INTEGER)
	r.writer.WriteString(strconv.Itoa(v.In))
	r.writer.WriteString("\r\n")
}
