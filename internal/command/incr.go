package command

import (
	"fmt"
	"rodis/internal/protocol/resp"
	"strconv"
)

type IncrCommand struct{}

func (c *IncrCommand) Execute(args []resp.Value, ctx *CommandContext) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'incr' command")
	}

	key := args[0].Bulk
	value, ok := ctx.k.Get(key)
	if !ok {
		ctx.k.Set(key, "1")
		return resp.NewInteger(1)
	}

	i64, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	i64++

	ctx.k.DelValue(key)
	ctx.k.Set(key, fmt.Sprint(i64))
	return resp.NewInteger(int(i64))
}
