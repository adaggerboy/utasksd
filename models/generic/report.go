package generic

type ReportBase struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type TasksReport struct {
	ReportBase
	Tasks        int `json:"tasks"`
	DoneTasks    int `json:"done_tasks"`
	Issues       int `json:"issues"`
	ClosedIssues int `json:"closed_issues"`
}
type TimeEfficiencyReport struct {
	ReportBase
	TrackedRecords  int     `json:"tasks"`
	SummaryHours    float64 `json:"summary"`
	AveragePerDay   float64 `json:"average_per_day"`
	AveragePerMonth float64 `json:"average_per_months"`
}
