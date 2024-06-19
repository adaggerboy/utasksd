package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/adaggerboy/utasksd/models/context"
	"github.com/adaggerboy/utasksd/models/generic"
	"github.com/adaggerboy/utasksd/models/reqresp"
)

//Reusable

func EnsureUserTaskManageCapability(ctx context.IContext, taskID int) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	_, _, m, err := ctx.GetDatabase().EnsureTaskCapability(ctx, userID, taskID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database EnsureUserTaskManageCapability(): %s", err))
	} else if !m {
		ctx.AddPublicError(http.StatusForbidden, "Access denied")
	}
}

func EnsureUserTaskEditCapability(ctx context.IContext, taskID int) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	_, e, _, err := ctx.GetDatabase().EnsureTaskCapability(ctx, userID, taskID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database EnsureUserTaskEditCapability(): %s", err))
	} else if !e {
		ctx.AddPublicError(http.StatusForbidden, "Access denied")
	}
}

func EnsureUserTaskReadCapability(ctx context.IContext, taskID int) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	r, _, _, err := ctx.GetDatabase().EnsureTaskCapability(ctx, userID, taskID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database EnsureUserTaskReadCapability(): %s", err))
	} else if !r {
		ctx.AddPublicError(http.StatusForbidden, "Access denied")
	}
}

//Task API Controllers

func GetTaskInfo(ctx context.IContext, taskID int, detailed bool) (task *generic.Task) {
	task = nil
	var err error
	EnsureUserTaskReadCapability(ctx, taskID)
	if !ctx.IsActive() {
		return
	}
	if !detailed {
		task, err = ctx.GetDatabase().ReadTask(ctx, taskID)
	} else {
		task, err = ctx.GetDatabase().ReadTaskRecursive(ctx, taskID)
	}
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database GetTaskInfo() detailed:%t: %s", detailed, err))
		task = nil
		return
	} else if task == nil {
		ctx.AddPublicError(http.StatusNotFound, "Task not found")
		return
	}
	return
}

func CreateTask(ctx context.IContext, task generic.Task) (resultId int) {
	resultId = 0
	if task.Project == nil {
		ctx.AddPublicError(http.StatusBadRequest, "Field \"project\" required")
		return
	}
	EnsureUserProjectEditTasksCapability(ctx, *task.Project)
	if !ctx.IsActive() {
		return
	}
	userID, _ := ctx.GetUserID()
	task.Assigner = &userID

	id, err := ctx.GetDatabase().CreateTask(ctx, task)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database CreateTask(): %s", err))
		return
	} else if id == nil {
		ctx.AddPublicError(http.StatusConflict, "Task can't be created, possibly already exists")
		return
	}
	return *id
}

func UpdateTask(ctx context.IContext, taskID int, task generic.Task) {
	task.Assigner = nil
	task.Project = nil

	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	_, e, m, err := ctx.GetDatabase().EnsureTaskCapability(ctx, userID, taskID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database UpdateTask(): %s", err))
		return
	}
	if !e {
		ctx.AddPublicError(http.StatusForbidden, "Access denied")
		return
	}
	if !m {
		task.Assignees = nil
		task.Name = nil
		task.Priority = nil
		task.StartDate = nil
		task.DueDate = nil
		task.LinkedIssues = nil
		task.DependentTasks = nil
		task.Attachments = nil
	}

	err = ctx.GetDatabase().UpdateTask(ctx, taskID, task)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database UpdateTask(): %s", err))
		return
	}
}

func DeleteTask(ctx context.IContext, taskID int) {
	EnsureUserTaskManageCapability(ctx, taskID)
	if !ctx.IsActive() {
		return
	}
	err := ctx.GetDatabase().DeleteTask(ctx, taskID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database DeleteTask(): %s", err))
		return
	}
}

func SearchTasks(ctx context.IContext, query reqresp.SearchTasks) (taskIDs []int) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}

	taskIDs, err := ctx.GetDatabase().SearchTasks(ctx, query, userID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database SearchTasks(): %s", err))
		return
	}
	return
}

func ReadComment(ctx context.IContext, commentID int) (comment *generic.Comment) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	comment, err := ctx.GetDatabase().ReadComment(ctx, commentID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database ReadComment(): %s", err))
		return
	} else if comment == nil {
		ctx.AddPublicError(http.StatusNotFound, "Comment not found")
		return
	}
	if userID != *comment.AuthorID {
		EnsureUserTaskReadCapability(ctx, *comment.TaskID)
		if !ctx.IsActive() {
			return
		}
	}
	return
}

func PublishComment(ctx context.IContext, taskID int, text string) (resultId int) {
	resultId = 0
	EnsureUserTaskReadCapability(ctx, taskID)
	if !ctx.IsActive() {
		return
	}
	userID, _ := ctx.GetUserID()
	id, err := ctx.GetDatabase().PublishComment(ctx, userID, taskID, text)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database PublishComment(): %s", err))
		return
	} else if id == nil {
		ctx.AddPublicError(http.StatusConflict, "Comment can't be created, idk why")
		return
	}
	return *id
}

func DeleteComment(ctx context.IContext, commentID int) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	comment, err := ctx.GetDatabase().ReadComment(ctx, commentID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database DeleteComment(): %s", err))
		return
	} else if comment == nil {
		ctx.AddPublicError(http.StatusNotFound, "Comment not found")
		return
	}
	if userID != *comment.AuthorID {
		EnsureUserTaskManageCapability(ctx, *comment.TaskID)
		if !ctx.IsActive() {
			return
		}
	}
	err = ctx.GetDatabase().DeleteComment(ctx, commentID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database DeleteComment(): %s", err))
		return
	}
}

func EnsureUserIsAssingeeOfTask(ctx context.IContext, taskID int) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	task := GetTaskInfo(ctx, taskID, true)
	if !ctx.IsActive() {
		return
	}
	found := false
	for _, v := range *task.Assignees {
		if v == userID {
			found = true
		}
	}
	if !found {
		ctx.AddPublicError(http.StatusNotAcceptable, "User is not an assigner of task")
	}
}

func TrackTaskActivity(ctx context.IContext, taskID int, text string, duration time.Duration) {

	err := ctx.GetDatabase().TrackTaskActivity(ctx, taskID, text, int(duration.Seconds()))
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database TrackTaskActivity(): %s", err))
		return
	}
}
