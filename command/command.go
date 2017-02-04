package command

import (
	"github.com/edwinhoksberg/tamarin/message"
)

type Command struct {
	request message.Request
}

func New(request message.Request) *Command {
	return &Command{
		request: request,
	}
}

func (c *Command) GenerateResponse() message.Response {
	switch c.request.GetCommand() {
	case "quit":
		return new(quitCommand).generate(c.request)
	case "help":
		return new(helpCommand).generate(c.request)
	case "capabilities":
		return new(capabilitiesCommand).generate(c.request)
	case "date":
		return new(dateCommand).generate(c.request)
	default:
		return *message.NewResponse(500, "command '"+c.request.GetCommand()+"' not recognized", "")
	}
}
