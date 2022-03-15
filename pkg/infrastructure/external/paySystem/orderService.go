package paySystem

import (
	"github.com/col3name/images-search/pkg/common/util/uuid"
	"github.com/col3name/images-search/pkg/domain/order"
	log "github.com/sirupsen/logrus"
)

type OrderService struct {
	repo order.Repo
}

func NewOrderService(repo order.Repo) *OrderService {
	s := new(OrderService)
	s.repo = repo
	return s
}

func (s *OrderService) OnHandle(paymentSystem order.PaymentSystem, event interface{}, remoteIp string) error {
	updateOrderDTO, err := paymentSystem.OnHandleEvent(event, remoteIp)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = s.repo.UpdateStatus(updateOrderDTO)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	//TODO
	// updateOrderStatus(orderId, status)
	// get /api/v1/users/{username}/orders
	// get /api/v1/users/{username}/orders/{orderId}
	return nil
}

func (s *OrderService) GetOrder(dto *order.GetOrderDTO) (*order.ReturnOrderDTO, error) {
	return s.repo.GetOrder(dto)
}
func (s *OrderService) Buy(paymentSystem order.PaymentSystem, dto *order.CreateOrderDTO) (string, error) {
	totalPrice := int64(len(dto.Data) * 100)
	dto.TotalPrice = totalPrice

	orderId := uuid.Generate()
	dto.OrderId = orderId
	orderIdStr := orderId.String()

	err := s.repo.Store(dto)
	if err != nil {
		return "", err
	}

	payResult, err := paymentSystem.Buy(dto)
	if err != nil {
		statusDTO := order.UpdateOrderStatusDTO{
			OrderId: dto.OrderId.String(),
			Status:  order.PayFailed,
		}
		err = s.repo.UpdateStatus(&statusDTO)
		if err != nil {
			return orderIdStr, err
		}
		return "", err
	}

	err = s.repo.SavePayResult(payResult)
	if err != nil {
		return orderIdStr, err
	}

	return orderIdStr, nil
}
