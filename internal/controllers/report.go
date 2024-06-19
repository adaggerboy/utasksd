package controllers

import (
	"fmt"
	"net/http"

	"github.com/adaggerboy/utasksd/models/context"
	"github.com/adaggerboy/utasksd/models/generic"
	"github.com/adaggerboy/utasksd/models/reqresp"
)

func GetTasksReport(ctx context.IContext, request reqresp.ReportRequest) (report []generic.TasksReport) {
	if request.Project != nil {
		EnsureUserProjectManageCapability(ctx, *request.Project)
	} else {
		EnsureAdminUsersManageCapability(ctx)
	}
	if !ctx.IsActive() {
		return
	}
	report, err := ctx.GetDatabase().TasksReport(ctx, request)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database GetTasksReport(): %s", err))
		return
	}
	return
}
func GetTimeEfficiencyReport(ctx context.IContext, request reqresp.ReportRequest) (report []generic.TimeEfficiencyReport) {
	if request.Project != nil {
		EnsureUserProjectManageCapability(ctx, *request.Project)
	} else {
		EnsureAdminUsersManageCapability(ctx)
	}
	if !ctx.IsActive() {
		return
	}
	report, err := ctx.GetDatabase().TimeEfficiencyReport(ctx, request)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("read database GetTimeEfficiencyReport(): %s", err))
		return
	}
	return
}
