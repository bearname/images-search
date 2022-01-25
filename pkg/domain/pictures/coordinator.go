package pictures

type CoordinatorService interface {
	PerformAddImage(picture Picture) error
}
