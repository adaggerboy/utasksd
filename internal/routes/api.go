package routes

import (
	"net/http"

	"github.com/adaggerboy/utasksd/config"
	"github.com/adaggerboy/utasksd/internal/controllers"
	"github.com/adaggerboy/utasksd/models/generic"
	"github.com/adaggerboy/utasksd/models/reqresp"
	ginhelpers "github.com/adaggerboy/utasksd/pkg/utils/ginHelpers"
	"github.com/gin-gonic/gin"
)

func DeployTaskRoutes(r gin.IRouter) {
	r.GET(":id", ginhelpers.WrapContextID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int) (code int, obj interface{}) {
			return http.StatusOK, controllers.GetTaskInfo(nctx, id, false)
		},
	))

	r.GET(":id/detailed", ginhelpers.WrapContextID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int) (code int, obj interface{}) {
			return http.StatusOK, controllers.GetTaskInfo(nctx, id, true)
		},
	))

	r.POST(":id/tracking-record", ginhelpers.WrapContextDataID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int, tr generic.TrackingRecord) (code int, obj interface{}) {
			if tr.Duration == nil || tr.Text == nil {
				nctx.AddPublicError(http.StatusBadRequest, "required fields was not provided")
				return http.StatusBadRequest, nil
			}
			controllers.TrackTaskActivity(nctx, id, *tr.Text, tr.Duration.Duration)
			return http.StatusOK, nil
		},
	))

	r.POST("", ginhelpers.WrapContextData(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, task generic.Task) (code int, obj interface{}) {
			return http.StatusOK, map[string]int{"id": controllers.CreateTask(nctx, task)}
		},
	))

	r.PUT(":id", ginhelpers.WrapContextDataID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int, task generic.Task) (code int, obj interface{}) {
			controllers.UpdateTask(nctx, id, task)
			return http.StatusOK, nil
		},
	))

	r.DELETE(":id", ginhelpers.WrapContextID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int) (code int, obj interface{}) {
			controllers.DeleteTask(nctx, id)
			return http.StatusOK, nil
		},
	))
}

func DeployUserRoutes(r gin.IRouter) {
	r.GET(":id", ginhelpers.WrapContextID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int) (code int, obj interface{}) {
			return http.StatusOK, controllers.GetUserInfo(nctx, id)
		},
	))

	r.GET("", ginhelpers.WrapContext(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper) (code int, obj interface{}) {
			return http.StatusOK, controllers.GetMyUserInfo(nctx)
		},
	))

	r.PUT(":id", ginhelpers.WrapContextDataID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int, user generic.User) (code int, obj interface{}) {
			controllers.UpdateUser(nctx, id, user)
			return http.StatusOK, nil
		},
	))

	r.DELETE(":id", ginhelpers.WrapContextID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int) (code int, obj interface{}) {
			controllers.DeleteUser(nctx, id)
			return http.StatusOK, nil
		},
	))

	r.PUT("", ginhelpers.WrapContextData(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, user generic.User) (code int, obj interface{}) {
			controllers.UpdateMyUser(nctx, user)
			return http.StatusOK, nil
		},
	))

	r.DELETE("", ginhelpers.WrapContext(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper) (code int, obj interface{}) {
			ginhelpers.ResetToken(ctx, config.GlobalConfig.HTTPServer.Domain, config.GlobalConfig.HTTPServer.CookieSecure)
			controllers.DeleteMyUser(nctx)
			return http.StatusOK, nil
		},
	))

	r.GET(":id/tasks", ginhelpers.WrapContextID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int) (code int, obj interface{}) {
			return http.StatusOK, controllers.GetUserTasks(nctx, id)
		},
	))

	r.GET(":id/projects", ginhelpers.WrapContextID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int) (code int, obj interface{}) {
			return http.StatusOK, controllers.GetUserProjects(nctx, id)
		},
	))

}

func DeployAuthRoutes(r gin.IRouter) {
	r.POST("with-password", ginhelpers.WrapContextDataPermitted(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, req reqresp.CreateUserWithPasswordRequest) (code int, obj interface{}) {
			return http.StatusOK, map[string]int{"id": controllers.CreateUserWithPassword(nctx, req)}
		},
	))

	r.PUT(":id/credentials", ginhelpers.WrapContextDataID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int, creds reqresp.UpdatePasswordRequest) (code int, obj interface{}) {
			controllers.UpdateCredentials(nctx, id, creds)
			return http.StatusOK, nil
		},
	))

	r.PUT(":id/permissions", ginhelpers.WrapContextDataID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int, user generic.User) (code int, obj interface{}) {
			controllers.SetUserPermissions(nctx, id, user)
			return http.StatusOK, nil
		},
	))

	r.PUT("credentials", ginhelpers.WrapContextData(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, creds reqresp.UpdatePasswordRequest) (code int, obj interface{}) {
			controllers.UpdateMyCredentials(nctx, creds)
			return http.StatusOK, nil
		},
	))

	r.Match([]string{"GET", "POST"}, "session", ginhelpers.WrapContextDataPermitted(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, login reqresp.LoginWithPasswordRequest) (code int, obj interface{}) {
			token := controllers.Login(nctx, login.Username, login.Secret)
			if nctx.IsActive() {
				ginhelpers.SetToken(ctx, config.GlobalConfig.HTTPServer.Domain, config.GlobalConfig.HTTPServer.CookieSecure, config.GlobalConfig.HTTPServer.TokenDuration, token)
				return http.StatusOK, map[string]string{"accessToken": token}
			}
			return
		},
	))
	r.DELETE("session", ginhelpers.WrapContext(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper) (code int, obj interface{}) {
			ginhelpers.ResetToken(ctx, config.GlobalConfig.HTTPServer.Domain, config.GlobalConfig.HTTPServer.CookieSecure)
			return http.StatusOK, nil
		},
	))
}

