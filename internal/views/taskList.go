package views

import (
	"fmt"
	"net/http"

	"github.com/adaggerboy/utasksd/internal/controllers"
	"github.com/adaggerboy/utasksd/models/context"
	"github.com/adaggerboy/utasksd/models/generic"
	"github.com/adaggerboy/utasksd/models/rendering"
	"github.com/adaggerboy/utasksd/models/reqresp"
)

//UI Controllers

func RenderTaskListAndTree(ctx context.IContext, tasks []int, projectName string, projectID int) (node []rendering.TreeNode, list []rendering.ListTask) {

	cacheUsers := map[int]generic.User{}

	getUser := func(id int) generic.User {
		user, ok := cacheUsers[id]
		if ok {
			return user
		}
		if userPtr := controllers.GetUserInfo(ctx, id); userPtr != nil {
			cacheUsers[userPtr.ID] = *userPtr
			return cacheUsers[userPtr.ID]
		} else {
			return generic.User{}
		}
	}

	node = []rendering.TreeNode{}
	list = []rendering.ListTask{}

	for _, l := range tasks {
		task := controllers.GetTaskInfo(ctx, l, true)
		if !ctx.IsActive() {
			return
		}
		node = append(node, rendering.TreeNode{
			ID:           l,
			Name:         *task.Name,
			VisibleClass: string(*task.Status) + "-node",
		})
		assigner := getUser(*task.Assigner)
		if !ctx.IsActive() {
			return
		}

		listTask := rendering.ListTask{
			FirePriority: false,
			ProjectName:  projectName,
			ProjectID:    projectID,
			Name:         *task.Name,
			Description:  *task.Description,
			Status:       GetRendererCoreStruct().Statuses[string(*task.Status)],
			VisibleClass: string(*task.Status) + "-task",
			ID:           l,
			Assigner: rendering.Account{
				PathToAvatar: *assigner.AvatarPath,
				Name:         *assigner.Firstname + " " + *assigner.Lastname,
			},
			Assignees: []rendering.Account{},
		}
		for _, s := range *task.Assignees {
			assignee := getUser(s)
			if !ctx.IsActive() {
				return
			}
			listTask.Assignees = append(listTask.Assignees, rendering.Account{
				PathToAvatar: *assignee.AvatarPath,
				Name:         *assignee.Firstname + " " + *assignee.Lastname,
			})
		}
		list = append(list, listTask)
	}
	return
}

func RenderListTaskAndTreeForProject(ctx context.IContext, projectID int, userID *int) (node rendering.TreeNode, list []rendering.ListTask) {

	project := controllers.GetProjectInfo(ctx, projectID, false)
	if !ctx.IsActive() {
		return
	}

	var tasks []int

	if userID != nil {
		tasks = controllers.GetUserProjectTasks(ctx, *userID, projectID)
	} else {
		tasks = controllers.GetProjectTasks(ctx, projectID)
	}
	if !ctx.IsActive() {
		return
	}

	nodes, list := RenderTaskListAndTree(ctx, tasks, *project.Name, projectID)
	if !ctx.IsActive() {
		return
	}

	node = rendering.TreeNode{
		ID:       projectID,
		Name:     *project.Name,
		Children: nil,
	}

	if len(nodes) > 0 {
		node.Children = &[]rendering.TreeNode{}
		(*node.Children) = append((*node.Children), nodes...)
	}

	return
}

func RenderProjectTaskListView(ctx context.IContext, projectID int) (data []byte) {

	core := GetRendererCoreStruct()
	user := RenderCoreInfo(ctx, &core)
	if !ctx.IsActive() {
		return
	}

	r, t, i, m := controllers.GetUserProjectCapabilities(ctx, user.ID, projectID)
	if !r {
		ctx.AddPublicError(http.StatusForbidden, "Access denied")
		return
	}

	project := controllers.GetProjectInfo(ctx, projectID, true)
	if !ctx.IsActive() {
		return
	}

	input := rendering.ListViewRoot{
		Core:              core,
		Tasks:             []rendering.ListTask{},
		AvailableProjects: []rendering.TreeNode{},
		ProjectData:       RenderProjectData(ctx, projectID, project),
	}

	if t {
		input.Core.Options = append(input.Core.Options, rendering.Option{
			Href: fmt.Sprintf("/web/project/%d/create_task", projectID), Label: "+ Create task",
		})
	}
	if i {
		input.Core.Options = append(input.Core.Options, rendering.Option{
			Href: fmt.Sprintf("/web/project/%d/create_issue", projectID), Label: "+ Register issue",
		})
	}
	if m {
		input.Core.Options = append(input.Core.Options, rendering.Option{
			Href: fmt.Sprintf("/web/project/%d/create_report", projectID), Label: "+ Create report",
		})
	}

	node, list := RenderListTaskAndTreeForProject(ctx, projectID, nil)
	if !ctx.IsActive() {
		return
	}
	input.Tasks = append(input.Tasks, list...)
	input.AvailableProjects = append(input.AvailableProjects, node)

	return RenderAs(ctx, "list_view", input)
}

func RenderUserTaskListView(ctx context.IContext) (data []byte) {

	core := GetRendererCoreStruct()
	user := RenderCoreInfo(ctx, &core)
	if !ctx.IsActive() {
		return
	}

	input := rendering.ListViewRoot{
		Core: core,
		ProjectData: rendering.ProjectData{
			InProject: false,
			Assigners: map[int]string{},
			Assignees: map[int]string{},
		},
		Tasks:             []rendering.ListTask{},
		AvailableProjects: []rendering.TreeNode{},
	}

	projects := controllers.GetUserProjects(ctx, user.ID)
	if !ctx.IsActive() {
		return
	}

	for _, v := range projects {
		node, list := RenderListTaskAndTreeForProject(ctx, v, &user.ID)
		if !ctx.IsActive() {
			return
		}
		input.Tasks = append(input.Tasks, list...)
		input.AvailableProjects = append(input.AvailableProjects, node)
	}

	return RenderAs(ctx, "list_view", input)

}

func RenderSearchProjectTaskListView(ctx context.IContext, search reqresp.SearchTasks) (data []byte) {

	core := GetRendererCoreStruct()
	RenderCoreInfo(ctx, &core)
	if !ctx.IsActive() {
		return
	}

	if search.Project == nil {
		ctx.AddPublicError(http.StatusBadRequest, "Project specification required")
	}

	project := controllers.GetProjectInfo(ctx, *search.Project, true)
	if !ctx.IsActive() {
		return
	}

	input := rendering.ListViewRoot{
		Core:              core,
		Tasks:             []rendering.ListTask{},
		AvailableProjects: []rendering.TreeNode{},
		ProjectData:       RenderProjectData(ctx, *search.Project, project),
	}

	tasks := controllers.SearchTasks(ctx, search)
	if !ctx.IsActive() {
		return
	}

	nodes, list := RenderTaskListAndTree(ctx, tasks, *project.Name, *search.Project)
	if !ctx.IsActive() {
		return
	}

	node := rendering.TreeNode{
		ID:       *search.Project,
		Name:     *project.Name,
		Children: nil,
	}

	if len(nodes) > 0 {
		node.Children = &[]rendering.TreeNode{}
		(*node.Children) = append((*node.Children), nodes...)
	}
	input.Tasks = append(input.Tasks, list...)
	input.AvailableProjects = append(input.AvailableProjects, node)

	return RenderAs(ctx, "list_view", input)

}
