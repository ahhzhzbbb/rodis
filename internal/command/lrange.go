package command

import (
	"rodis/internal/protocol/resp"
)

type LrangeCommand struct{}

func (c *LrangeCommand) Execute(args []resp.Payload, ctx *CommandContext) resp.Payload {
	// if len(args) < 2 {
	// 	return resp.NewError("ERR wrong number of arguments for 'rpush' command")
	// }

	// if ctx == nil || ctx.k == nil {
	// 	return resp.NewError("ERR internal server error")
	// }

	key := args[0].Bulk
	elements, found, err := ctx.k.GetList(key)
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
