// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ModelsEnvironmentRequest models environment request
//
// swagger:model models.EnvironmentRequest
type ModelsEnvironmentRequest struct {

	// The requested target for the request
	Target string `json:"target,omitempty"`
}

// Validate validates this models environment request
func (m *ModelsEnvironmentRequest) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this models environment request based on context it is used
func (m *ModelsEnvironmentRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ModelsEnvironmentRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModelsEnvironmentRequest) UnmarshalBinary(b []byte) error {
	var res ModelsEnvironmentRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
