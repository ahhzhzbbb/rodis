package command

import (
	"rodis/internal/protocol/resp"
	"time"
)

type GetCommand struct{}

func (c *GetCommand) Execute(args []resp.Value, ctx *CommandContext) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'get' command")
	}
	if ctx == nil || ctx.kv == nil {
		return resp.NewError("ERR internal server error")
	}

	key := args[0].Bulk

	ctx.et.Mu.RLock()
	t, ok := ctx.et.Et[key]
	ctx.et.Mu.RUnlock()

	if ok && t.Before(time.Now()) {
		ctx.kv.Mu.Lock()
		delete(ctx.kv.Kv, key)
		ctx.kv.Mu.Unlock()
		ctx.et.Mu.Lock()
		delete(ctx.et.Et, key)
		ctx.et.Mu.Unlock()
		return resp.NewBulk("")
	}

	ctx.kv.Mu.RLock()
	value, exists := ctx.kv.Kv[key]
	ctx.kv.Mu.RUnlock()
	if !exists {
		return resp.Value{Typ: "null"}
	}

	return resp.NewBulk(value)
}
