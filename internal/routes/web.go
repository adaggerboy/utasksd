package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/adaggerboy/utasksd/internal/views"
	"github.com/adaggerboy/utasksd/models/generic"
	"github.com/adaggerboy/utasksd/models/reqresp"
	ginhelpers "github.com/adaggerboy/utasksd/pkg/utils/ginHelpers"
	"github.com/gin-gonic/gin"
)

type StructResponseJSON struct {
	Rate float32 `json:"rate"`
}

type StructSubscriptionRequest struct {
	Email string `json:"email"`
}

func GinPrintError(ctx *gin.Context, code int, message string) {
	ctx.String(code, message)
}

func DeployWebRoutes(r gin.IRouter) {

	r.GET("task/:id", ginhelpers.WrapContextWebID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, taskID int) (code int, cType string, data []byte) {
			return http.StatusOK, "text/html", views.RenderTaskView(nctx, taskID)
		},
	))

	r.GET("issue/:id", ginhelpers.WrapContextWebID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, taskID int) (code int, cType string, data []byte) {
			return http.StatusOK, "text/html", views.RenderIssueView(nctx, taskID)
		},
	))

	r.GET("login", ginhelpers.WrapContextWebPermitted(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper) (code int, cType string, data []byte) {
			return http.StatusOK, "text/html", views.RenderLoginView(nctx)
		},
	))

	r.GET("register", ginhelpers.WrapContextWebPermitted(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper) (code int, cType string, data []byte) {
			return http.StatusOK, "text/html", views.RenderRegisterView(nctx)
		},
	))

	r.GET("tasks", ginhelpers.WrapContextWeb(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper) (code int, cType string, data []byte) {
			return http.StatusOK, "text/html", views.RenderUserTaskListView(nctx)
		},
	))

	r.GET("projects", ginhelpers.WrapContextWeb(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper) (code int, cType string, data []byte) {
			return http.StatusOK, "text/html", views.RenderProjectListView(nctx)
		},
	))

	r.GET("issues", ginhelpers.WrapContextWeb(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper) (code int, cType string, data []byte) {
			return http.StatusOK, "text/html", views.RenderUserIssueListView(nctx)
		},
	))

	r.GET("users", ginhelpers.WrapContextWeb(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper) (code int, cType string, data []byte) {
			return http.StatusOK, "text/html", views.RenderUserListView(nctx)
		},
	))

	r.GET("project/:id", ginhelpers.WrapContextWebID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, projectID int) (code int, cType string, data []byte) {
			return http.StatusOK, "text/html", views.RenderProjectView(nctx, projectID)
		},
	))

	r.GET("user/:id", ginhelpers.WrapContextWebID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, userID int) (code int, cType string, data []byte) {
			return http.StatusOK, "text/html", views.RenderUserView(nctx, userID)
		},
	))

	r.GET("project/:id/create_report", ginhelpers.WrapContextWebID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, projectID int) (code int, cType string, data []byte) {
			return http.StatusOK, "text/html", views.RenderCreateReportView(nctx, projectID)
		},
	))

	r.GET("project/:id/tasks", ginhelpers.WrapContextWebID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, projectID int) (code int, cType string, data []byte) {
			return http.StatusOK, "text/html", views.RenderProjectTaskListView(nctx, projectID)
		},
	))

	r.GET("project/:id/issues", ginhelpers.WrapContextWebID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, projectID int) (code int, cType string, data []byte) {
			return http.StatusOK, "text/html", views.RenderProjectIssueListView(nctx, projectID)
		},
	))

	r.GET("project/:id/create_task", ginhelpers.WrapContextWebID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, projectID int) (code int, cType string, data []byte) {
			return http.StatusOK, "text/html", views.RenderCreateTaskView(nctx, projectID)
		},
	))

	r.GET("project/:id/create_issue", ginhelpers.WrapContextWebID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, projectID int) (code int, cType string, data []byte) {
			return http.StatusOK, "text/html", views.RenderCreateIssueView(nctx, projectID)
		},
	))

	r.GET("create_project", ginhelpers.WrapContextWeb(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper) (code int, cType string, data []byte) {
			return http.StatusOK, "text/html", views.RenderCreateProjectView(nctx)
		},
	))

	r.GET("search_tasks/list", ginhelpers.WrapContextWeb(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper) (code int, cType string, data []byte) {
			var searchParams reqresp.SearchTasks

			if val, err := strconv.Atoi(ctx.Query("project")); err == nil {
				searchParams.Project = &val
			}
			if val := ctx.Query("name"); val != "" {
				searchParams.Name = &val
			}
			if val := ctx.Query("status"); val != "" {
				status := generic.TaskStatus(val)
				searchParams.Status = &status
			}
			if val := ctx.Query("priority"); val != "" {
				priority := generic.TaskPriority(val)
				searchParams.Priority = &priority
			}
			if val, err := strconv.Atoi(ctx.Query("assigner")); err == nil {
				searchParams.Assigner = &val
			}
			if val, err := strconv.Atoi(ctx.Query("assignee")); err == nil {
				searchParams.Assignee = &val
			}
			return http.StatusOK, "text/html", views.RenderSearchProjectTaskListView(nctx, searchParams)
		},
	))

	r.GET("search_issues", ginhelpers.WrapContextWeb(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper) (code int, cType string, data []byte) {
			var searchParams reqresp.SearchIssues

			if val, err := strconv.Atoi(ctx.Query("project")); err == nil {
				searchParams.Project = &val
			}
			if val := ctx.Query("name"); val != "" {
				searchParams.Name = &val
			}
			if val := ctx.Query("status"); val != "" {
				status := generic.IssueStatus(val)
				searchParams.Status = &status
			}
			if val, err := strconv.Atoi(ctx.Query("registrar")); err == nil {
				searchParams.Registrar = &val
			}
			return http.StatusOK, "text/html", views.RenderSearchProjectIssueListView(nctx, searchParams)
		},
	))

	r.GET("report", ginhelpers.WrapContextWeb(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper) (code int, cType string, data []byte) {
			request := reqresp.ReportRequest{
				Project:   nil,
				StartDate: nil,
				DueDate:   nil,
			}

			if val, err := strconv.Atoi(ctx.Query("project")); err == nil {
				request.Project = &val
			}
			if val := ctx.Query("startDate"); val != "" {
				startDate, err := time.Parse("2006-01", val)
				if err != nil {
					nctx.AddPublicError(http.StatusBadRequest, "Invalid start date")
				}
				request.StartDate = &generic.CommonDate{Time: startDate}
			}
			if val := ctx.Query("dueDate"); val != "" {
				dueDate, err := time.Parse("2006-01", val)
				if err != nil {
					nctx.AddPublicError(http.StatusBadRequest, "Invalid start date")
				}
				request.DueDate = &generic.CommonDate{Time: dueDate}
			}
			return http.StatusOK, "text/html", views.RenderReport(nctx, request)
		},
	))

	r.GET("", func(ctx *gin.Context) {
		ctx.Header("Location", "web/index")
		ctx.Status(http.StatusSeeOther)
	})

	r.GET("index", func(ctx *gin.Context) {
		if _, ok := ctx.Get("user_id"); ok {
			ctx.Header("Location", "tasks")
			ctx.Status(http.StatusSeeOther)
		} else {
			ctx.Header("Location", "login")
			ctx.Status(http.StatusSeeOther)
		}
	})
}
