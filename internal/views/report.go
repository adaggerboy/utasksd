package views

import (
	"strconv"

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
		TimeEfficiency: []rendering.TimeEfficiencyReport{},
	}

	for _, v := range report2 {
		input.TimeEfficiency = append(input.TimeEfficiency, rendering.TimeEfficiencyReport{
			ID:              v.ID,
			Firstname:       v.Firstname,
			Username:        v.Username,
			Lastname:        v.Lastname,
			TrackedRecords:  v.TrackedRecords,
			SummaryHours:    strconv.FormatFloat(v.SummaryHours, 'f', 2, 64),
			AveragePerDay:   strconv.FormatFloat(v.AveragePerDay, 'f', 2, 64),
			AveragePerMonth: strconv.FormatFloat(v.AveragePerMonth, 'f', 2, 64),
		})
	}

	return RenderAs(ctx, "report_view", input)
}
