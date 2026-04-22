package command

import (
	"rodis/internal/protocol/resp"
)

type DelCommand struct{}

func (c *DelCommand) Execute(args []resp.Payload, ctx *CommandContext) resp.Payload {
	var count int
	if len(args) == 0 {
		return resp.NewError("ERR wrong number of arguments for 'del' command")
	}

	if ctx == nil || ctx.k == nil {
		return resp.NewError("ERR internal server error")
	}

	for _, arg := range args {
		key := arg.Bulk

		if ctx.k.Del(key) {
			count++
		}
	}
	return resp.NewInteger(count)
}
