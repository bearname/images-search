package event

type Event struct {
	Id       int
	Name     string
	Location string
}

func NewEvent(id int, name string, location string) *Event {
	t := new(Event)
	t.Id = id
	t.Name = name
	t.Location = location
	return t
}
