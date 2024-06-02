package database

import "github.com/google/wire"

var WireSet = wire.NewSet(
	InitializeDB,
	InitializeGoquDB,
	NewDatabaseAccessor,
	NewMigrator,
	NewAccountPasswordDataAccessor,
	NewDownloadTaskDataAccessor,
	NewTokenPublicKeyDataAccessor,
)
