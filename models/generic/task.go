package generic

import "time"

type TaskStatus string
type TaskPriority string
type TaskDependencyType string

const (
	ToDo       TaskStatus = "to-do"
	InProgress TaskStatus = "in-progress"
	Testing    TaskStatus = "testing"
	Waiting    TaskStatus = "waiting"
	Done       TaskStatus = "done"
)

const (
	Highest TaskPriority = "highest"
	High    TaskPriority = "high"
	Middle  TaskPriority = "middle"
	Low     TaskPriority = "low"
	Lowest  TaskPriority = "lowest"
)

const (
	BlockedBy TaskDependencyType = "blocked-by"
	Includes  TaskDependencyType = "includes"
)

type DependentTask struct {
	ID                 int    `json:"id"`
	TaskDependencyType string `json:"dependency_type"`
}

type Comment struct {
	ID       *int           `json:"id,omitempty"`
	TaskID   *int           `json:"task_id,omitempty"`
	AuthorID *int           `json:"author_id,omitempty"`
	TimeAgo  *time.Duration `json:"time_ago,omitempty"`
	Text     *string        `json:"text,omitempty"`
}

func AllocateComment() *Comment {
	return &Comment{
		ID:       new(int),
		TaskID:   new(int),
		AuthorID: new(int),
		TimeAgo:  new(time.Duration),
		Text:     new(string),
	}
}

type TrackingRecord struct {
	TaskID   *int            `json:"task_id,omitempty"`
	EndDate  *CommonDate     `json:"end_date,omitempty"`
	Duration *CommonDuration `json:"duration,omitempty"`
	Text     *string         `json:"text,omitempty"`
}

func AllocateTrackingRecord() *TrackingRecord {
	return &TrackingRecord{
		TaskID:   new(int),
		EndDate:  new(CommonDate),
		Duration: new(CommonDuration),
		Text:     new(string),
	}
}

type Task struct {
	// ID        int           `json:"id"`
	Project     *int          `json:"project,omitempty"`
	Name        *string       `json:"name,omitempty"`
	Description *string       `json:"description,omitempty"`
	Status      *TaskStatus   `json:"status,omitempty"`
	Priority    *TaskPriority `json:"priority,omitempty"`
	StartDate   *CommonDate   `json:"start_date,omitempty"`
	DueDate     *CommonDate   `json:"due_date,omitempty"`
	Assigner    *int          `json:"assigner,omitempty"`
	Assignees   *[]int        `json:"assignees,omitempty"`

	//Recursively acquired data
	LinkedIssues   *[]int            `json:"linked_issues,omitempty"`
	DependentTasks *[]DependentTask  `json:"dependent_tasks,omitempty"`
	Comments       *[]Comment        `json:"comments,omitempty"`
	TrackedRecords *[]TrackingRecord `json:"tracked_records,omitempty"`
	Attachments    *[]string         `json:"attachments,omitempty"`
}

func AllocateTask() *Task {
	return &Task{
		Project:     new(int),
		Name:        new(string),
		Description: new(string),
		Status:      new(TaskStatus),
		Priority:    new(TaskPriority),
		StartDate:   new(CommonDate),
		DueDate:     new(CommonDate),
		Assigner:    new(int),
		Assignees:   &[]int{},

		LinkedIssues:   &[]int{},
		DependentTasks: &[]DependentTask{},
		Comments:       &[]Comment{},
		TrackedRecords: &[]TrackingRecord{},
		Attachments:    &[]string{},
	}
}
