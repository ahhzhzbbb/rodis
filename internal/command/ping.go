package command

import "rodis/internal/protocol/resp"

type PingCommand struct{}

func (c *PingCommand) Execute(args []resp.Value, ctx *CommandContext) resp.Value {
	if len(args) > 1 {
		return resp.NewError("ERR wrong number of arguments for 'ping' command")
	}
	if len(args) == 1 {
		return args[0]
	}

	return resp.NewBulk("PONG")
}
