package message

import (
	"strconv"
	"strings"
)

type Response struct {
	responseCode    int
	responseMessage string
	responseBody    string
}

func NewResponse(responseCode int, responseMessage string, responseBody string) *Response  {
	return &Response{
		responseCode: responseCode,
		responseMessage: responseMessage,
		responseBody: responseBody,
	}
}

func (r *Response) SetResponseCode(responseCode int) {
	r.responseCode = responseCode
}

func (r *Response) SetResponseMessage(responseMessage string) {
	r.responseMessage = responseMessage
}

func (r *Response) SetResponseBody(responseBody string) {
	r.responseBody = responseBody
}

func (r *Response) ToString() string {
	if r.responseBody != "" && !strings.HasSuffix(r.responseBody, "\r\n") {
		r.responseBody += "\r\n"
	}

	return strconv.Itoa(r.responseCode) +
		" " +
		r.responseMessage +
		"\r\n" +
		r.responseBody
}
