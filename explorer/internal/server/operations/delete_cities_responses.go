// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// DeleteCitiesOKCode is the HTTP code returned for type DeleteCitiesOK
const DeleteCitiesOKCode int = 200

/*
DeleteCitiesOK OK

swagger:response deleteCitiesOK
*/
type DeleteCitiesOK struct {
}

// NewDeleteCitiesOK creates DeleteCitiesOK with default headers values
func NewDeleteCitiesOK() *DeleteCitiesOK {

	return &DeleteCitiesOK{}
}

// WriteResponse to the client
func (o *DeleteCitiesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// DeleteCitiesBadRequestCode is the HTTP code returned for type DeleteCitiesBadRequest
const DeleteCitiesBadRequestCode int = 400

/*
DeleteCitiesBadRequest Bad request

swagger:response deleteCitiesBadRequest
*/
type DeleteCitiesBadRequest struct {
}

// NewDeleteCitiesBadRequest creates DeleteCitiesBadRequest with default headers values
func NewDeleteCitiesBadRequest() *DeleteCitiesBadRequest {

	return &DeleteCitiesBadRequest{}
}

// WriteResponse to the client
func (o *DeleteCitiesBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// DeleteCitiesNotFoundCode is the HTTP code returned for type DeleteCitiesNotFound
const DeleteCitiesNotFoundCode int = 404

/*
DeleteCitiesNotFound Not found

swagger:response deleteCitiesNotFound
*/
type DeleteCitiesNotFound struct {
}

// NewDeleteCitiesNotFound creates DeleteCitiesNotFound with default headers values
func NewDeleteCitiesNotFound() *DeleteCitiesNotFound {

	return &DeleteCitiesNotFound{}
}

// WriteResponse to the client
func (o *DeleteCitiesNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// DeleteCitiesInternalServerErrorCode is the HTTP code returned for type DeleteCitiesInternalServerError
const DeleteCitiesInternalServerErrorCode int = 500

/*
DeleteCitiesInternalServerError Server error

swagger:response deleteCitiesInternalServerError
*/
type DeleteCitiesInternalServerError struct {
}

// NewDeleteCitiesInternalServerError creates DeleteCitiesInternalServerError with default headers values
func NewDeleteCitiesInternalServerError() *DeleteCitiesInternalServerError {

	return &DeleteCitiesInternalServerError{}
}

// WriteResponse to the client
func (o *DeleteCitiesInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
