package order

type Repo interface {
	Store(order *CreateOrderDTO) error
	UpdateStatus(order *UpdateOrderStatusDTO) error
	SavePayResult(result *PayResultDTO) error
	GetOrder(dto *GetOrderDTO) (*ReturnOrderDTO, error)
}
