package main

import (
	"sync"

	"github.com/apex/log"
	loghandlers "github.com/apex/log/handlers/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/contributor-ninja/infra/dynamodb"
	"github.com/contributor-ninja/infra/github"
	"github.com/contributor-ninja/infra/protocol"
)

const (
	FETCH_CONCURRENCY = 10
)

var (
	indices = map[string]*protocol.IssueIdIndex{
		"php":      protocol.NewEmptyIssueIdIndex("default-php"),
		"js":       protocol.NewEmptyIssueIdIndex("default-js"),
		"html-css": protocol.NewEmptyIssueIdIndex("default-html"),
		"ruby":     protocol.NewEmptyIssueIdIndex("default-ruby"),
		"rust":     protocol.NewEmptyIssueIdIndex("default-rust"),
		"go":       protocol.NewEmptyIssueIdIndex("default-go"),
		"scala":    protocol.NewEmptyIssueIdIndex("default-scala"),
	}
)

func init() {
	log.SetHandler(loghandlers.Default)
}

func main() {
	lambda.Start(crawlerAll)
}

func crawlerAll() {
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
			aws.String("language"),
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
				githubProject.Language = *currentItem["language"].S

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

					if index, ok := indices[githubProject.Language]; ok {
						if !index.HasIssueId(issue.Id) {
							index.Ids = append(index.Ids, issue.Id)
						}
					} else {
						log.Error("no index for language " + githubProject.Language)
					}
				}

			}(item)
		}

		wg.Wait()
	}

	log.Info("processed all projects")

	/*
	  Update IssueIndices
	*/
	for _, index := range indices {
		inputQuery := dynamodb.MakePutItemIssueIdIndex(index)
		_, putErr := svc.PutItem(&inputQuery)

		if putErr != nil {
			log.
				WithField("entries", len(index.Ids)).
				WithField("name", index.Name).
				WithError(putErr).
				Info("put IndexIssueId failed")
		}

	}
}
