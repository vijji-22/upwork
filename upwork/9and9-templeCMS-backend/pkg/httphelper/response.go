package httphelper

import (
	"net/http"

	"encoding/json"
)

type errorWithStatusCode struct {
	err        error
	msg        string
	statusCode int
}

func NewErrorWithStatusCode(err error, msg string, statusCode int) error {
	return &errorWithStatusCode{err: err, msg: msg, statusCode: statusCode}
}

func (e *errorWithStatusCode) GetMessage() string {
	return e.msg
}

func (e *errorWithStatusCode) GetStatusCode() int {
	return e.statusCode
}

func (e *errorWithStatusCode) Error() string {
	return e.err.Error()
}

type ResponseMessageGetter interface {
	GetMessage() string
}

type StatusCodeGetter interface {
	GetStatusCode() int
}

// Response is response structure which will be sent to the client
type Response struct {
	Success  bool        `json:"success"`
	Message  string      `json:"message,omitempty"`
	Data     interface{} `json:"data"`
	Errors   []string    `json:"error,omitempty"`
	Warnings []string    `json:"warning,omitempty"`
}

// NewResponse is used to create new response
func NewResponse() *Response {
	return &Response{}
}

// Sucessfull will set success flag as true in response
func (r *Response) Sucessfull() *Response {
	r.Success = true
	return r
}

// Failed will set success flag as false in response
func (r *Response) Failed() *Response {
	r.Success = false
	return r
}

// SetMessage is used to set message in response
func (r *Response) SetMessage(message string) *Response {
	r.Message = message
	return r
}

// SetData is used to set data in response
func (r *Response) SetData(data interface{}) *Response {
	r.Data = data
	return r
}

// AddError is used to add error in response
func (r *Response) AddError(err error) *Response {
	r.Errors = append(r.Errors, err.Error())
	return r
}

// AddWarning is used to add warning in response
func (r *Response) AddWarning(warning string) *Response {
	r.Warnings = append(r.Warnings, warning)
	return r
}

// Send is used to send response to client
func (r *Response) Send(statusCode int, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(r)
}

// SendError is used to send error response to client. this will find status code from error and send response
func SendError(message string, err error, w http.ResponseWriter) error {
	resp := NewResponse().Failed()
	statusCode := http.StatusInternalServerError

	if m, ok := err.(ResponseMessageGetter); ok {
		resp.SetMessage(m.GetMessage())
	}
	if s, ok := err.(StatusCodeGetter); ok {
		statusCode = s.GetStatusCode()
	}

	return resp.Failed().AddError(err).SetMessage(message).Send(statusCode, w)
}

// SendData is used to send data response to client
func SendData(data interface{}, w http.ResponseWriter) error {
	return NewResponse().Sucessfull().SetData(data).Send(http.StatusOK, w)
}

// Send404 is used to send 404 response to client
func Send404(w http.ResponseWriter) error {
	return NewResponse().Failed().SetMessage("Not found").Send(http.StatusNotFound, w)
}
