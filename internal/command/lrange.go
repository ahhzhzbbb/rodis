package command

import (
	"rodis/internal/protocol/resp"
)

type LrangeCommand struct{}

func (c *LrangeCommand) Execute(args []resp.Payload, ctx *CommandContext) resp.Payload {
	if len(args) != 3 {
		return resp.NewError("ERR wrong number of arguments for 'lrange' command")
	}

	if ctx == nil || ctx.k == nil {
		return resp.NewError("ERR internal server error")
	}

	key := args[0].Bulk
	start := args[1].Bulk
	stop := args[2].Bulk

	elements, found, err := ctx.k.GetListBetween(key, start, stop)
	if err != nil {
		return resp.NewError(err.Error())
	}
	if !found {
		return resp.NewArray(nil)
	}
	var res []resp.Payload
	for _, e := range elements {
		res = append(res, resp.NewBulk(e))
	}
	return resp.NewArray(res)
}
