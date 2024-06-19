package views

import (
	"github.com/adaggerboy/utasksd/internal/controllers"
	"github.com/adaggerboy/utasksd/models/context"
	"github.com/adaggerboy/utasksd/models/rendering"
	"github.com/adaggerboy/utasksd/models/reqresp"
)

func RenderReport(ctx context.IContext, request reqresp.ReportRequest) (data []byte) {
	report1 := controllers.GetTasksReport(ctx, request)
	if !ctx.IsActive() {
		return
	}
	report2 := controllers.GetTimeEfficiencyReport(ctx, request)
	if !ctx.IsActive() {
		return
	}

	input := rendering.ReportRoot{
		Tasks:          report1,
		TimeEfficiency: report2,
	}

	return RenderAs(ctx, "report_view", input)
}
