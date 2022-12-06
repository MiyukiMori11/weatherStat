// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewGetTempParams creates a new GetTempParams object
//
// There are no default values defined in the spec.
func NewGetTempParams() GetTempParams {

	return GetTempParams{}
}

// GetTempParams contains all the bound params for the get temp operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetTemp
type GetTempParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  Max Length: 25
	  Pattern: ^[\w\-\' ]+
	  In: query
	*/
	Name string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetTempParams() beforehand.
func (o *GetTempParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qName, qhkName, _ := qs.GetOK("name")
	if err := o.bindName(qName, qhkName, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindName binds and validates parameter Name from query.
func (o *GetTempParams) bindName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("name", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("name", "query", raw); err != nil {
		return err
	}
	o.Name = raw

	if err := o.validateName(formats); err != nil {
		return err
	}

	return nil
}

// validateName carries on validations for parameter Name
func (o *GetTempParams) validateName(formats strfmt.Registry) error {

	if err := validate.MaxLength("name", "query", o.Name, 25); err != nil {
		return err
	}

	if err := validate.Pattern("name", "query", o.Name, `^[\w\-\' ]+`); err != nil {
		return err
	}

	return nil
}
