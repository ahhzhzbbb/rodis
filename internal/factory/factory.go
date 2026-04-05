package factory

import "rodis/internal/command"

type Factory interface {
	CreateCommand() command.Command
}

var CommandRegistry = map[string]func() command.Command{
	"PING":    func() command.Command { return &command.PingCommand{} },
	"COMMAND": func() command.Command { return &command.CommandDocs{} },
}
