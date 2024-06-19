package database

import (
	"context"
	"fmt"

	"github.com/adaggerboy/utasksd/config"
	model "github.com/adaggerboy/utasksd/models/config"
)

type DatabaseFabricFunc[R any] func(ctx context.Context, conf model.DatabaseEndpointConfig, user string, password string) (db R, err error)

type DatabaseFabric[R any] struct {
	readDBFabrics map[string]DatabaseFabricFunc[R]
}

func NewDatabaseFabric[R any]() *DatabaseFabric[R] {
	return &DatabaseFabric[R]{
		readDBFabrics: map[string]DatabaseFabricFunc[R]{},
	}
}

func (c *DatabaseFabric[R]) RegisterDatabaseFabric(key string, fun DatabaseFabricFunc[R]) {
	c.readDBFabrics[key] = fun
}

func (c *DatabaseFabric[R]) NewDatabaseController(ctx context.Context, conf model.DatabaseEndpointConfig, user string, password string) (controller R, err error) {
	readDBFabric, ok := c.readDBFabrics[conf.Driver]
	if !ok {
		err = fmt.Errorf("driver not found: %s", conf.Driver)
		return
	}
	if controller, err = readDBFabric(ctx, conf, user, password); err != nil {
		err = fmt.Errorf("read db controller init: %s", err)
		return
	}
	return
}

var (
	databaseFabric *DatabaseFabric[IDatabase] = NewDatabaseFabric[IDatabase]()
)

func GetDatabaseFabric() *DatabaseFabric[IDatabase] {
	return databaseFabric
}

func ConnectDatabase(ctx context.Context, user string, password string) (db IDatabase, err error) {
	return databaseFabric.NewDatabaseController(ctx, config.GlobalConfig.Database, user, password)
}

func GetDatabase(ctx context.Context) (db IDatabase, err error) {
	return databaseFabric.NewDatabaseController(ctx, config.GlobalConfig.Database, config.GlobalConfig.Database.User, config.GlobalConfig.Database.Password)
}
