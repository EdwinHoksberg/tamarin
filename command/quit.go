package command

import (
	"github.com/edwinhoksberg/tamarin/message"
)

type quitCommand struct{}

func (c *quitCommand) generate(request message.Request) message.Response {
	return *message.NewResponse(205, "Connection closing", "")
}
