package command

import (
	"rodis/internal/protocol/resp"
	"strconv"
	"time"
)

type ExpireCommand struct{}

func (c *ExpireCommand) Execute(args []resp.Value, ctx *CommandContext) resp.Value {
	if len(args) == 0 || len(args) == 1 {
		return resp.NewBulk("(error) ERR wrong number of arguments for 'expire' command")
	}

	var result int
	key := args[0].Bulk
	i64, _ := strconv.ParseInt((args[1].Bulk), 10, 64)
	t := time.Now().Add(time.Duration(i64) * time.Second)

	_, ok := ctx.kv.Kv[key]
	if ok {
		ctx.et.Mu.Lock()
		ctx.et.Et[key] = t
		ctx.et.Mu.Unlock()
		result++
	}

	return resp.NewInteger(result)
}
