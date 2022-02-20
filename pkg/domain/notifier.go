package domain

type Message struct {
	Message string
}

type NotifierService interface {
	Notify(msg Message) error
}
