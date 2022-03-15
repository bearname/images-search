package paySystem

import (
	"errors"
	"github.com/col3name/images-search/pkg/domain/order"
	yookassa2 "github.com/col3name/images-search/pkg/infrastructure/external/yookassa"
	log "github.com/sirupsen/logrus"
	"net"
	"strings"
)

var ErrInvalidIP = errors.New("unsupported ip address")
var ErrInvalidEvent = errors.New("invalid event")
var ErrFailedHandle = errors.New("failed handle payment")

type YookassaService struct {
	sdk             *yookassa2.SDK
	appUrl          string
	successWebhook  yookassa2.Webhook
	canceledWebhook yookassa2.Webhook
}

func NewYookassaService(shopId int, apiKey string) (*YookassaService, error) {
	s := new(YookassaService)
	s.appUrl = "https://adb3-188-187-176-103.ngrok.io"
	s.sdk = yookassa2.NewYookassaSDK(shopId, apiKey, s.appUrl, "https://yookassa.ru")
	err := s.initWebhook()
	if err != nil {
		return s, err
	}

	return s, nil
}

func (s *YookassaService) checkIP(requestIp string) bool {
	validIPList := []string{
		"185.71.76.0/27",
		"185.71.77.0/27",
		"77.75.153.0/25",
		"77.75.156.11",
		"77.75.156.35",
		"77.75.154.128/25",
		"2a02:5180:0:1509::/64",
		"2a02:5180:0:2655::/64",
		"2a02:5180:0:1533::/64",
		"2a02:5180:0:2669::/64",
	}
	isFound := false
	ipB := net.ParseIP(requestIp)

	for _, validIpRange := range validIPList {
		if s.cidrRangeContains(validIpRange, ipB) {
			isFound = true
			break
		}
	}
	return isFound
}

func (s *YookassaService) cidrRangeContains(validIpRange string, ipB net.IP) bool {
	_, ipnetA, err := net.ParseCIDR(validIpRange)
	return err == nil && ipnetA.Contains(ipB)
}

func (s *YookassaService) OnHandleEvent(event interface{}, remoteIp string) (*order.UpdateOrderStatusDTO, error) {
	if !s.checkIP(remoteIp) {
		return nil, ErrInvalidIP
	}
	var yookassaEvent yookassa2.NotificationEvent
	switch event.(type) {
	case yookassa2.NotificationEvent:
		yookassaEvent = event.(yookassa2.NotificationEvent)
	default:
		return nil, ErrInvalidEvent
	}

	updateOrderDTO, isOk := s.handlePayment(yookassaEvent)
	if !isOk {
		log.Println(ErrFailedHandle.Error(), yookassaEvent)
		return nil, ErrFailedHandle
	}
	return updateOrderDTO, nil
}

func (s *YookassaService) Buy(dto *order.CreateOrderDTO) (*order.PayResultDTO, error) {
	meta := map[string]string{}
	meta["orderId"] = dto.OrderId.String()

	paymentResp, err := s.sdk.Pay(&yookassa2.CreatePaymentDTO{
		OrderId:   dto.OrderId.String(),
		Currency:  yookassa2.USD,
		Value:     dto.TotalPrice,
		ReturnUrl: s.appUrl + "/api/v1/yookassa",
	})

	if err != nil {
		return nil, err
	}
	return &order.PayResultDTO{
		ID:         paymentResp.ID,
		Status:     paymentResp.Status,
		ConfirmUrl: paymentResp.Confirmation.ConfirmationURL,
	}, nil
}

func (s *YookassaService) initWebhook() error {
	list, err := s.sdk.GetWebhookList()
	if err != nil {
		return err
	}
	if len(list.Items) == 0 {
		webhook, err := s.createWebhook(yookassa2.PaymentSucceeded)
		if err != nil {
			return err
		}
		s.successWebhook = *webhook

		webhook, err = s.createWebhook(yookassa2.PaymentCanceled)
		if err != nil {
			return err
		}
		s.canceledWebhook = *webhook
	}

	isFoundSucceed, isFoundCanceled := s.checkExistedWebhook(list)
	if !isFoundSucceed {
		webhook, err := s.createWebhook(yookassa2.PaymentSucceeded)
		if err != nil {
			return err
		}
		s.successWebhook = *webhook
	}

	if !isFoundCanceled {
		webhook, err := s.createWebhook(yookassa2.PaymentCanceled)
		if err != nil {
			return err
		}
		s.canceledWebhook = *webhook
	}
	return nil
}

func (s *YookassaService) checkExistedWebhook(list *yookassa2.WebhookListResp) (bool, bool) {
	isFoundSucceed := false
	isFoundCanceled := false
	for _, item := range list.Items {
		if strings.Contains(item.URL, s.appUrl+"/api/v1/yookassa") {
			switch item.Event {
			case yookassa2.PaymentCanceled:
				isFoundCanceled = true
			case yookassa2.PaymentSucceeded:
				isFoundSucceed = true
			}
		}
		if isFoundCanceled && isFoundSucceed {
			break
		}
	}
	return isFoundSucceed, isFoundCanceled
}

func (s *YookassaService) createWebhook(event yookassa2.Event) (*yookassa2.Webhook, error) {
	webhook, err := s.sdk.CreateWebhook(&yookassa2.BaseWebhook{
		Event: event,
		URL:   s.appUrl + "/api/v1/yookassa",
	})

	return webhook, err
}

func (s *YookassaService) handlePayment(event yookassa2.NotificationEvent) (*order.UpdateOrderStatusDTO, bool) {
	var updateOrderDTO order.UpdateOrderStatusDTO
	metadata := event.Object.Metadata
	orderId, ok := metadata["orderId"]
	if !ok {
		return nil, false
	}
	updateOrderDTO.OrderId = orderId
	switch event.Event {
	case string(yookassa2.PaymentSucceeded):
		updateOrderDTO.Status = order.PaySuccess
	case string(yookassa2.PaymentCanceled):
		updateOrderDTO.Status = order.PayCanceled
	default:
		return nil, false
	}
	log.Println("Order: ", updateOrderDTO, event.Type, " was successful!")

	return &updateOrderDTO, true
}
