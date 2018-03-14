package protocol

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
		},
	}
)

type Dashboard struct {
	Columns []Column
}
