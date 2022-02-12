package broker

//CheckExist(outboxId string) error
//Store(imageTextDetectionDto *Outbox) (int, error)
type Repo interface {
	//Store(outbox *Outbox) error
	UpdateStatus(outboxId string, status ProcessingStatus) error
	CheckExist(outboxId string) error
	FindNotCompletedOutboxList(limit int) (*[]Outbox, error)
	//Delete(outboxId string) error
}
