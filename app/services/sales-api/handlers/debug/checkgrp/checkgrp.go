package checkgrp

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Handlers struct {
	Build string
	Log   *zap.SugaredLogger
}

func (h Handlers) Readiness(w http.ResponseWriter, r *http.Request) {

	status := "ok"
	statusCode := http.StatusOK

	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	if err := response(w, statusCode, data); err != nil {
		h.Log.Errorw("readiness", "ERROR", err)
	}

	h.Log.Infow("readiness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
}

func (h Handlers) Liveness(w http.ResponseWriter, r *http.Request) {

	data := struct {
		Status string `json:"status,omitempty"`
		Build  string `json:"build,omitempty"`
	}{
		Status: "up",
		Build:  h.Build,
	}

	statusCode := http.StatusOK
	if err := response(w, statusCode, data); err != nil {
		h.Log.Errorw("liveness", "ERROR", err)
	}

	h.Log.Infow("liveness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
}

func response(w http.ResponseWriter, statusCode int, data interface{}) error {

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
