package mq

import (
	"github.com/google/wire"
	"goload/internal/dataAccess/mq/consumer"
	"goload/internal/dataAccess/mq/producer"
)

var WireSet = wire.NewSet(
	consumer.WireSet,
	producer.WireSet,
)
