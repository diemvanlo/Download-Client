package wiring

import (
	"github.com/google/wire"
	"goload/internal/configs"
	"goload/internal/dataAccess/database"
	"goload/internal/handler"
	"goload/internal/handler/grpc"
	"goload/internal/logic"
)

var WireSet = wire.NewSet(
	configs.WireSet,
	database.WireSet,
	logic.WireSet,
	handler.WireSet,
)

func initializeGRPCServer(configFilePath configs.ConfigFilePath) (grpc.Server, func(), error) {
	wire.Build(WireSet)
	return nil, nil, nil
}
