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
	CityName string
	/*
	  Required: true
	  Max Length: 25
	  Pattern: ^[\w\-\' ]+
	  In: query
	*/
	CountryName string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetTempParams() beforehand.
func (o *GetTempParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qCityName, qhkCityName, _ := qs.GetOK("city_name")
	if err := o.bindCityName(qCityName, qhkCityName, route.Formats); err != nil {
		res = append(res, err)
	}

	qCountryName, qhkCountryName, _ := qs.GetOK("country_name")
	if err := o.bindCountryName(qCountryName, qhkCountryName, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindCityName binds and validates parameter CityName from query.
func (o *GetTempParams) bindCityName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("city_name", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("city_name", "query", raw); err != nil {
		return err
	}
	o.CityName = raw

	if err := o.validateCityName(formats); err != nil {
		return err
	}

	return nil
}

// validateCityName carries on validations for parameter CityName
func (o *GetTempParams) validateCityName(formats strfmt.Registry) error {

	if err := validate.MaxLength("city_name", "query", o.CityName, 25); err != nil {
		return err
	}

	if err := validate.Pattern("city_name", "query", o.CityName, `^[\w\-\' ]+`); err != nil {
		return err
	}

	return nil
}

// bindCountryName binds and validates parameter CountryName from query.
func (o *GetTempParams) bindCountryName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("country_name", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("country_name", "query", raw); err != nil {
		return err
	}
	o.CountryName = raw

	if err := o.validateCountryName(formats); err != nil {
		return err
	}

	return nil
}

// validateCountryName carries on validations for parameter CountryName
func (o *GetTempParams) validateCountryName(formats strfmt.Registry) error {

	if err := validate.MaxLength("country_name", "query", o.CountryName, 25); err != nil {
		return err
	}

	if err := validate.Pattern("country_name", "query", o.CountryName, `^[\w\-\' ]+`); err != nil {
		return err
	}

	return nil
}
