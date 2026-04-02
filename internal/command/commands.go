package command

import "rodis/internal/protocol/resp"

func Commands(value resp.Value) resp.Value {
	values := value.Array
	command := values[0]
	switch command.Bulk {
	case "PING":
		return ping(values)
	default:
		return resp.NewBulk("OK")
	}
}

func ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.NewError("ERR wrong number of arguments for 'ping' command")
	}
	return resp.NewBulk("PONG")
}
