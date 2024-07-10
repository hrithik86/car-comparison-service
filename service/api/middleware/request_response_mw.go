package middleware

import (
	"car-comparison-service/errors"
	"car-comparison-service/serdes"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/go-playground/validator.v9/non-standard/validators"
	"io/ioutil"
	"net/http"
)

type ApiHandlerFunc[T any] func(ctx context.Context, r serdes.Request[T]) (serdes.Response, error)

func RequestResponseMw[T any](h ApiHandlerFunc[T]) http.Handler {
	validatorV9 = validator.New()
	_ = validatorV9.RegisterValidation("notblank", validators.NotBlank)
	return &requestResponseLoggerMwHandler[T]{next: h, isNilBody: false}
}

func NilRequestResponseMw(h ApiHandlerFunc[serdes.NilBody]) http.Handler {
	validatorV9 = validator.New()
	_ = validatorV9.RegisterValidation("notblank", validators.NotBlank)
	return &requestResponseLoggerMwHandler[serdes.NilBody]{next: h, isNilBody: true}
}

type requestResponseLoggerMwHandler[T any] struct {
	next      ApiHandlerFunc[T]
	isNilBody bool
}

func (h *requestResponseLoggerMwHandler[T]) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.serveHTTP(ctx, rw, r)
	if err != nil {
		errorHandler(ctx, rw, err)
	}
}

func (h *requestResponseLoggerMwHandler[T]) serveHTTP(ctx context.Context, rw http.ResponseWriter, r *http.Request) (err error) {
	var body T
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.FAILED_TO_READ_REQUEST_BODY
	}

	if !h.isNilBody {
		err = json.Unmarshal(reqBytes, &body)
		if err != nil {
			return errors.NewServiceError("REQUEST_DESERIALIZATION_ERROR: "+err.Error(), 400)
		}
		err = validatorV9.Struct(body)
		if err != nil {
			if _, ok := err.(validator.ValidationErrors); ok {
				return errors.NewServiceError("REQUEST_BODY_VALIDATION_FAILED: "+err.Error(), 400)
			}
		}
	}

	request := serdes.NewHttpRequest(mux.Vars(r), r.URL.Query(), body, r.URL.Path, r.Header)

	res, err := h.next(ctx, request)
	if err != nil {
		return err
	}
	var bytes []byte
	bytes, err = json.Marshal(res.Body())
	if err != nil {
		return errors.RESPONSE_SERIALIZATION_ERROR
	}
	writeResponse(ctx, rw, bytes, res.Status())
	return nil
}
