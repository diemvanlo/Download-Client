package wiring

import (
	"github.com/google/wire"
	"goload/internal/app"
	"goload/internal/configs"
	"goload/internal/dataAccess/cache"
	"goload/internal/dataAccess/database"
	"goload/internal/handler"
	"goload/internal/logic"
	"goload/internal/utils"
)

var WireSet = wire.NewSet(
	configs.WireSet,
	utils.WireSet,
	database.WireSet,
	cache.WireSet,
	logic.WireSet,
	handler.WireSet,
	app.WireSet,
)

func InitializeGRPCServer(configFilePath configs.ConfigFilePath) (*app.Server, func(), error) {
	wire.Build(WireSet)
	return nil, nil, nil
}
