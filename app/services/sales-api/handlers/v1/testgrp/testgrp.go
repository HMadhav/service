package testgrp

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Handlers struct {
	Build string
	Log   *zap.SugaredLogger
}

func (h Handlers) Test(w http.ResponseWriter, r *http.Request) {

	status := "ok"
	statusCode := http.StatusOK

	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	json.NewEncoder(w).Encode(data)

	h.Log.Infow("Test", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
}
