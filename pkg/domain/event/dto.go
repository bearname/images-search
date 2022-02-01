package event

type CreateEventInputDto struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Date     string `json:"date"`
}
