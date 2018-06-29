package protocol

// issueIdIndex content is filled later, we can leave it empty here
var (
	DefaultDashboard = Dashboard{
		Columns: []Column{
			{
				Id:           "default-php",
				Language:     Language{"php"},
				IssueIdIndex: nil,
			},

			{
				Id:           "default-js",
				Language:     Language{"js"},
				IssueIdIndex: nil,
			},

			{
				Id:           "default-html",
				Language:     Language{"html"},
				IssueIdIndex: nil,
			},

			{
				Id:           "default-ruby",
				Language:     Language{"ruby"},
				IssueIdIndex: nil,
			},

			{
				Id:           "default-rust",
				Language:     Language{"rust"},
				IssueIdIndex: nil,
			},

			{
				Id:           "default-go",
				Language:     Language{"go"},
				IssueIdIndex: nil,
			},

			{
				Id:           "default-scala",
				Language:     Language{"scala"},
				IssueIdIndex: nil,
			},
		},
	}
)

type Dashboard struct {
	Columns []Column
}
