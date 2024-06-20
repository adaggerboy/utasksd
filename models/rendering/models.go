package rendering

import (
	"html/template"

	"github.com/adaggerboy/utasksd/models/generic"
)

type Option struct {
	Href  string
	Label string
}
type ProjectData struct {
	Assigners map[int]string
	Assignees map[int]string
	Supports  map[int]string
	ID        int
	InProject bool
}

type Project struct {
	Name        string
	ID          int
	Description string
	Logo        string
	Role        string
	Owner       Account
	Users       []Account
}

type Account struct {
	ID           int
	PathToAvatar string
	Name         string
	Email        string
	Role         string
}

type Core struct {
	Statuses            map[string]string
	IssueStatuses       map[string]string
	Priorities          map[string]string
	TaskDependencyTypes map[string]string

	CSSVariables    map[string]string
	JSVariablesJSON template.JS

	Account Account
	Options []Option
}

type Comment struct {
	Name      string
	TimeAgo   string
	Text      string
	ID        int
	Avatar    string
	CanRemove bool
}

type DependentTask struct {
	Name           string
	ID             int
	DependencyType string
}

type LinkedIssue struct {
	Name string
	ID   int
}

type TrackingRecord struct {
	Text     string
	Duration string
}

type Task struct {
	Name            string
	Description     string
	ID              int
	Status          string
	Priority        string
	Assigner        string
	AssignerID      int
	Assignee        string
	AssigneeID      int
	Attachments     []string
	StartDate       string
	DueDate         string
	DependentTasks  []DependentTask
	LinkedIssues    []LinkedIssue
	Comments        []Comment
	TrackingRecords []TrackingRecord
	OverallDuration string

	ChangeCapability bool
	EditCapability   bool
	ManageCapability bool

	Created bool
}

type TaskRoot struct {
	ProjectData ProjectData
	Core        Core
	Task        Task
}

type TreeNode struct {
	ID           int
	Name         string
	Children     *[]TreeNode
	VisibleClass string
}

type ListTask struct {
	FirePriority bool
	ProjectName  string
	ProjectID    int
	Name         string
	ID           int
	Status       string
	Description  string
	VisibleClass string
	Assigner     Account
	Assignees    []Account
}

type ListViewRoot struct {
	// JSVariablesJSON   template.JS
	// Account           Account
	ProjectData       ProjectData
	Core              Core
	Tasks             []ListTask
	AvailableProjects []TreeNode
}

type Issue struct {
	Name        string
	Description string
	ID          int
	Status      string
	Registrar   string
	Reporter    string
	Attachments []string
	StartDate   string
	DueDate     string

	EditCapability bool

	Created bool
	Closed  bool
}

type IssueRoot struct {
	ProjectData ProjectData
	Core        Core
	Issue       Issue
}

type ListIssue struct {
	ProjectName  string
	ProjectID    int
	Name         string
	ID           int
	Status       string
	Description  string
	VisibleClass string
	Registrar    Account
	Reporter     string
}

type ListIssueViewRoot struct {
	// JSVariablesJSON   template.JS
	// Account           Account
	ProjectData       ProjectData
	Core              Core
	Issues            []ListIssue
	AvailableProjects []TreeNode
}

type GenericUserWrapper struct {
	generic.User
	Admin, Director, Active bool
}

type ListUserViewRoot struct {
	Core    Core
	Users   []GenericUserWrapper
	IsAdmin bool
}

type ReportRoot struct {
	Tasks          []generic.TasksReport
	TimeEfficiency []generic.TimeEfficiencyReport
}

type ListProjectViewRoot struct {
	Core     Core
	Projects []Project
}

type ProjectRoot struct {
	Core    Core
	Project Project
	IsOwner bool
}

type UserRoot struct {
	User  generic.User
	ItsMe bool
	Core  Core
}
