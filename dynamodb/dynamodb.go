package dynamodb

import (
	"os"
	"strconv"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/contributor-ninja/infra/protocol"
)

var (
	dynamodbRegion = "us-east-1"

	id     = os.Getenv("DB_AWS_ACCESS_KEY_ID")
	secret = os.Getenv("DB_AWS_ACCESS_KEY")
	token  = ""

	maxIssuesPerCol = 50
)

func NewClient() (*awsdynamodb.DynamoDB, error) {
	creds := credentials.NewStaticCredentials(id, secret, token)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(dynamodbRegion),
		Credentials: creds,
	})

	if err != nil {
		return nil, err
	}

	return awsdynamodb.New(sess), err
}

func MakeFindQuery(p protocol.GitHubProject) awsdynamodb.QueryInput {
	return awsdynamodb.QueryInput{
		TableName: protocol.GitHubProjectTable,

		KeyConditions: map[string]*awsdynamodb.Condition{
			"id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*awsdynamodb.AttributeValue{
					{
						S: aws.String(p.GetId()),
					},
				},
			},
		},
	}
}

func MakeBatchGetItemByIssueIndex(p protocol.IssueIdIndex) awsdynamodb.BatchGetItemInput {
	keysArray := make([]map[string]*awsdynamodb.AttributeValue, 0)

	for _, id := range p.Ids {

		keys := map[string]*awsdynamodb.AttributeValue{
			"id": {
				N: aws.String(strconv.Itoa(id)),
			},
		}

		keysArray = append(keysArray, keys)

		if len(keysArray) > maxIssuesPerCol {
			log.Info("more than " + strconv.Itoa(maxIssuesPerCol) + " elements in the batch, truncating")
			break
		}
	}

	issueTableName := *protocol.IssueTable

	requestItems := make(map[string]*awsdynamodb.KeysAndAttributes)

	requestItems[issueTableName] = &awsdynamodb.KeysAndAttributes{
		ProjectionExpression: aws.String("id,title,body,html_url,gh_user"),
		Keys:                 keysArray,
	}

	return awsdynamodb.BatchGetItemInput{
		RequestItems: requestItems,
	}
}

func MakeGetIssueIdIndexQuery(name string) awsdynamodb.QueryInput {
	return awsdynamodb.QueryInput{
		TableName: protocol.IssueIdIndexTable,

		KeyConditions: map[string]*awsdynamodb.Condition{
			"name": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*awsdynamodb.AttributeValue{
					{
						S: aws.String(name),
					},
				},
			},
		},
	}
}

func MakePutItemIssueIdIndex(p *protocol.IssueIdIndex) awsdynamodb.PutItemInput {

	idNS := make([]*string, 0)

	for _, id := range p.Ids {
		idNS = append(idNS, aws.String(strconv.Itoa(id)))
	}

	return awsdynamodb.PutItemInput{
		TableName: protocol.IssueIdIndexTable,

		Item: map[string]*awsdynamodb.AttributeValue{
			"name": &awsdynamodb.AttributeValue{
				S: aws.String(p.Name),
			},
			"lastupdated": &awsdynamodb.AttributeValue{
				S: aws.String(p.LastUpdated),
			},
			"ids": &awsdynamodb.AttributeValue{
				NS: idNS,
			},
		},
	}
}

func MakePutGitHubProject(p protocol.GitHubProject) awsdynamodb.PutItemInput {
	labels := make([]*string, 0)

	for _, l := range p.Labels {
		labels = append(labels, aws.String(l))
	}

	return awsdynamodb.PutItemInput{
		TableName: protocol.GitHubProjectTable,

		Item: map[string]*awsdynamodb.AttributeValue{
			"id": &awsdynamodb.AttributeValue{
				S: aws.String(p.Id),
			},
			"labels": &awsdynamodb.AttributeValue{
				SS: labels,
			},
			"language": &awsdynamodb.AttributeValue{
				S: aws.String(p.Language),
			},
			"name": &awsdynamodb.AttributeValue{
				S: aws.String(p.Name),
			},
			"org": &awsdynamodb.AttributeValue{
				S: aws.String(p.Org),
			},
		},
	}
}
