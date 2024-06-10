package database

import "github.com/google/wire"

var WireSet = wire.NewSet(
	InitializeAndMigrateUpDB,
	InitializeGoquDB,
	NewAccountDataAccessor,
	NewMigrator,
	NewAccountPasswordDataAccessor,
	NewDownloadTaskDataAccessor,
	NewTokenPublicKeyDataAccessor,
)
