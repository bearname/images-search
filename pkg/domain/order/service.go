package order

type Service interface {
	Buy(paymentSystem PaymentSystem, dto CreateOrderDTO) (string, error)
	OnHandle(paymentSystem PaymentSystem, event interface{}, remoteIp string) error
	GetOrder(dto *GetOrderDTO) (*ReturnOrderDTO, error)
}

type PaymentSystem interface {
	OnHandleEvent(event interface{}, remoteIp string) (*UpdateOrderStatusDTO, error)
	Buy(dto CreateOrderDTO) (*PayResultDTO, error)
}
