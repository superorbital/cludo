// Code generated by go-swagger; DO NOT EDIT.

package environment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/superorbital/cludo/models"
)

// GenerateEnvironmentOKCode is the HTTP code returned for type GenerateEnvironmentOK
const GenerateEnvironmentOKCode int = 200

/*GenerateEnvironmentOK OK

swagger:response generateEnvironmentOK
*/
type GenerateEnvironmentOK struct {

	/*
	  In: Body
	*/
	Payload *models.ModelsEnvironmentResponse `json:"body,omitempty"`
}

// NewGenerateEnvironmentOK creates GenerateEnvironmentOK with default headers values
func NewGenerateEnvironmentOK() *GenerateEnvironmentOK {

	return &GenerateEnvironmentOK{}
}

// WithPayload adds the payload to the generate environment o k response
func (o *GenerateEnvironmentOK) WithPayload(payload *models.ModelsEnvironmentResponse) *GenerateEnvironmentOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the generate environment o k response
func (o *GenerateEnvironmentOK) SetPayload(payload *models.ModelsEnvironmentResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GenerateEnvironmentOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GenerateEnvironmentBadRequestCode is the HTTP code returned for type GenerateEnvironmentBadRequest
const GenerateEnvironmentBadRequestCode int = 400

/*GenerateEnvironmentBadRequest Bad Request

swagger:response generateEnvironmentBadRequest
*/
type GenerateEnvironmentBadRequest struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGenerateEnvironmentBadRequest creates GenerateEnvironmentBadRequest with default headers values
func NewGenerateEnvironmentBadRequest() *GenerateEnvironmentBadRequest {

	return &GenerateEnvironmentBadRequest{}
}

// WithPayload adds the payload to the generate environment bad request response
func (o *GenerateEnvironmentBadRequest) WithPayload(payload string) *GenerateEnvironmentBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the generate environment bad request response
func (o *GenerateEnvironmentBadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GenerateEnvironmentBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*GenerateEnvironmentDefault generic error response

swagger:response generateEnvironmentDefault
*/
type GenerateEnvironmentDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGenerateEnvironmentDefault creates GenerateEnvironmentDefault with default headers values
func NewGenerateEnvironmentDefault(code int) *GenerateEnvironmentDefault {
	if code <= 0 {
		code = 500
	}

	return &GenerateEnvironmentDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the generate environment default response
func (o *GenerateEnvironmentDefault) WithStatusCode(code int) *GenerateEnvironmentDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the generate environment default response
func (o *GenerateEnvironmentDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the generate environment default response
func (o *GenerateEnvironmentDefault) WithPayload(payload *models.Error) *GenerateEnvironmentDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the generate environment default response
func (o *GenerateEnvironmentDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GenerateEnvironmentDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
