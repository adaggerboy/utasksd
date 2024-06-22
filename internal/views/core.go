package views

import (
	"fmt"
	"net/http"

	"github.com/adaggerboy/utasksd/internal/controllers"
	"github.com/adaggerboy/utasksd/models/context"
	"github.com/adaggerboy/utasksd/models/generic"
	"github.com/adaggerboy/utasksd/models/rendering"
	"github.com/adaggerboy/utasksd/pkg/utils/renderer"
)

func GetRendererCoreStruct() rendering.Core {
	return rendering.Core{
		Statuses: map[string]string{
			string(generic.ToDo):       "To do",
			string(generic.InProgress): "In progress",
			string(generic.Waiting):    "On waiting",
			string(generic.Testing):    "Testing",
			string(generic.Done):       "Done",
		},
		IssueStatuses: map[string]string{
			string(generic.Open):     "Open",
			string(generic.Closed):   "closed",
			string(generic.Reopened): "Reopened",
		},
		Priorities: map[string]string{
			string(generic.Highest): "Highest",
			string(generic.High):    "High",
			string(generic.Middle):  "Middle",
			string(generic.Low):     "Low",
			string(generic.Lowest):  "Lowest",
		},
		TaskDependencyTypes: map[string]string{
			string(generic.BlockedBy): "Blocked by",
			string(generic.Includes):  "Include",
		},
		CSSVariables: map[string]string{
			"accent1":      "#3048a7",
			"accent2":      "#3f51b5",
			"accent2hover": "#5f71d7",
			"accent3":      "#afb5d7",

			"task-done":       "#3ab847",
			"task-todo":       "#acacac",
			"task-inprogress": "#3d7ec8",
			"task-waiting":    "#c8a13d",
			"task-testing":    "#803dc8",
		},
		JSVariablesJSON: `{api_base_link: "http://localhost:10200"}`,
	}
}

func RenderCoreInfo(ctx context.IContext, core *rendering.Core) (user *generic.User) {
	user = controllers.GetMyUserInfo(ctx)
	if !ctx.IsActive() {
		return
	}
	core.Account = rendering.Account{
		Name:         *user.Firstname + " " + *user.Lastname,
		PathToAvatar: *user.AvatarPath,
		ID:           user.ID,
	}

	core.Options = []rendering.Option{
		{Href: "/web/tasks", Label: "Tasks"},
		{Href: "/web/issues", Label: "Issues"},
		{Href: "/web/projects", Label: "Projects"},
		{Href: "/web/users", Label: "Users"},
	}

	if *user.IsDirector {
		core.Options = append(core.Options, rendering.Option{Href: "/web/create_project", Label: "+ Create project"})
	}
	return
}

func RenderAs(ctx context.IContext, templateName string, input any) (data []byte) {
	rend, err := renderer.GetRenderer(templateName)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("loading renderer RenderTaskView(): %s", err))
		return
	}
	data, err = rend.Render(input)
	if err != nil {
		ctx.AddPrivateError(http.StatusInternalServerError, fmt.Errorf("rendering RenderTaskView(): %s", err))
		return
	}
	return data
}

func RenderProjectData(ctx context.IContext, projectID int, project *generic.Project) rendering.ProjectData {
	result := rendering.ProjectData{
		Assigners: map[int]string{},
		Assignees: map[int]string{},
		Supports:  map[int]string{},
		ID:        projectID,
		InProject: true,
	}

	owner := controllers.GetUserInfo(ctx, *project.OwnerID)
	if !ctx.IsActive() {
		return result
	}

	result.Assigners[*project.OwnerID] = *owner.Firstname + " " + *owner.Lastname
	result.Assignees[*project.OwnerID] = *owner.Firstname + " " + *owner.Lastname
	result.Supports[*project.OwnerID] = *owner.Firstname + " " + *owner.Lastname

	for _, v := range *project.Managers {
		manager := controllers.GetUserInfo(ctx, v)
		if !ctx.IsActive() {
			return result
		}
		result.Assigners[v] = *manager.Firstname + " " + *manager.Lastname
		result.Assignees[v] = *manager.Firstname + " " + *manager.Lastname
	}

	for _, v := range *project.SupportAgents {
		manager := controllers.GetUserInfo(ctx, v)
		if !ctx.IsActive() {
			return result
		}
		result.Supports[v] = *manager.Firstname + " " + *manager.Lastname
	}

	for _, v := range *project.Workers {
		worker := controllers.GetUserInfo(ctx, v)
		if !ctx.IsActive() {
			return result
		}
		result.Assignees[v] = *worker.Firstname + " " + *worker.Lastname
	}
	return result
}
