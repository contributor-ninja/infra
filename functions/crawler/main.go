package main

import (
	"github.com/apex/log"
	loghandlers "github.com/apex/log/handlers/json"
	"github.com/aws/aws-lambda-go/lambda"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/contributor-ninja/infra/dynamodb"
	"github.com/contributor-ninja/infra/github"
	"github.com/contributor-ninja/infra/protocol"
)

const (
	FETCH_CONCURRENCY = 10
)

func init() {
	log.SetHandler(loghandlers.Default)
}

func main() {
	lambda.Start(crawler)
}

func crawler(in *protocol.GitHubProject) {
	githubClient := github.NewClient()
	svc, err := dynamodb.NewClient()

	if err != nil {
		log.WithError(err).Fatal("connection to dynamodb failed")
	}

	if in == nil {
		log.Fatal("could not decode input")
	}

	githubProject := *in

	query := dynamodb.MakeFindQuery(githubProject)
	resp, findQueryErr := svc.Query(&query)

	if findQueryErr != nil {
		log.WithError(findQueryErr).Fatal("could not find project")
	}

	if len(resp.Items) == 0 {
		log.WithField("project", githubProject).Fatal("could not find project")
	}

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
	}
}
