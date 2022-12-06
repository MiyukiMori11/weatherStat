// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// PostCitiesHandlerFunc turns a function with the right signature into a post cities handler
type PostCitiesHandlerFunc func(PostCitiesParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostCitiesHandlerFunc) Handle(params PostCitiesParams) middleware.Responder {
	return fn(params)
}

// PostCitiesHandler interface for that can handle valid post cities params
type PostCitiesHandler interface {
	Handle(PostCitiesParams) middleware.Responder
}

// NewPostCities creates a new http.Handler for the post cities operation
func NewPostCities(ctx *middleware.Context, handler PostCitiesHandler) *PostCities {
	return &PostCities{Context: ctx, Handler: handler}
}

/*
	PostCities swagger:route POST /cities postCities

Adds a new city into subscripton list
*/
type PostCities struct {
	Context *middleware.Context
	Handler PostCitiesHandler
}

func (o *PostCities) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostCitiesParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
