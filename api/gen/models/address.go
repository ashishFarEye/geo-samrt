// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Address address
//
// swagger:model Address
type Address struct {

	// country code
	CountryCode string `json:"countryCode,omitempty"`

	// landmark
	Landmark string `json:"landmark,omitempty"`

	// line1
	Line1 string `json:"line1,omitempty"`

	// line2
	Line2 string `json:"line2,omitempty"`

	// pincode
	Pincode string `json:"pincode,omitempty"`

	// uuid
	UUID string `json:"uuid,omitempty"`
}

// Validate validates this address
func (m *Address) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this address based on context it is used
func (m *Address) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Address) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Address) UnmarshalBinary(b []byte) error {
	var res Address
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
