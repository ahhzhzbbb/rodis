package command

import (
	"rodis/internal/protocol/resp"
	"strconv"
	"time"
)

type ExpireCommand struct{}

func (c *ExpireCommand) Execute(args []resp.Value, ctx *CommandContext) resp.Value {
	if len(args) != 2 {
		return resp.NewError("ERR wrong number of arguments for 'expire' command")
	}

	var result int
	key := args[0].Bulk
	i64, err := strconv.ParseInt((args[1].Bulk), 10, 64)
	if err != nil {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	t := time.Now().Add(time.Duration(i64) * time.Second)

	if ctx.k.SetExpireTime(key, t) {
		result++
	}

	return resp.NewInteger(result)
}
