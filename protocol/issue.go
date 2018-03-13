package protocol

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

var (
	IssueTable = aws.String("NewIssues")
)

type User struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

type Issue struct {
	Id        int       `json:"id"` //  primary key
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	HTMLURL   string    `json:"html_url"`
	AvatarURL string    `json:"avatarURL"`

	Project GitHubProject `json:"project"`

	User User `json:"user"`
}
