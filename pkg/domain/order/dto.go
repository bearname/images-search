package order

import "photofinish/pkg/common/util/uuid"

//type CreateOrderDTO struct {
//
//}

type StripeToken struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	Card   struct {
		ID                 string      `json:"id"`
		Object             string      `json:"object"`
		AddressCity        interface{} `json:"address_city"`
		AddressCountry     interface{} `json:"address_country"`
		AddressLine1       interface{} `json:"address_line1"`
		AddressLine1Check  interface{} `json:"address_line1_check"`
		AddressLine2       interface{} `json:"address_line2"`
		AddressState       interface{} `json:"address_state"`
		AddressZip         interface{} `json:"address_zip"`
		AddressZipCheck    interface{} `json:"address_zip_check"`
		Brand              string      `json:"brand"`
		Country            string      `json:"country"`
		CvcCheck           string      `json:"cvc_check"`
		DynamicLast4       interface{} `json:"dynamic_last4"`
		ExpMonth           int         `json:"exp_month"`
		ExpYear            int         `json:"exp_year"`
		Funding            string      `json:"funding"`
		Last4              string      `json:"last4"`
		Name               interface{} `json:"name"`
		TokenizationMethod interface{} `json:"tokenization_method"`
	} `json:"card"`
	ClientIP string `json:"client_ip"`
	Created  int    `json:"created"`
	Livemode bool   `json:"livemode"`
	Type     string `json:"type"`
	Used     bool   `json:"used"`
}

type CreateOrderDTO struct {
	ReceiptEmail string `json:"receiptEmail"`
	StripeToken  string `json:"stripeToken"`
	Data         []struct {
		PictureID string `json:"pictureId"`
	} `json:"data"`
	UserId     int
	TotalPrice int64
	OrderId    uuid.UUID
}

type UpdateOrderStatusDTO struct {
	OrderId string
	Status  PayStatus
}

type PayResultDTO struct {
	OrderId       string
	ID            string
	ReceiptURL    string
	ReceiptNumber string
	Status        string
	ConfirmUrl    string
}

type StripeError struct {
	Code      string `json:"code"`
	DocURL    string `json:"doc_url"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Param     string `json:"param"`
	RequestID string `json:"request_id"`
	Type      string `json:"type"`
}

type GetOrderDTO struct {
	OrderId string
	UserId  int
}

type ReturnOrderDTO struct {
	OrderId       string
	UserId        int
	Status        PayStatus
	TotalPrice    int
	PayId         string
	ReceiptUrl    string
	ReceiptNumber string
	CountPictures int
}
