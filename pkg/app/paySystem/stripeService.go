package paySystem

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"photofinish/pkg/domain/order"
	"reflect"
)

type StripeService struct {
	secretKey string
}

func NewStripeService(secretKey string) *StripeService {
	s := new(StripeService)
	s.secretKey = secretKey
	return s
}

func (s *StripeService) OnHandleEvent(event interface{}, _ string) (*order.UpdateOrderStatusDTO, error) {
	var stripeEvent stripe.Event
	of := reflect.TypeOf(event)
	fmt.Println(of)
	switch event.(type) {
	case stripe.Event:
		stripeEvent = event.(stripe.Event)
	default:
		return nil, errors.New("invalid event")
	}

	updateOrderDTO, isOk := s.handlePayment(stripeEvent)
	if !isOk {
		return nil, errors.New("failed handle payment" + stripeEvent.ID)
	}
	return updateOrderDTO, nil
}

func (s *StripeService) Buy(dto order.CreateOrderDTO) (*order.PayResultDTO, error) {
	stripe.Key = s.secretKey
	// Attempt to make the charge.
	// We are setting the charge response to _
	// as we are not using it.
	token := stripe.String("tok_visa")
	params := stripe.ChargeParams{
		Amount:       stripe.Int64(dto.TotalPrice),
		Currency:     stripe.String(string(stripe.CurrencyUSD)),
		Source:       &stripe.SourceParams{Token: token}, // this should come from clientside
		ReceiptEmail: stripe.String(dto.ReceiptEmail),
	}
	data, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	orderIdStr := dto.OrderId.String()
	params.AddMetadata("order", string(data))
	params.AddMetadata("orderId", orderIdStr)
	resp, err := charge.New(&params)
	if err != nil {
		var stripeErr order.StripeError
		err = json.Unmarshal([]byte(err.Error()), &stripeErr)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(stripeErr)
		}
		return nil, err
	}

	fmt.Println(resp.FailureMessage)
	fmt.Println(resp.FailureCode)

	fmt.Println(resp)

	return &order.PayResultDTO{
		OrderId:       orderIdStr,
		ID:            resp.ID,
		ReceiptURL:    resp.ReceiptURL,
		ReceiptNumber: resp.ReceiptNumber,
		Status:        resp.Status,
	}, nil
}

func (s *StripeService) handlePayment(event stripe.Event) (*order.UpdateOrderStatusDTO, bool) {
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		return nil, false
	}
	payData := paymentIntent.Metadata["order"]
	var chargeJson order.CreateOrderDTO
	err = json.Unmarshal([]byte(payData), &chargeJson)
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	orderId := paymentIntent.Metadata["orderId"]

	var updateOrderDTO order.UpdateOrderStatusDTO
	updateOrderDTO.OrderId = orderId
	switch event.Type {
	case "charge.succeeded":
		updateOrderDTO.Status = order.PaySuccess
	case "charge.failed":
		updateOrderDTO.Status = order.PayFailed
	default:
		return nil, false
	}
	fmt.Println("Order: ", updateOrderDTO)
	fmt.Println(event.Type + " was successful!")
	return &updateOrderDTO, true
}
