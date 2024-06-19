package views

import (
	"time"

	"github.com/adaggerboy/utasksd/internal/controllers"
	"github.com/adaggerboy/utasksd/models/context"
	"github.com/adaggerboy/utasksd/models/generic"
	"github.com/adaggerboy/utasksd/models/rendering"
)

func RenderIssue(ctx context.IContext, issueID int, issue *generic.Issue, projectOwner int) rendering.Issue {

	userID, ok := ctx.GetUserID()
	if !ok {
		return rendering.Issue{}
	}

	assigner := controllers.GetUserInfo(ctx, *issue.Registrar)
	if !ctx.IsActive() {
		return rendering.Issue{}
	}

	result := rendering.Issue{
		Name:        *issue.Name,
		Description: *issue.Description,
		ID:          issueID,
		Attachments: []string{},
		StartDate:   (*issue.StartDate).Format("2006-01-02"),
		DueDate:     "",
		Registrar:   *assigner.Firstname + " " + *assigner.Lastname,
		Reporter:    "",
		Status:      string(*issue.Status),
		Closed:      *issue.Status == "closed",

		Created: true,

		EditCapability: userID == projectOwner || userID == *issue.Registrar,
	}

	if issue.Reporter != nil {
		result.Reporter = *issue.Reporter

	}

	if issue.CloseDate != nil {
		result.DueDate = (*issue.CloseDate).Format("2006-01-02")
	}

	if issue.Attachments != nil {
		result.Attachments = *issue.Attachments
	}

	return result
}

//UI Controllers

func RenderIssueView(ctx context.IContext, issueID int) (data []byte) {

	core := GetRendererCoreStruct()
	RenderCoreInfo(ctx, &core)
	if !ctx.IsActive() {
		return
	}
	issue := controllers.GetIssueInfo(ctx, issueID, true)
	if !ctx.IsActive() {
		return
	}

	project := controllers.GetProjectInfo(ctx, *issue.Project, true)
	if !ctx.IsActive() {
		return
	}

	renderIssue := RenderIssue(ctx, issueID, issue, *project.OwnerID)
	if !ctx.IsActive() {
		return
	}

	input := rendering.IssueRoot{
		Core:        core,
		ProjectData: RenderProjectData(ctx, *issue.Project, project),
		Issue:       renderIssue,
	}
	if !ctx.IsActive() {
		return
	}

	return RenderAs(ctx, "issue", input)
}

func RenderCreateIssueView(ctx context.IContext, projectID int) (data []byte) {

	core := GetRendererCoreStruct()
	user := RenderCoreInfo(ctx, &core)
	if !ctx.IsActive() {
		return
	}

	project := controllers.GetProjectInfo(ctx, projectID, true)
	if !ctx.IsActive() {
		return
	}

	input := rendering.IssueRoot{
		Core:        core,
		ProjectData: RenderProjectData(ctx, projectID, project),
		Issue: rendering.Issue{
			Name:        "New issue",
			Reporter:    "",
			Registrar:   *user.Firstname + " " + *user.Lastname,
			Description: "Click on task name or on this text to edit...",
			ID:          0,
			Attachments: []string{},
			StartDate:   time.Now().Format("2006-01-02"),
			DueDate:     "",
			Status:      string(generic.Open),
			Created:     false,
			Closed:      false,

			EditCapability: true,
		},
	}
	if !ctx.IsActive() {
		return
	}

	return RenderAs(ctx, "issue", input)
}
