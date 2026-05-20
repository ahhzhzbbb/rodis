package command

import "rodis/internal/protocol/resp"

type RpushCommand struct{}

func (c *RpushCommand) Execute(args []resp.Payload, ctx *CommandContext) resp.Payload {
	if len(args) < 2 {
		return resp.NewError("ERR wrong number of arguments for 'rpush' command")
	}

	if ctx == nil || ctx.k == nil {
		return resp.NewError("ERR internal server error")
	}

	key := args[0].Bulk

	elements := make([]string, len(args)-1)

	for i := 1; i < len(args); i++ {
		elements[i-1] = args[i].Bulk
	}

	res, err := ctx.k.SetList(key, false, elements)
	if err != nil {
		return resp.NewError(err.Error())
	}

	return resp.NewInteger(res)
}
