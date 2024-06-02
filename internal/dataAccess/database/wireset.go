package database

import "github.com/google/wire"

var WireSet = wire.NewSet(
	InitializeAndMigrateUpDB,
	InitializeGoquDB,
	NewDatabaseAccessor,
	NewMigrator,
	NewAccountPasswordDataAccessor,
	NewDownloadTaskDataAccessor,
	NewTokenPublicKeyDataAccessor,
)
