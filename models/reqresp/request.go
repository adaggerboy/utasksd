package reqresp

import "github.com/adaggerboy/utasksd/models/generic"

type UpdatePasswordRequest struct {
	Secret string `json:"secret"`
}

type CreateUserWithPasswordRequest struct {
	User   generic.User `json:"user"`
	Secret string       `json:"secret"`
}

type LoginWithPasswordRequest struct {
	Username string `json:"username"`
	Secret   string `json:"secret"`
}

type SearchTasks struct {
	Project   *int                  `json:"project,omitempty"`
	Name      *string               `json:"name,omitempty"`
	Status    *generic.TaskStatus   `json:"status,omitempty"`
	Priority  *generic.TaskPriority `json:"priority,omitempty"`
	StartDate *generic.CommonDate   `json:"start_date,omitempty"`
	DueDate   *generic.CommonDate   `json:"due_date,omitempty"`
	Assigner  *int                  `json:"assigner,omitempty"`
	Assignee  *int                  `json:"assignee,omitempty"`
	Issue     *int                  `json:"issue,omitempty"`
}

type SearchIssues struct {
	Project   *int                 `json:"project,omitempty"`
	Name      *string              `json:"name,omitempty"`
	Status    *generic.IssueStatus `json:"status,omitempty"`
	StartDate *generic.CommonDate  `json:"start_date,omitempty"`
	DueDate   *generic.CommonDate  `json:"due_date,omitempty"`
	Registrar *int                 `json:"registrar,omitempty"`
}

type ReportRequest struct {
	Project   *int                `json:"project,omitempty"`
	StartDate *generic.CommonDate `json:"start_date,omitempty"`
	DueDate   *generic.CommonDate `json:"due_date,omitempty"`
}
