package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/apex/log"
	loghandlers "github.com/apex/log/handlers/json"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"

	"github.com/contributor-ninja/infra/api"
	"github.com/contributor-ninja/infra/dynamodb"
	"github.com/contributor-ninja/infra/protocol"
)

var (
	port = os.Getenv("PORT")
)

func init() {
	log.SetHandler(loghandlers.Default)
}

func main() {
	addr := ":" + port

	r := mux.NewRouter()

	r.HandleFunc("/status", getStatusHandler)

	r.
		HandleFunc("/repo/{org}/{name}", putRepoHandler).
		Methods("PUT")

	if err := http.ListenAndServe(addr, r); err != nil {
		log.WithError(err).Fatal("error listening")
	}
}

func getStatusHandler(w http.ResponseWriter, r *http.Request) {
	svc, err := dynamodb.NewClient()

	if err != nil {
		log.WithError(err).Fatal("connection to dynamodb failed")
	}

	_, listTablesErr := svc.ListTables(&awsdynamodb.ListTablesInput{})

	if listTablesErr != nil {
		log.WithError(listTablesErr).Fatal("connection to dynamodb failed")
	}

	res := api.Response{"ok"}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		log.WithError(err).Fatal("could not send response")
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(jsonBytes); err != nil {
		log.WithError(err).Fatal("could not send response")
	}
}

func putRepoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["org"] == "" || vars["name"] == "" {
		log.Fatal("malformed request")
	}

	svc, err := dynamodb.NewClient()

	if err != nil {
		log.WithError(err).Fatal("connection to dynamodb failed")
	}

	/*
	   Check if project already exists
	*/
	predicat := protocol.GitHubProject{
		Org:  vars["org"],
		Name: vars["name"],
	}

	query := dynamodb.MakeFindQuery(predicat)
	resp, findQueryErr := svc.Query(&query)

	if findQueryErr != nil {
		log.WithError(findQueryErr).Fatal("could not find project")
	}

	if len(resp.Items) > 0 {
		// Project already exists
		api.SendResponse(api.Response{"project already exists"}, w)
		return
	}

	/*
		Insert new project
	*/
	project := protocol.MakeGitHubProject(vars["org"], vars["name"])

	av, err := dynamodbattribute.MarshalMap(project)

	input := &awsdynamodb.PutItemInput{
		Item:      av,
		TableName: protocol.GitHubProjectTable,
	}

	_, putErr := svc.PutItem(input)

	if putErr != nil {
		log.WithError(putErr).Fatal("could not send response")
	}

	log.
		WithFields(log.Fields{
			"id": project.Id,
		}).
		Info("added item in index")

	api.SendResponse(api.Response{"added " + project.Id}, w)
}
