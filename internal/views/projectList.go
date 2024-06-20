package views

import (
	"github.com/adaggerboy/utasksd/internal/controllers"
	"github.com/adaggerboy/utasksd/models/context"
	"github.com/adaggerboy/utasksd/models/rendering"
)

//UI Controllers

func RenderProjectListView(ctx context.IContext) (data []byte) {

	core := GetRendererCoreStruct()
	user := RenderCoreInfo(ctx, &core)
	if !ctx.IsActive() {
		return
	}

	projects := controllers.GetUserProjects(ctx, user.ID)
	if !ctx.IsActive() {
		return
	}

	input := rendering.ListProjectViewRoot{
		Core:     core,
		Projects: []rendering.Project{},
	}

	for _, v := range projects {
		role := "none"
		proj := controllers.GetProjectInfo(ctx, v, true)
		if !ctx.IsActive() {
			return
		}
		owner := controllers.GetUserInfo(ctx, *proj.OwnerID)
		if !ctx.IsActive() {
			return
		}
		if user.ID == *proj.OwnerID {
			role = "owner"
		} else {
			for _, k := range *proj.Workers {
				if user.ID == k {
					role = "worker"
				}
			}
			for _, k := range *proj.SupportAgents {
				if user.ID == k {
					role = "support agent"
				}
			}
			for _, k := range *proj.Managers {
				if user.ID == k {
					role = "manager"
				}
			}
		}
		input.Projects = append(input.Projects, rendering.Project{
			Name:        *proj.Name,
			ID:          v,
			Description: *proj.Description,
			Role:        role,
			Logo:        *proj.LogoPath,
			Owner: rendering.Account{
				Name:         *owner.Firstname + " " + *owner.Lastname,
				PathToAvatar: *owner.AvatarPath,
			},
		})
	}

	return RenderAs(ctx, "project_list_view", input)
}

func RenderProjectView(ctx context.IContext, projectID int) (data []byte) {

	core := GetRendererCoreStruct()
	user := RenderCoreInfo(ctx, &core)
	if !ctx.IsActive() {
		return
	}

	project := controllers.GetProjectInfo(ctx, projectID, true)
	if !ctx.IsActive() {
		return
	}

	input := rendering.ProjectRoot{
		Core:    core,
		IsOwner: *project.OwnerID == user.ID,
		Project: rendering.Project{
			Name:        *project.Name,
			ID:          projectID,
			Description: *project.Description,
			Logo:        *project.LogoPath,
			Users:       []rendering.Account{},
		},
	}

	for _, k := range *project.Workers {
		projUser := controllers.GetUserInfo(ctx, k)
		if !ctx.IsActive() {
			return
		}
		input.Project.Users = append(input.Project.Users, rendering.Account{
			ID:           k,
			Name:         *projUser.Firstname + " " + *projUser.Lastname,
			Email:        *projUser.Email,
			PathToAvatar: *projUser.AvatarPath,
			Role:         "worker",
		})
	}

	for _, k := range *project.Managers {
		projUser := controllers.GetUserInfo(ctx, k)
		if !ctx.IsActive() {
			return
		}
		input.Project.Users = append(input.Project.Users, rendering.Account{
			ID:           k,
			Name:         *projUser.Firstname + " " + *projUser.Lastname,
			Email:        *projUser.Email,
			PathToAvatar: *projUser.AvatarPath,
			Role:         "manager",
		})
	}

	for _, k := range *project.SupportAgents {
		projUser := controllers.GetUserInfo(ctx, k)
		if !ctx.IsActive() {
			return
		}
		input.Project.Users = append(input.Project.Users, rendering.Account{
			ID:           k,
			Name:         *projUser.Firstname + " " + *projUser.Lastname,
			Email:        *projUser.Email,
			PathToAvatar: *projUser.AvatarPath,
			Role:         "support",
		})
	}

	return RenderAs(ctx, "project", input)
}
