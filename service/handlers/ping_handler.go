package handlers

import (
	"encoding/json"
	"net/http"
)

type PingResponse struct {
	Status string `json:"status"`
}

func PingHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		payload, _ := json.Marshal(PingResponse{
			Status: "healthy",
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
	}
}
