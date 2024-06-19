package views

import (
	"github.com/adaggerboy/utasksd/internal/controllers"
	"github.com/adaggerboy/utasksd/models/context"
	"github.com/adaggerboy/utasksd/models/rendering"
)

func RenderLoginView(ctx context.IContext) (data []byte) {
	return RenderAs(ctx, "login", struct {
		Core            rendering.Core
		JSVariablesJSON string
	}{
		Core:            GetRendererCoreStruct(),
		JSVariablesJSON: "",
	})
}

func RenderRegisterView(ctx context.IContext) (data []byte) {
	return RenderAs(ctx, "register", struct {
		Core            rendering.Core
		JSVariablesJSON string
	}{
		Core:            GetRendererCoreStruct(),
		JSVariablesJSON: "",
	})
}

func RenderCreateProjectView(ctx context.IContext) (data []byte) {
	return RenderAs(ctx, "create_project", struct {
		Core            rendering.Core
		JSVariablesJSON string
	}{
		Core:            GetRendererCoreStruct(),
		JSVariablesJSON: "",
	})
}
func RenderCreateReportView(ctx context.IContext, projectID int) (data []byte) {

	project := controllers.GetProjectInfo(ctx, projectID, false)
	if !ctx.IsActive() {
		return
	}

	return RenderAs(ctx, "create_report", struct {
		Core            rendering.Core
		JSVariablesJSON string
		ID              int
		Name            string
	}{
		Core:            GetRendererCoreStruct(),
		JSVariablesJSON: "",
		ID:              projectID,
		Name:            *project.Name,
	})
}
