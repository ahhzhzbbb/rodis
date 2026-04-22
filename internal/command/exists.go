package command

import "rodis/internal/protocol/resp"

type ExistsCommand struct{}

func (c *ExistsCommand) Execute(args []resp.Payload, ctx *CommandContext) resp.Payload {
	if len(args) == 0 {
		return resp.NewError("ERR wrong number of arguments for 'exists' command")
	}

	if ctx == nil || ctx.k == nil {
		return resp.NewError("ERR internal server error")
	}

	var result int

	for _, arg := range args {
		_, ok := ctx.k.Get(arg.Bulk)
		if ok {
			result++
		}
	}
	return resp.NewInteger(result)
}
