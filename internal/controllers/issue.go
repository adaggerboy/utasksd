package controllers

import (
	"fmt"
	"net/http"

	"github.com/adaggerboy/utasksd/models/context"
	"github.com/adaggerboy/utasksd/models/generic"
	"github.com/adaggerboy/utasksd/models/reqresp"
)

//Reusable

func EnsureUserIssueEditCapability(ctx context.IContext, issueID int) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	_, e, err := ctx.GetDatabase().EnsureIssueCapability(ctx, userID, issueID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database EnsureUserIssueEditCapability(): %s", err))
	} else if !e {
		ctx.AddPublicError(http.StatusForbidden, "Access denied")
	}
}

func EnsureUserIssueReadCapability(ctx context.IContext, issueID int) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	r, _, err := ctx.GetDatabase().EnsureIssueCapability(ctx, userID, issueID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database EnsureUserIssueReadCapability(): %s", err))
	} else if !r {
		ctx.AddPublicError(http.StatusForbidden, "Access denied")
	}
}

//Issue API Controllers

func GetIssueInfo(ctx context.IContext, issueID int, detailed bool) (issue *generic.Issue) {
	issue = nil
	var err error
	EnsureUserIssueReadCapability(ctx, issueID)
	if !ctx.IsActive() {
		return
	}
	if !detailed {
		issue, err = ctx.GetDatabase().ReadIssue(ctx, issueID)
	} else {
		issue, err = ctx.GetDatabase().ReadIssueRecursive(ctx, issueID)
	}
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database GetIssueInfo(): %s", err))
		issue = nil
		return
	} else if issue == nil {
		ctx.AddPublicError(http.StatusNotFound, "Task not found")
		return
	}
	return
}

func CreateIssue(ctx context.IContext, issue generic.Issue) (resultId int) {
	resultId = 0
	if issue.Project == nil {
		ctx.AddPublicError(http.StatusBadRequest, "Field \"project\" required")
		return
	}
	EnsureUserProjectEditIssuesCapability(ctx, *issue.Project)
	if !ctx.IsActive() {
		return
	}
	userID, _ := ctx.GetUserID()
	issue.Registrar = &userID

	id, err := ctx.GetDatabase().CreateIssue(ctx, issue)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database CreateIssue(): %s", err))
		return
	} else if id == nil {
		ctx.AddPublicError(http.StatusConflict, "Task can't be created, possibly already exists")
		return
	}
	return *id
}

func UpdateIssue(ctx context.IContext, issueID int, issue generic.Issue) {
	issue.Registrar = nil
	issue.Project = nil

	EnsureUserIssueEditCapability(ctx, issueID)
	if !ctx.IsActive() {
		return
	}

	err := ctx.GetDatabase().UpdateIssue(ctx, issueID, issue)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database UpdateIssue(): %s", err))
		return
	}
}

func DeleteIssue(ctx context.IContext, issueID int) {
	EnsureUserIssueEditCapability(ctx, issueID)
	if !ctx.IsActive() {
		return
	}
	err := ctx.GetDatabase().DeleteIssue(ctx, issueID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database DeleteIssue(): %s", err))
		return
	}
}

func SearchIssues(ctx context.IContext, query reqresp.SearchIssues) (issueIDs []int) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}

	issueIDs, err := ctx.GetDatabase().SearchIssues(ctx, query, userID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database SearchIssues(): %s", err))
		return
	}
	return
}
