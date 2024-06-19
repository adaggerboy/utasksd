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

func RenderIssueListAndTree(ctx context.IContext, issues []int, projectName string, projectID int) (node []rendering.TreeNode, list []rendering.ListIssue) {

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
	list = []rendering.ListIssue{}

	for _, l := range issues {
		issue := controllers.GetIssueInfo(ctx, l, true)
		if !ctx.IsActive() {
			return
		}
		node = append(node, rendering.TreeNode{
			ID:           l,
			Name:         *issue.Name,
			VisibleClass: string(*issue.Status) + "-node",
		})
		assigner := getUser(*issue.Registrar)
		if !ctx.IsActive() {
			return
		}

		listIssue := rendering.ListIssue{
			ProjectName:  projectName,
			ProjectID:    projectID,
			Name:         *issue.Name,
			Description:  *issue.Description,
			Status:       GetRendererCoreStruct().IssueStatuses[string(*issue.Status)],
			VisibleClass: string(*issue.Status) + "-issue",
			ID:           l,
			Registrar: rendering.Account{
				PathToAvatar: *assigner.AvatarPath,
				Name:         *assigner.Firstname + " " + *assigner.Lastname,
			},
		}

		if issue.Reporter != nil {
			listIssue.Reporter = *issue.Reporter
		}
		list = append(list, listIssue)
	}
	return
}

func RenderListIssueAndTreeForProject(ctx context.IContext, projectID int, userID *int) (node rendering.TreeNode, list []rendering.ListIssue) {

	project := controllers.GetProjectInfo(ctx, projectID, false)
	if !ctx.IsActive() {
		return
	}

	issues := controllers.GetProjectIssues(ctx, projectID)
	if !ctx.IsActive() {
		return
	}

	nodes, list := RenderIssueListAndTree(ctx, issues, *project.Name, projectID)
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

func RenderProjectIssueListView(ctx context.IContext, projectID int) (data []byte) {

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

	input := rendering.ListIssueViewRoot{
		Core:              core,
		Issues:            []rendering.ListIssue{},
		AvailableProjects: []rendering.TreeNode{},
		ProjectData:       RenderProjectData(ctx, projectID, project),
	}

	if t {
		input.Core.Options = append(input.Core.Options, rendering.Option{
			Href: fmt.Sprintf("/web/project/%d/create_issue", projectID), Label: "+ Create issue",
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

	node, list := RenderListIssueAndTreeForProject(ctx, projectID, nil)
	if !ctx.IsActive() {
		return
	}
	input.Issues = append(input.Issues, list...)
	input.AvailableProjects = append(input.AvailableProjects, node)

	return RenderAs(ctx, "issue_list_view", input)
}

func RenderUserIssueListView(ctx context.IContext) (data []byte) {

	core := GetRendererCoreStruct()
	user := RenderCoreInfo(ctx, &core)
	if !ctx.IsActive() {
		return
	}

	input := rendering.ListIssueViewRoot{
		Core: core,
		ProjectData: rendering.ProjectData{
			InProject: false,
			Assigners: map[int]string{},
			Assignees: map[int]string{},
		},
		Issues:            []rendering.ListIssue{},
		AvailableProjects: []rendering.TreeNode{},
	}

	projects := controllers.GetUserProjects(ctx, user.ID)
	if !ctx.IsActive() {
		return
	}

	for _, v := range projects {
		node, list := RenderListIssueAndTreeForProject(ctx, v, &user.ID)
		if !ctx.IsActive() {
			return
		}
		input.Issues = append(input.Issues, list...)
		input.AvailableProjects = append(input.AvailableProjects, node)
	}

	return RenderAs(ctx, "issue_list_view", input)

}

func RenderSearchProjectIssueListView(ctx context.IContext, search reqresp.SearchIssues) (data []byte) {

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

	input := rendering.ListIssueViewRoot{
		Core:              core,
		Issues:            []rendering.ListIssue{},
		AvailableProjects: []rendering.TreeNode{},
		ProjectData:       RenderProjectData(ctx, *search.Project, project),
	}

	issues := controllers.SearchIssues(ctx, search)
	if !ctx.IsActive() {
		return
	}

	nodes, list := RenderIssueListAndTree(ctx, issues, *project.Name, *search.Project)
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
	input.Issues = append(input.Issues, list...)
	input.AvailableProjects = append(input.AvailableProjects, node)

	return RenderAs(ctx, "issue_list_view", input)

}
