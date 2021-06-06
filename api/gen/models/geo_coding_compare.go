// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// GeoCodingCompare geo coding compare
//
// swagger:model GeoCodingCompare
type GeoCodingCompare struct {

	// env company Id
	EnvCompanyID string `json:"envCompanyId,omitempty"`

	// request
	Request map[string]Address `json:"request,omitempty"`

	// response
	Response map[string]Success `json:"response,omitempty"`
}

// Validate validates this geo coding compare
func (m *GeoCodingCompare) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRequest(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResponse(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GeoCodingCompare) validateRequest(formats strfmt.Registry) error {
	if swag.IsZero(m.Request) { // not required
		return nil
	}

	for k := range m.Request {

		if err := validate.Required("request"+"."+k, "body", m.Request[k]); err != nil {
			return err
		}
		if val, ok := m.Request[k]; ok {
			if err := val.Validate(formats); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *GeoCodingCompare) validateResponse(formats strfmt.Registry) error {
	if swag.IsZero(m.Response) { // not required
		return nil
	}

	for k := range m.Response {

		if err := validate.Required("response"+"."+k, "body", m.Response[k]); err != nil {
			return err
		}
		if val, ok := m.Response[k]; ok {
			if err := val.Validate(formats); err != nil {
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this geo coding compare based on the context it is used
func (m *GeoCodingCompare) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateRequest(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateResponse(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GeoCodingCompare) contextValidateRequest(ctx context.Context, formats strfmt.Registry) error {

	for k := range m.Request {

		if val, ok := m.Request[k]; ok {
			if err := val.ContextValidate(ctx, formats); err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *GeoCodingCompare) contextValidateResponse(ctx context.Context, formats strfmt.Registry) error {

	for k := range m.Response {

		if val, ok := m.Response[k]; ok {
			if err := val.ContextValidate(ctx, formats); err != nil {
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *GeoCodingCompare) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GeoCodingCompare) UnmarshalBinary(b []byte) error {
	var res GeoCodingCompare
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
