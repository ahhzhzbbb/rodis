package command

import "rodis/internal/protocol/resp"

type CommandDocs struct{}

func (c *CommandDocs) Execute(args []resp.Value, ctx *CommandContext) resp.Value {
	return resp.NewBulk("OK")
}
