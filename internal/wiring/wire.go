package wiring

import (
	"github.com/google/wire"
	"goload/internal/configs"
	"goload/internal/dataAccess/database"
	"goload/internal/handler"
	"goload/internal/handler/grpc"
	"goload/internal/logic"
	"goload/internal/utils"
)

var WireSet = wire.NewSet(
	configs.WireSet,
	utils.WireSet,
	database.WireSet,
	logic.WireSet,
	handler.WireSet,
)

func initializeGRPCServer(configFilePath configs.ConfigFilePath) (grpc.Server, func(), error) {
	wire.Build(WireSet)
	return nil, nil, nil
}
