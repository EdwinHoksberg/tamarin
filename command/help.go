package command

import (
	"github.com/edwinhoksberg/tamarin/message"
	"strings"
)

type helpCommand struct{}

func (c *helpCommand) generate(request message.Request) message.Response {
	//commandMap := getCommandMap()

	commands := make([]string, 0)
	/*for key := range commandMap {
		commands = append(commands, key)
	}*/
	commands = append(commands, ".")

	return *message.NewResponse(100, "Help text follows", strings.Join(commands, "\r\n"))
}
