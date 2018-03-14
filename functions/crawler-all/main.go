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
	defaultphpIssueIdIndex  = protocol.MakeEmptyIssueIdIndex("default-php")
	defaultjsIssueIdIndex   = protocol.MakeEmptyIssueIdIndex("default-js")
	defaulthtmlIssueIdIndex = protocol.MakeEmptyIssueIdIndex("default-html")
	defaultrubyIssueIdIndex = protocol.MakeEmptyIssueIdIndex("default-ruby")
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

					log.WithField("title", issue.Title).Info("insered issue")

					if githubProject.Language == "js" {
						defaultjsIssueIdIndex.Ids = append(defaultjsIssueIdIndex.Ids, issue.Id)
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
	inputQuery := dynamodb.MakePutItemIssueIdIndex(defaultphpIssueIdIndex)
	_, putErr := svc.PutItem(&inputQuery)

	if putErr != nil {
		log.
			WithField("name", "default-php").
			WithError(putErr).
			Info("put IndexIssueId failed")
	}

	inputQuery = dynamodb.MakePutItemIssueIdIndex(defaultjsIssueIdIndex)
	_, putErr = svc.PutItem(&inputQuery)

	if putErr != nil {
		log.
			WithField("name", "default-js").
			WithError(putErr).
			Info("put IndexIssueId failed")
	}

	inputQuery = dynamodb.MakePutItemIssueIdIndex(defaulthtmlIssueIdIndex)
	_, putErr = svc.PutItem(&inputQuery)

	if putErr != nil {
		log.
			WithField("name", "default-html").
			WithError(putErr).
			Info("put IndexIssueId failed")
	}

	inputQuery = dynamodb.MakePutItemIssueIdIndex(defaultrubyIssueIdIndex)
	_, putErr = svc.PutItem(&inputQuery)

	if putErr != nil {
		log.
			WithField("name", "default-ruby").
			WithError(putErr).
			Info("put IndexIssueId failed")
	}
}
