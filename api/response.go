package api

import (
	"encoding/json"
	"net/http"

	"github.com/apex/log"
)

type Response struct {
	Status string `json:"status"`
}

func SendResponse(res Response, w http.ResponseWriter) {
	jsonBytes, err := json.Marshal(res)

	if err != nil {
		log.WithError(err).Fatal("could not send response")
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(jsonBytes); err != nil {
		log.WithError(err).Fatal("could not send response")
	}
}
