package generic

type IssueStatus string

const (
	Open     IssueStatus = "open"
	Closed   IssueStatus = "closed"
	Reopened IssueStatus = "reopened"
)

type Issue struct {
	// ID        int           `json:"id"`
	Project     *int         `json:"project,omitempty"`
	Name        *string      `json:"name,omitempty"`
	Description *string      `json:"description,omitempty"`
	Status      *IssueStatus `json:"status,omitempty"`
	StartDate   *CommonDate  `json:"start_date,omitempty"`
	CloseDate   *CommonDate  `json:"close_date,omitempty"`
	Registrar   *int         `json:"registrar,omitempty"`
	Reporter    *string      `json:"reporter,omitempty"`

	//Recursively acquired data
	Attachments *[]string `json:"attachments,omitempty"`
}

func AllocateIssue() *Issue {
	return &Issue{
		Project:     new(int),
		Name:        new(string),
		Reporter:    new(string),
		Description: new(string),
		Status:      new(IssueStatus),
		StartDate:   new(CommonDate),
		CloseDate:   new(CommonDate),
		Registrar:   new(int),

		Attachments: &[]string{},
	}
}
