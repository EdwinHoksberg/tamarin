package message

import (
	"strings"
)

type Request struct {
	command    string
	parameters []string
}

func NewRequest(command string, parameters []string) Request {
	return Request{
		command:    command,
		parameters: parameters,
	}
}

func (r *Request) GetCommand() string {
	return r.command
}

func (r *Request) GetParameters() []string {
	return r.parameters
}

func (r *Request) ToString() string {
	request := r.command

	if len(r.parameters) > 0 {
		request += strings.Join(r.parameters, " ")
	}

	return request
}
