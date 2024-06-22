package controllers

import (
	"fmt"
	"net/http"

	"github.com/adaggerboy/utasksd/models/context"
	"github.com/adaggerboy/utasksd/models/generic"
)

//Reusable

func GetUserProjectCapabilities(ctx context.IContext, userID, projectID int) (read, tasks, issues, manage bool) {
	read, tasks, issues, manage, err := ctx.GetDatabase().EnsureProjectCapability(ctx, userID, projectID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database EnsureUserProjectManageCapability(): %s", err))
	}
	return
}

func EnsureUserProjectManageCapability(ctx context.IContext, projectID int) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	_, _, _, m, err := ctx.GetDatabase().EnsureProjectCapability(ctx, userID, projectID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database EnsureUserProjectManageCapability(): %s", err))
	} else if !m {
		ctx.AddPublicError(http.StatusForbidden, "Access denied")
	}
}

func EnsureUserProjectEditTasksCapability(ctx context.IContext, projectID int) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	_, t, _, _, err := ctx.GetDatabase().EnsureProjectCapability(ctx, userID, projectID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database EnsureUserProjectEditCapability(): %s", err))
	} else if !t {
		ctx.AddPublicError(http.StatusForbidden, "Access denied")
	}
}

func EnsureUserProjectEditIssuesCapability(ctx context.IContext, projectID int) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	_, _, i, _, err := ctx.GetDatabase().EnsureProjectCapability(ctx, userID, projectID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database EnsureUserProjectEditCapability(): %s", err))
	} else if !i {
		ctx.AddPublicError(http.StatusForbidden, "Access denied")
	}
}

func EnsureUserProjectReadCapability(ctx context.IContext, projectID int) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	r, _, _, _, err := ctx.GetDatabase().EnsureProjectCapability(ctx, userID, projectID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database EnsureUserProjectReadCapability(): %s", err))
	} else if !r {
		ctx.AddPublicError(http.StatusForbidden, "Access denied")
	}
}

//Project API Controllers

func GetProjectInfo(ctx context.IContext, projectID int, detailed bool) (project *generic.Project) {
	var err error
	project = nil
	EnsureUserProjectReadCapability(ctx, projectID)
	if !ctx.IsActive() {
		return
	}
	if !detailed {
		project, err = ctx.GetDatabase().ReadProject(ctx, projectID)
	} else {
		project, err = ctx.GetDatabase().ReadProjectRecursive(ctx, projectID)
	}
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database GetProjectInfo(): %s", err))
		project = nil
		return
	} else if project == nil {
		ctx.AddPublicError(http.StatusNotFound, "Project not found")
		return
	}
	return
}

func CreateProject(ctx context.IContext, project generic.Project) (resultId int) {
	if project.LogoPath == nil {
		project.LogoPath = new(string)
		*project.LogoPath = "null"
	} else if (*project.LogoPath)[0] == '/' {
		*project.LogoPath = (*project.LogoPath)[1:]
	}
	EnsureUsersDirectorCapability(ctx)
	if !ctx.IsActive() {
		return
	}
	userID, _ := ctx.GetUserID()
	project.OwnerID = &userID
	id, err := ctx.GetDatabase().CreateProject(ctx, project)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database CreateProject(): %s", err))
		return
	} else if id == nil {
		ctx.AddPublicError(http.StatusConflict, "Project can't be created")
		return
	}
	return *id
}

func UpdateProject(ctx context.IContext, projectID int, project generic.Project) {
	EnsureUserProjectManageCapability(ctx, projectID)
	if !ctx.IsActive() {
		return
	}
	if (*project.LogoPath)[0] == '/' {
		*project.LogoPath = (*project.LogoPath)[1:]
	}
	err := ctx.GetDatabase().UpdateProject(ctx, projectID, project)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database UpdateProject(): %s", err))
		return
	}
}

func DeleteProject(ctx context.IContext, projectID int) {
	EnsureUserProjectManageCapability(ctx, projectID)
	if !ctx.IsActive() {
		return
	}
	err := ctx.GetDatabase().DeleteProject(ctx, projectID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database DeleteProject(): %s", err))
		return
	}
}

func GetProjectTasks(ctx context.IContext, projectID int) (taskIDs []int) {
	EnsureUserProjectReadCapability(ctx, projectID)
	if !ctx.IsActive() {
		return
	}
	taskIDs, err := ctx.GetDatabase().GetProjectTasks(ctx, projectID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database GetProjectTasks(): %s", err))
		return
	}
	return
}

func GetProjectIssues(ctx context.IContext, projectID int) (issueIDs []int) {
	EnsureUserProjectReadCapability(ctx, projectID)
	if !ctx.IsActive() {
		return
	}
	issueIDs, err := ctx.GetDatabase().GetProjectIssues(ctx, projectID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database GetProjectIssues(): %s", err))
		return
	}
	return
}
