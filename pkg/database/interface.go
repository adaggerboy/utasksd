package database

import (
	"context"

	"github.com/adaggerboy/utasksd/models/generic"
	"github.com/adaggerboy/utasksd/models/reqresp"
)

type IDatabase interface {
	CreateTask(ctx context.Context, task generic.Task) (taskID *int, err error)
	UpdateTask(ctx context.Context, taskID int, task generic.Task) (err error)
	DeleteTask(ctx context.Context, taskID int) (err error)
	ReadTask(ctx context.Context, taskID int) (task *generic.Task, err error)
	ReadTaskRecursive(ctx context.Context, taskID int) (task *generic.Task, err error)

	SearchTasks(ctx context.Context, query reqresp.SearchTasks, userID int) (taskIDs []int, err error)

	ReadComment(ctx context.Context, commentID int) (comment *generic.Comment, err error)
	PublishComment(ctx context.Context, userID int, taskID int, text string) (commentID *int, err error)
	DeleteComment(ctx context.Context, commentID int) (err error)

	CreateUser(ctx context.Context, user generic.User) (userID *int, err error)
	UpdateUser(ctx context.Context, userID int, user generic.User) (err error)
	DeactivateUser(ctx context.Context, userID int) (err error)
	ReadUser(ctx context.Context, userID int) (user *generic.User, err error)
	GetUserByUsername(ctx context.Context, username string) (userID *int, isActive bool, err error)

	GetUserProjects(ctx context.Context, userID int) (projectIDs []int, err error)

	GetUserTasks(ctx context.Context, userID int) (taskIDs []int, err error)
	GetUserProjectTasks(ctx context.Context, userID, projectID int) (taskIDs []int, err error)
	GetProjectTasks(ctx context.Context, projectID int) (taskIDs []int, err error)

	GetUserIssues(ctx context.Context, userID int) (issueIDs []int, err error)
	GetProjectIssues(ctx context.Context, projectID int) (issueIDs []int, err error)

	CreateProject(ctx context.Context, project generic.Project) (projectID *int, err error)
	UpdateProject(ctx context.Context, projectID int, project generic.Project) (err error)
	DeleteProject(ctx context.Context, projectID int) (err error)
	ReadProject(ctx context.Context, projectID int) (project *generic.Project, err error)
	ReadProjectRecursive(ctx context.Context, projectID int) (project *generic.Project, err error)

	EnsureTaskCapability(ctx context.Context, userID, taskID int) (readCapability, editCapability, manageCapability bool, err error)
	EnsureIssueCapability(ctx context.Context, userID, issueID int) (readCapability, editCapability bool, err error)
	EnsureUsersCapability(ctx context.Context, userID int) (manageCapability, directorCapability bool, err error)
	EnsureProjectCapability(ctx context.Context, userID, projectID int) (readCapability, tasksCapability, issuesCapability, manageCapability bool, err error)

	RegisterAttachment(ctx context.Context, attachment string) (err error)
	TrackTaskActivity(ctx context.Context, taskID int, text string, duration int) (err error)

	CreateIssue(ctx context.Context, issue generic.Issue) (issueID *int, err error)
	UpdateIssue(ctx context.Context, issueID int, issue generic.Issue) (err error)
	DeleteIssue(ctx context.Context, issueID int) (err error)
	ReadIssue(ctx context.Context, issueID int) (issue *generic.Issue, err error)
	ReadIssueRecursive(ctx context.Context, issueID int) (issue *generic.Issue, err error)

	SearchIssues(ctx context.Context, query reqresp.SearchIssues, userID int) (issueID []int, err error)

	GetAllUsers(ctx context.Context) (users []generic.User, err error)

	SetUserPermissions(ctx context.Context, userID int, username string, isActive, isAdmin, isDirector bool) (err error)
	Close() (err error)
	RollbackClose() (err error)
	DeploySchema() (err error)
	CreateDBUser(ctx context.Context, username string, password string) (err error)
	DeleteDBUser(ctx context.Context, username string) (err error)
	ChangeDBUserPassword(ctx context.Context, username string, password string) (err error)

	TasksReport(ctx context.Context, request reqresp.ReportRequest) (perUser []generic.TasksReport, err error)
	TimeEfficiencyReport(ctx context.Context, request reqresp.ReportRequest) (perUser []generic.TimeEfficiencyReport, err error)
}
