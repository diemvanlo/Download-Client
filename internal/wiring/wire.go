//go:build wireinject
// +build wireinject

//
//go:generate go run github.com/google/wire/cmd/wire
package wiring

import (
	"github.com/google/wire"
	"goload/internal/app"
	"goload/internal/configs"
	dataaccess "goload/internal/dataAccess"
	"goload/internal/handler"
	"goload/internal/logic"
	"goload/internal/utils"
)

var WireSet = wire.NewSet(
	configs.WireSet,
	utils.WireSet,
	dataaccess.WireSet,
	logic.WireSet,
	handler.WireSet,
	app.WireSet,
)

func InitializeGRPCServer(configFilePath configs.ConfigFilePath) (*app.Server, func(), error) {
	wire.Build(WireSet)
	return nil, nil, nil
}
