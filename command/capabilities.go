package command

import (
	"fmt"
	"github.com/edwinhoksberg/tamarin/message"
	"strings"
)

type capabilitiesCommand struct{}

func (c *capabilitiesCommand) generate(request message.Request) message.Response {
	capabilities := []string{
		"VERSION 2",
		fmt.Sprintf("IMPLEMENTATION %s %s", "test", "test"),
		".",
	}

	return *message.NewResponse(101, "Capability list follows", strings.Join(capabilities, "\r\n"))
}
