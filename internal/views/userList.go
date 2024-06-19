package views

import (
	"github.com/adaggerboy/utasksd/internal/controllers"
	"github.com/adaggerboy/utasksd/models/context"
	"github.com/adaggerboy/utasksd/models/rendering"
)

//UI Controllers

func RenderUserListView(ctx context.IContext) (data []byte) {

	core := GetRendererCoreStruct()
	user := RenderCoreInfo(ctx, &core)
	if !ctx.IsActive() {
		return
	}

	users := controllers.GetAllUsers(ctx)
	if !ctx.IsActive() {
		return
	}

	input := rendering.ListUserViewRoot{
		Core:    core,
		IsAdmin: *user.IsAdmin,
		Users:   []rendering.GenericUserWrapper{},
	}

	for _, v := range users {
		input.Users = append(input.Users, rendering.GenericUserWrapper{
			User:     v,
			Admin:    *v.IsAdmin,
			Active:   *v.IsActive,
			Director: *v.IsDirector,
		})
	}

	return RenderAs(ctx, "users_list_view", input)
}
