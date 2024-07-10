package serdes

import (
	"net/http"
	"net/url"
)

type NilBody struct {
}

type Request[T any] interface {
	Param(key string) string
	Query(key string) string
	Body() T
	QueryParams() url.Values
	Header() http.Header
	Path() string
}

type HttpRequest[T any] struct {
	params  map[string]string
	queries url.Values
	body    T
	path    string
	header  http.Header
}

func (request HttpRequest[T]) Header() http.Header {
	return request.header
}

func (request HttpRequest[T]) Path() string {
	return request.path
}

func (request HttpRequest[T]) QueryParams() url.Values {
	return request.queries
}

func (request HttpRequest[T]) Param(key string) string {
	return request.params[key]
}

func (request HttpRequest[T]) Query(key string) string {
	return request.queries.Get(key)
}

func (request HttpRequest[T]) Body() T {
	return request.body
}

func NewHttpRequest[T any](params map[string]string, queries url.Values, body T, path string, header http.Header) Request[T] {
	return &HttpRequest[T]{params: params, queries: queries, body: body, path: path, header: header}
}

type RequestBodyFunc func() interface{}

type Response interface {
	Status() int
	Body() any
}

type HttpResponse struct {
	status int
	body   any
}

func (response *HttpResponse) Status() int {
	return response.status
}

func (response *HttpResponse) Body() any {
	return response.body
}

func NewHttpResponse(status int, body any) *HttpResponse {
	return &HttpResponse{status: status, body: body}
}
