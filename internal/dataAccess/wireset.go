package dataAccess

import (
	"github.com/google/wire"
	"goload/internal/dataAccess/database"
)

var WireSet = wire.NewSet(
	database.WireSet)
