package dataaccess

import (
	"github.com/google/wire"
	"goload/internal/dataAccess/cache"
	"goload/internal/dataAccess/database"
	"goload/internal/dataAccess/mq"
)

var WireSet = wire.NewSet(
	cache.WireSet,
	database.WireSet,
	mq.WireSet,
)
