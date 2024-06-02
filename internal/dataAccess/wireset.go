package dataAccess

import (
	"github.com/google/wire"
	"goload/internal/dataAccess/cache"
	"goload/internal/dataAccess/database"
)

var WireSet = wire.NewSet(
	cache.WireSet,
	database.WireSet,
)
