package resp

func Commands(value Value) Value {
	values := value.array
	command := values[0]
	switch command.bulk {
	case "PING":
		return ping(values)
	}
	return Value{}
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'ping' command"}
	}
	return Value{typ: "string", str: "PONG"}
}
