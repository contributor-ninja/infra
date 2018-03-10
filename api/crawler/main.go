package main

import (
	"net/http"
	"os"
	"sync"

	"github.com/apex/log"
	loghandlers "github.com/apex/log/handlers/json"
	"github.com/aws/aws-sdk-go/aws"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"

	"github.com/contributor-ninja/infra/api"
	"github.com/contributor-ninja/infra/dynamodb"
	"github.com/contributor-ninja/infra/github"
	"github.com/contributor-ninja/infra/protocol"
)

var (
	port = os.Getenv("PORT")
)

const (
	FETCH_CONCURRENCY = 10
)

func init() {
	log.SetHandler(loghandlers.Default)
}

func main() {
	addr := ":" + port

	r := mux.NewRouter()

	r.
		HandleFunc("/crawl/all", postCrawlAllHandler).
		Methods("POST")

	if err := http.ListenAndServe(addr, r); err != nil {
		log.WithError(err).Fatal("error listening")
	}
}

func postCrawlAllHandler(w http.ResponseWriter, r *http.Request) {
	githubClient := github.NewClient()
	svc, err := dynamodb.NewClient()

	if err != nil {
		log.WithError(err).Fatal("connection to dynamodb failed")
	}

	/*
	  Scan projects
	*/
	scanParams := &awsdynamodb.ScanInput{
		TableName: protocol.GitHubProjectTable,
		AttributesToGet: []*string{
			aws.String("name"),
			aws.String("org"),
			aws.String("labels"),
		},
	}

	resp, scanErr := svc.Scan(scanParams)

	if scanErr != nil {
		log.WithError(scanErr).Fatal("could not get projects")
	}

	var items = resp.Items

	for {
		if len(items) == 0 {
			break
		}

		var projectFetchBatch []map[string]*awsdynamodb.AttributeValue

		if len(items) < FETCH_CONCURRENCY {
			projectFetchBatch = items
			items = items[len(items):] // empty array
		} else {
			projectFetchBatch = items[:FETCH_CONCURRENCY]
			items = items[FETCH_CONCURRENCY:]
		}

		var wg sync.WaitGroup
		wg.Add(len(projectFetchBatch))

		for _, item := range projectFetchBatch {

			go func(item map[string]*awsdynamodb.AttributeValue) {
				defer wg.Done()

				currentItem := item

				/*
					Fetch issues
				*/

				githubProject := protocol.MakeGitHubProject(*currentItem["org"].S, *currentItem["name"].S)

				log.Info("processing project " + githubProject.GetId())

				issues, fetchErr := githubClient.FetchIssues(githubProject)

				if fetchErr != nil {
					log.WithError(fetchErr).Info("could not fetch issues for " + githubProject.GetId())
				}

				/*
					Insert issues
				*/

				for _, issue := range issues {
					av, err := dynamodbattribute.MarshalMap(issue)

					if err != nil {
						log.WithError(err).Info("could not add in index")
					}

					input := &awsdynamodb.PutItemInput{
						Item:      av,
						TableName: protocol.IssueTable,
					}

					_, putErr := svc.PutItem(input)

					if putErr != nil {
						log.WithError(putErr).Info("could not add in index")
					}
				}

			}(item)
		}

		wg.Wait()
	}

	log.Info("processed all projects")

	api.SendResponse(api.Response{"ok"}, w)
}
