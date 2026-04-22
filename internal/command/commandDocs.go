package command

import "rodis/internal/protocol/resp"

type CommandDocs struct{}

func (c *CommandDocs) Execute(args []resp.Payload, ctx *CommandContext) resp.Payload {
	return resp.NewBulk("OK")
}
