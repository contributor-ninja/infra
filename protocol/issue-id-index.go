package protocol

import "github.com/aws/aws-sdk-go/aws"

var (
	IssueIdIndexTable = aws.String("issueids")
)

type IssueIdIndex struct {
	Name string `json:"name"` // primary key
	Ids  []int  `json:"ids"`
}
