package middleware

import (
	"car-comparison-service/errors"
	"context"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"net/http"
)

func errorHandler(ctx context.Context, rw http.ResponseWriter, err error) {
	serviceError := errors.UNKNOWN
	errors2.As(err, &serviceError)

	bytes, err := json.Marshal(serviceError)
	if err != nil {
		fmt.Printf("Failed to serialize Error Response" + err.Error())
	}
	writeResponse(ctx, rw, bytes, serviceError.Status)
}

func writeResponse(ctx context.Context, rw http.ResponseWriter, bytes []byte, statusCode int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	_, err := rw.Write(bytes)
	if err != nil {
		fmt.Printf("Failed to write response" + err.Error())
	}
}
