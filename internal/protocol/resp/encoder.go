package resp

import "strconv"

func (r *Resp) Marshal(v Value) error {
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

func (r *Resp) marshalArray(v Value) error {
	if err := r.writer.WriteByte(ARRAY); err != nil {
		return err
	}
	if _, err := r.writer.WriteString(strconv.Itoa(len(v.Array))); err != nil {
		return err
	}
	if _, err := r.writer.WriteString("\r\n"); err != nil {
		return err
	}
	for _, element := range v.Array {
		if err := r.Marshal(element); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resp) marshalBulk(v Value) error {
	if err := r.writer.WriteByte(BULK); err != nil {
		return err
	}
	if _, err := r.writer.WriteString(strconv.Itoa(len(v.Bulk))); err != nil {
		return err
	}
	if _, err := r.writer.WriteString("\r\n"); err != nil {
		return err
	}
	if _, err := r.writer.WriteString(v.Bulk); err != nil {
		return err
	}
	if _, err := r.writer.WriteString("\r\n"); err != nil {
		return err
	}
	return nil
}

func (r *Resp) marshalString(v Value) error {
	if err := r.writer.WriteByte(STRING); err != nil {
		return err
	}
	if _, err := r.writer.WriteString(v.Str); err != nil {
		return err
	}
	if _, err := r.writer.WriteString("\r\n"); err != nil {
		return err
	}
	return nil
}

func (r *Resp) marshallNullBulk() error {
	_, err := r.writer.WriteString("$-1\r\n")
	return err
}

func (r *Resp) marshalError(v Value) error {
	if err := r.writer.WriteByte(ERROR); err != nil {
		return err
	}
	if _, err := r.writer.WriteString(v.Er); err != nil {
		return err
	}
	if _, err := r.writer.WriteString("\r\n"); err != nil {
		return err
	}
	return nil
}

func (r *Resp) marshalInteger(v Value) error {
	if err := r.writer.WriteByte(INTEGER); err != nil {
		return err
	}
	if _, err := r.writer.WriteString(strconv.Itoa(v.In)); err != nil {
		return err
	}
	if _, err := r.writer.WriteString("\r\n"); err != nil {
		return err
	}
	return nil
}
