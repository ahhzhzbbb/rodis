package command

import (
	"rodis/internal/protocol/resp"
)

type DelCommand struct{}

func (c *DelCommand) Execute(args []resp.Value, ctx *CommandContext) resp.Value {
	var count int
	if len(args) == 0 {
		return resp.NewError("ERR wrong number of arguments for 'del' command")
	}

	for _, arg := range args {
		key := arg.Bulk

		if ctx.k.Del(key) {
			count++
		}
	}
	return resp.NewInteger(count)
}
