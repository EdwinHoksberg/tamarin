package command

import (
	"github.com/edwinhoksberg/tamarin/message"
	"time"
	"fmt"
)

type dateCommand struct{}

func (c *dateCommand) generate(request message.Request) message.Response {
	now := time.Now()
	dateTime := fmt.Sprintf("%04d%02d%02d%02d%02d%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())

	return *message.NewResponse(111, dateTime, "")
}
