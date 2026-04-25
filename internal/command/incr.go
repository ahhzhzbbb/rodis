package command

import (
	"rodis/internal/protocol/resp"
)

type IncrCommand struct{}

func (c *IncrCommand) Execute(args []resp.Payload, ctx *CommandContext) resp.Payload {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'incr' command")
	}
	if ctx == nil || ctx.k == nil {
		return resp.NewError("ERR internal server error")
	}

	key := args[0].Bulk
	i64, err := ctx.k.IncrString(key)
	if err != nil {
		// if errors.Is(err, engine.ErrNotInteger) {
		// 	return resp.NewError("ERR value is not an integer or out of range")
		// }
		// return resp.NewError("ERR internal server error")
		return resp.NewError(err.Error())
	}

	if i64 > int64(^uint(0)>>1) {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	return resp.NewInteger(int(i64))
}
