package main

import (
	"context"
	"flag"

	"github.com/adaggerboy/utasksd/config"
	"github.com/adaggerboy/utasksd/internal/routes"
	"github.com/adaggerboy/utasksd/models/generic"
	"github.com/adaggerboy/utasksd/pkg/database"
	_ "github.com/adaggerboy/utasksd/pkg/init"
	"github.com/adaggerboy/utasksd/pkg/utils/jwt"
	"github.com/adaggerboy/utasksd/pkg/utils/renderer"
	"github.com/gin-gonic/gin"
)

func migrate() {
	db, err := database.GetDatabase(context.Background())
	if err != nil {
		panic(err)
	}

	err = db.DeploySchema()
	if err != nil {
		db.RollbackClose()
		panic(err)
	}
	superadmin := generic.AllocateUser()
	*superadmin.Username = config.GlobalConfig.Database.User
	*superadmin.Firstname = "Super"
	*superadmin.Lastname = "Admin"
	*superadmin.Email = "superadmin@local.net"
	*superadmin.AvatarPath = "null"

	id, err := db.CreateUser(context.Background(), *superadmin)
	if err != nil {
		db.RollbackClose()
		panic(err)
	}
	if id != nil {
		err = db.SetUserPermissions(context.Background(), *id, *superadmin.Username, true, true, true)
		if err != nil {
			db.RollbackClose()
			panic(err)
		}
	}
	err = db.Close()
	if err != nil {
		panic(err)
	}

}

func loadTemplates() {
	err := renderer.LoadToPool(map[string]string{
		"task":              config.GlobalConfig.TemplatesLocation + "/task.tpl",
		"project":           config.GlobalConfig.TemplatesLocation + "/project.tpl",
		"issue":             config.GlobalConfig.TemplatesLocation + "/issue.tpl",
		"login":             config.GlobalConfig.TemplatesLocation + "/login.tpl",
		"register":          config.GlobalConfig.TemplatesLocation + "/register.tpl",
		"list_view":         config.GlobalConfig.TemplatesLocation + "/list_view.tpl",
		"user":              config.GlobalConfig.TemplatesLocation + "/user.tpl",
		"issue_list_view":   config.GlobalConfig.TemplatesLocation + "/issue_list_view.tpl",
		"project_list_view": config.GlobalConfig.TemplatesLocation + "/project_list_view.tpl",
		"users_list_view":   config.GlobalConfig.TemplatesLocation + "/users_list_view.tpl",
		"create_project":    config.GlobalConfig.TemplatesLocation + "/create_project.tpl",
		"create_report":     config.GlobalConfig.TemplatesLocation + "/create_report.tpl",
		"report_view":       config.GlobalConfig.TemplatesLocation + "/report_view.tpl",
	})
	if err != nil {
		panic(err)
	}
}

func main() {

	configPath := flag.String("c", "/etc/utasksd/config.yaml", "path to the config file")
	flag.Parse()

	config.Apply(config.Load(*configPath))

	loadTemplates()

	migrate()

	jwt.InitJWT(config.GlobalConfig.JWT)

	engine := gin.New()
	engine.Use(routes.GetAuthorizationMiddleware(config.GlobalConfig.HTTPServer.Domain, config.GlobalConfig.HTTPServer.CookieSecure))
	engine.Use(routes.GetCORSMiddleware(config.GlobalConfig.HTTPServer.Domain))
	engine.Static("static", config.GlobalConfig.StaticLocation)
	routes.DeployWebRoutes(engine.Group("web"))
	routes.DeployAPIRoutes(engine.Group("api/v1"))
	routes.DeployAttachmentsRoutes(engine.Group("files"), config.GlobalConfig.DataLocation)
	routes.DeployRootRoutes(engine.Group(""), config.GlobalConfig.StaticLocation)
	err := engine.Run(config.GlobalConfig.HTTPServer.Endpoints[0])
	if err != nil {
		panic(err)
	}
}
