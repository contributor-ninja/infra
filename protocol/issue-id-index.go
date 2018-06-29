package protocol

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

var (
	IssueIdIndexTable = aws.String("IssueIndices")
)

type IssueIdIndex struct {
	Name        string `json:"name"` // primary key
	Ids         []int  `json:"ids"`
	LastUpdated string `json:"lastupdated"`
}

func MakeEmptyIssueIdIndex(name string) IssueIdIndex {
	LastUpdated := time.Now().Format("2006-01-02 15:04:05")

	return IssueIdIndex{
		Name:        name,
		Ids:         make([]int, 0),
		LastUpdated: LastUpdated,
	}
}
