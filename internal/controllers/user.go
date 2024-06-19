package controllers

import (
	"fmt"
	"net/http"

	"github.com/adaggerboy/utasksd/models/context"
	"github.com/adaggerboy/utasksd/models/generic"
)

//Reusable

func EnsureAdminUsersManageCapability(ctx context.IContext) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	m, _, err := ctx.GetDatabase().EnsureUsersCapability(ctx, userID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database EnsureAdminUsersManageCapability(): %s", err))
	} else if !m {
		ctx.AddPublicError(http.StatusForbidden, "Access denied")
	}
}

func EnsureUsersDirectorCapability(ctx context.IContext) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	m, _, err := ctx.GetDatabase().EnsureUsersCapability(ctx, userID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database EnsureAdminUsersManageCapability(): %s", err))
	} else if !m && false {
		ctx.AddPublicError(http.StatusForbidden, "Access denied")
	}
}

//User API controllers

func GetUserInfo(ctx context.IContext, userID int) (user *generic.User) {
	_, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	user, err := ctx.GetDatabase().ReadUser(ctx, userID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database GetUserInfo(): %s", err))
		user = nil
		return
	} else if user == nil {
		ctx.AddPublicError(http.StatusNotFound, "User not found")
		return
	}
	return
}

func GetMyUserInfo(ctx context.IContext) (user *generic.User) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	user, err := ctx.GetDatabase().ReadUser(ctx, userID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database GetMyUserInfo(): %s", err))
		user = nil
		return
	} else if user == nil {
		ctx.AddPublicError(http.StatusNotFound, "User not found")
		return
	}
	user.ID = userID
	return
}

func UpdateMyUser(ctx context.IContext, user generic.User) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	err := ctx.GetDatabase().UpdateUser(ctx, userID, user)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database UpdateMyUser(): %s", err))
		return
	}
}

func DeleteMyUser(ctx context.IContext) {
	userID, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	err := ctx.GetDatabase().DeactivateUser(ctx, userID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database DeleteMyUser(): %s", err))
		return
	}
}

func UpdateUser(ctx context.IContext, userID int, user generic.User) {
	EnsureAdminUsersManageCapability(ctx)
	if !ctx.IsActive() {
		return
	}
	err := ctx.GetDatabase().UpdateUser(ctx, userID, user)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database UpdateUser(): %s", err))
		return
	}
}

func DeleteUser(ctx context.IContext, userID int) {
	EnsureAdminUsersManageCapability(ctx)
	if !ctx.IsActive() {
		return
	}
	err := ctx.GetDatabase().DeactivateUser(ctx, userID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("write database DeleteUser(): %s", err))
		return
	}
}

func GetUserTasks(ctx context.IContext, userID int) (taskIDs []int) {
	if !RequireStrictUser(ctx, userID) {
		return
	}
	taskIDs, err := ctx.GetDatabase().GetUserTasks(ctx, userID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database GetUserTasks(): %s", err))
		return
	}
	return
}

func GetUserProjects(ctx context.IContext, userID int) (projectIDs []int) {
	if !RequireStrictUser(ctx, userID) {
		return
	}
	projectIDs, err := ctx.GetDatabase().GetUserProjects(ctx, userID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database GetUserProjects(): %s", err))
		return
	}
	return
}

func GetUserProjectTasks(ctx context.IContext, userID int, projectID int) (taskIDs []int) {
	if !RequireStrictUser(ctx, userID) {
		return
	}
	taskIDs, err := ctx.GetDatabase().GetUserProjectTasks(ctx, userID, projectID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database GetUserProjectTasks(): %s", err))
		return
	}
	return
}

func GetUserIssues(ctx context.IContext, userID int) (issueIDs []int) {
	if !RequireStrictUser(ctx, userID) {
		return
	}
	issueIDs, err := ctx.GetDatabase().GetUserIssues(ctx, userID)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database GetUserIssues(): %s", err))
		return
	}
	return
}

func GetAllUsers(ctx context.IContext) (users []generic.User) {
	_, ok := RequireAuth(ctx)
	if !ok {
		return
	}
	users, err := ctx.GetDatabase().GetAllUsers(ctx)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database GetUserTasks(): %s", err))
		return
	}
	return
}
