package dataaccess

import (
	"github.com/google/wire"
	"goload/internal/dataaccess/cache"
	"goload/internal/dataaccess/database"
	"goload/internal/dataaccess/file"
	"goload/internal/dataaccess/mq"
)

var WireSet = wire.NewSet(
	cache.WireSet,
	database.WireSet,
	mq.WireSet,
	file.WireSet,
)