func DeployProjectRoutes(r gin.IRouter) {
	r.GET(":id", ginhelpers.WrapContextID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int) (code int, obj interface{}) {
			return http.StatusOK, controllers.GetProjectInfo(nctx, id, false)
		},
	))

	r.GET(":id/detailed", ginhelpers.WrapContextID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int) (code int, obj interface{}) {
			return http.StatusOK, controllers.GetProjectInfo(nctx, id, true)
		},
	))

	r.POST("", ginhelpers.WrapContextData(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, project generic.Project) (code int, obj interface{}) {
			return http.StatusOK, map[string]int{"id": controllers.CreateProject(nctx, project)}
		},
	))

	r.PUT(":id", ginhelpers.WrapContextDataID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int, project generic.Project) (code int, obj interface{}) {
			controllers.UpdateProject(nctx, id, project)
			return http.StatusOK, nil
		},
	))

	r.DELETE(":id", ginhelpers.WrapContextID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int) (code int, obj interface{}) {
			controllers.DeleteProject(nctx, id)
			return http.StatusOK, nil
		},
	))

	r.GET(":id/tasks", ginhelpers.WrapContextID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int) (code int, obj interface{}) {
			return http.StatusOK, controllers.GetProjectTasks(nctx, id)
		},
	))

	r.GET(":id/issues", ginhelpers.WrapContextID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int) (code int, obj interface{}) {
			return http.StatusOK, controllers.GetProjectIssues(nctx, id)
		},
	))
}

func DeployCommentRoutes(r gin.IRouter) {
	r.GET(":id", ginhelpers.WrapContextID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int) (code int, obj interface{}) {
			return http.StatusOK, controllers.ReadComment(nctx, id)
		},
	))

	r.POST("", ginhelpers.WrapContextData(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, comment generic.Comment) (code int, obj interface{}) {
			if comment.TaskID == nil || comment.Text == nil {
				nctx.AddPublicError(http.StatusBadRequest, "required fields was not provided")
				return http.StatusBadRequest, nil
			}
			return http.StatusOK, map[string]int{"id": controllers.PublishComment(nctx, *comment.TaskID, *comment.Text)}
		},
	))

	r.DELETE(":id", ginhelpers.WrapContextID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int) (code int, obj interface{}) {
			controllers.DeleteComment(nctx, id)
			return http.StatusOK, nil
		},
	))
}

func DeployIssueRoutes(r gin.IRouter) {
	r.GET(":id", ginhelpers.WrapContextID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int) (code int, obj interface{}) {
			return http.StatusOK, controllers.GetIssueInfo(nctx, id, false)
		},
	))

	r.GET(":id/detailed", ginhelpers.WrapContextID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int) (code int, obj interface{}) {
			return http.StatusOK, controllers.GetIssueInfo(nctx, id, true)
		},
	))

	r.POST("", ginhelpers.WrapContextData(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, issue generic.Issue) (code int, obj interface{}) {
			return http.StatusOK, map[string]int{"id": controllers.CreateIssue(nctx, issue)}
		},
	))

	r.PUT(":id", ginhelpers.WrapContextDataID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int, issue generic.Issue) (code int, obj interface{}) {
			controllers.UpdateIssue(nctx, id, issue)
			return http.StatusOK, nil
		},
	))

	r.DELETE(":id", ginhelpers.WrapContextID(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, id int) (code int, obj interface{}) {
			controllers.DeleteIssue(nctx, id)
			return http.StatusOK, nil
		},
	))
}

func DeploySearchRoutes(r gin.IRouter) {
	r.GET("tasks", ginhelpers.WrapContextData(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, query reqresp.SearchTasks) (code int, obj interface{}) {
			return http.StatusOK, controllers.SearchTasks(nctx, query)
		},
	))
	r.GET("issues", ginhelpers.WrapContextData(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, query reqresp.SearchIssues) (code int, obj interface{}) {
			return http.StatusOK, controllers.SearchIssues(nctx, query)
		},
	))
}

func DeployReportRoutes(r gin.IRouter) {
	r.GET("tasks", ginhelpers.WrapContextData(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, query reqresp.ReportRequest) (code int, obj interface{}) {
			return http.StatusOK, controllers.GetTasksReport(nctx, query)
		},
	))

	r.GET("time-efficiency", ginhelpers.WrapContextData(
		func(ctx *gin.Context, nctx *ginhelpers.GinWrapper, query reqresp.ReportRequest) (code int, obj interface{}) {
			return http.StatusOK, controllers.GetTimeEfficiencyReport(nctx, query)
		},
	))
}

func DeployAPIRoutes(r gin.IRouter) {
	DeployTaskRoutes(r.Group("task"))
	DeployIssueRoutes(r.Group("issue"))
	DeployProjectRoutes(r.Group("project"))
	DeployUserRoutes(r.Group("user"))
	DeployCommentRoutes(r.Group("comment"))
	DeployAuthRoutes(r.Group("auth"))
	DeploySearchRoutes(r.Group("search"))
	DeployReportRoutes(r.Group("report"))
}
