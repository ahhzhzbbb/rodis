package command

import "rodis/internal/protocol/resp"

type LpushCommand struct{}

func (c *LpushCommand) Execute(args []resp.Payload, ctx *CommandContext) resp.Payload {
	if len(args) < 2 {
		return resp.NewError("ERR wrong number of arguments for 'rpush' command")
	}

	if ctx == nil || ctx.k == nil {
		return resp.NewError("ERR internal server error")
	}

	key := args[0].Bulk

	elements := make([]string, len(args)-1)

	for i := len(args) - 1; i > 0; i-- {
		elements[len(args)-i-1] = args[i].Bulk
	}

	res, err := ctx.k.SetList(key, true, elements)
	if err != nil {
		return resp.NewError(err.Error())
	}

	return resp.NewInteger(res)
}
