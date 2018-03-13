package graphql

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/apex/log"
	graphql "github.com/graph-gophers/graphql-go"
)

func GetSchema(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

type GraphQL struct {
	Schema *graphql.Schema
}

func (h *GraphQL) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	ctx := r.Context()

	response := h.Schema.Exec(ctx, params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.WithError(err).Info("response error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
