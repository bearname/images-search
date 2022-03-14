package broker

import (
	"github.com/col3name/images-search/pkg/common/util/uuid"
	"time"
)

type ProcessingStatus int

const (
	OutboxDone ProcessingStatus = iota
	OutboxProcessing
	OutboxNotProcessing
)

type Outbox struct {
	Id          uuid.UUID
	BrokerTopic string
	BrokerKey   string
	BrokerValue string
	Status      ProcessingStatus
	UpdatedAt   time.Time
}
