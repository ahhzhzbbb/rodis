package resp

import "strconv"

func (v Value) Marshal() []byte {
	switch v.Typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshallNull()
	case "error":
		return v.marshalError()
	case "integer":
		return v.marshalInteger()
	default:
		return []byte{}
	}
}

func (v Value) marshalArray() []byte {
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len(v.Array))...)
	bytes = append(bytes, '\r', '\n')

	for _, element := range v.Array {
		bytes = append(bytes, element.Marshal()...)
	}

	return bytes
}
func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.Bulk))...)
	bytes = append(bytes, '\r', '\n')

	bytes = append(bytes, v.Bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}
func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}
func (v Value) marshallNull() []byte {
	return []byte("$-1\r\n")
}
func (v Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.Er...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}
func (v Value) marshalInteger() []byte {
	bytes := []byte(string(INTEGER) + strconv.Itoa(v.In) + "\r\n")
	return bytes
}
