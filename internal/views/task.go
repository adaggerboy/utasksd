package views

import (
	"fmt"
	"time"

	"github.com/adaggerboy/utasksd/internal/controllers"
	"github.com/adaggerboy/utasksd/models/context"
	"github.com/adaggerboy/utasksd/models/generic"
	"github.com/adaggerboy/utasksd/models/rendering"
)

func formatDuration(d time.Duration) string {
	seconds := int(d.Seconds())
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}

	minutes := seconds / 60
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}

	hours := minutes / 60
	if hours < 24 {
		return fmt.Sprintf("%dh", hours)
	}

	days := hours / 24
	return fmt.Sprintf("%dd", days)
}

func RenderTask(ctx context.IContext, taskID int, task *generic.Task, projectOwner int) rendering.Task {

	userID, ok := ctx.GetUserID()
	if !ok {
		return rendering.Task{}
	}

	assigner := controllers.GetUserInfo(ctx, *task.Assigner)
	if !ctx.IsActive() {
		return rendering.Task{}
	}

	result := rendering.Task{
		Name:            *task.Name,
		Description:     *task.Description,
		ID:              taskID,
		Attachments:     []string{},
		StartDate:       (*task.StartDate).Format("2006-01-02"),
		DueDate:         (*task.DueDate).Format("2006-01-02"),
		DependentTasks:  []rendering.DependentTask{},
		Status:          string(*task.Status),
		Priority:        string(*task.Priority),
		Comments:        []rendering.Comment{},
		TrackingRecords: []rendering.TrackingRecord{},
		LinkedIssues:    []rendering.LinkedIssue{},
		Assigner:        *assigner.Firstname + " " + *assigner.Lastname,
		AssignerID:      assigner.ID,
		Assignee:        "",
		AssigneeID:      0,
		Created:         true,

		ManageCapability: userID == projectOwner,
		EditCapability:   userID == projectOwner || userID == *task.Assigner,
		ChangeCapability: userID == projectOwner || userID == *task.Assigner,
	}

	if len(*task.Assignees) > 0 {
		assignee := controllers.GetUserInfo(ctx, (*task.Assignees)[0])
		if !ctx.IsActive() {
			return rendering.Task{}
		}
		result.Assignee = *assignee.Firstname + " " + *assignee.Lastname
		result.AssigneeID = assignee.ID
	}

	if !result.ChangeCapability {
		for _, v := range *task.Assignees {
			if v == userID {
				result.ChangeCapability = true
				break
			}
		}
	}

	if task.Attachments != nil {
		result.Attachments = *task.Attachments
	}

	for _, v := range *task.DependentTasks {
		depTask := controllers.GetTaskInfo(ctx, v.ID, false)
		if !ctx.IsActive() {
			return rendering.Task{}
		}
		result.DependentTasks = append(result.DependentTasks, rendering.DependentTask{
			Name:           *depTask.Name,
			ID:             v.ID,
			DependencyType: v.TaskDependencyType,
		})
	}

	for _, v := range *task.LinkedIssues {
		depTask := controllers.GetIssueInfo(ctx, v, false)
		if !ctx.IsActive() {
			return rendering.Task{}
		}
		result.LinkedIssues = append(result.LinkedIssues, rendering.LinkedIssue{
			Name: *depTask.Name,
			ID:   v,
		})
	}

	overall := time.Duration(0)

	for _, v := range *task.TrackedRecords {
		overall += v.Duration.Duration
		result.TrackingRecords = append([]rendering.TrackingRecord{{
			Duration: formatDuration(v.Duration.Duration),
			Text:     *v.Text,
		}}, result.TrackingRecords...)
	}

	result.OverallDuration = overall.String()

	for _, v := range *task.Comments {
		author := controllers.GetUserInfo(ctx, *v.AuthorID)
		if !ctx.IsActive() {
			return rendering.Task{}
		}
		result.Comments = append([]rendering.Comment{{
			ID:        *v.ID,
			CanRemove: userID == projectOwner || userID == *v.AuthorID,
			Name:      *author.Firstname + " " + *author.Lastname,
			TimeAgo:   formatDuration(*v.TimeAgo),
			Text:      *v.Text,
			Avatar:    *author.AvatarPath,
		}}, result.Comments...)
	}

	return result
}

//UI Controllers

func RenderTaskView(ctx context.IContext, taskID int) (data []byte) {

	core := GetRendererCoreStruct()
	RenderCoreInfo(ctx, &core)
	if !ctx.IsActive() {
		return
	}

	task := controllers.GetTaskInfo(ctx, taskID, true)
	if !ctx.IsActive() {
		return
	}

	project := controllers.GetProjectInfo(ctx, *task.Project, true)
	if !ctx.IsActive() {
		return
	}

	renderTask := RenderTask(ctx, taskID, task, *project.OwnerID)
	if !ctx.IsActive() {
		return
	}

	input := rendering.TaskRoot{
		Core:        core,
		ProjectData: RenderProjectData(ctx, *task.Project, project),
		Task:        renderTask,
	}
	if !ctx.IsActive() {
		return
	}

	return RenderAs(ctx, "task", input)
}

func RenderCreateTaskView(ctx context.IContext, projectID int) (data []byte) {

	core := GetRendererCoreStruct()
	user := RenderCoreInfo(ctx, &core)
	if !ctx.IsActive() {
		return
	}

	project := controllers.GetProjectInfo(ctx, projectID, true)
	if !ctx.IsActive() {
		return
	}

	input := rendering.TaskRoot{
		Core:        core,
		ProjectData: RenderProjectData(ctx, projectID, project),
		Task: rendering.Task{
			Name:           "New task",
			Description:    "Click on task name or on this text to edit...",
			ID:             0,
			Attachments:    []string{},
			StartDate:      time.Now().Format("2006-01-02"),
			DueDate:        time.Now().Format("2006-01-02"),
			DependentTasks: []rendering.DependentTask{},
			Status:         string(generic.ToDo),
			Priority:       string(generic.Middle),
			Comments:       []rendering.Comment{},
			LinkedIssues:   []rendering.LinkedIssue{},
			Assigner:       *user.Firstname + " " + *user.Lastname,
			AssignerID:     user.ID,
			Assignee:       "",
			AssigneeID:     0,
			Created:        false,

			ManageCapability: user.ID == *project.OwnerID,
			EditCapability:   true,
			ChangeCapability: true,
		},
	}
	if !ctx.IsActive() {
		return
	}

	return RenderAs(ctx, "task", input)
}
