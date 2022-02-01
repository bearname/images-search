package yookassa

import "time"

type PaymentResp struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Recipient struct {
		AccountID string `json:"account_id"`
		GatewayID string `json:"gateway_id"`
	} `json:"recipient"`
	PaymentMethod struct {
		Type  string `json:"type"`
		ID    string `json:"id"`
		Saved bool   `json:"saved"`
	} `json:"payment_method"`
	CreatedAt    time.Time `json:"created_at"`
	Confirmation struct {
		Type            string `json:"type"`
		ReturnURL       string `json:"return_url"`
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation"`
	Test       bool `json:"test"`
	Paid       bool `json:"paid"`
	Refundable bool `json:"refundable"`
	Metadata   struct {
	} `json:"metadata"`
}

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}
type Confirmation struct {
	Type      string `json:"type"`
	ReturnURL string `json:"return_url"`
}

type MetadataItem struct {
	Key   string
	Value string
}

type PaymentMethodData struct {
	Type string `json:"type"`
}
type PaymentReq struct {
	Amount        Amount            `json:"amount"`
	PaymentMethod PaymentMethodData `json:"payment_method_data"`
	Confirmation  Confirmation      `json:"confirmation"`
	Description   string            `json:"description"`
	Metadata      map[string]string `json:"metadata"`
}

type Webhook struct {
	BaseWebhook
	ID string `json:"id"`
}

type BaseWebhook struct {
	Event Event  `json:"event"`
	URL   string `json:"url"`
}

type WebhookListResp struct {
	Type  string    `json:"type"`
	Items []Webhook `json:"items"`
}

type Currency string

const (
	USD Currency = "USD"
	RUB Currency = "RUB"
)

type Event string

const (
	PaymentSucceeded         Event = "payment.succeeded"
	PaymentCanceled          Event = "payment.canceled"
	RefundSucceeded          Event = "refund.succeeded"
	PaymentWaitingForCapture Event = "payment.waiting_for_capture"
)

type CreatePaymentDTO struct {
	Value     int64
	Currency  Currency
	OrderId   string
	Meta      map[string]string
	ReturnUrl string
}

type Card struct {
	First6        string `json:"first6"`
	Last4         string `json:"last4"`
	ExpiryMonth   string `json:"expiry_month"`
	ExpiryYear    string `json:"expiry_year"`
	CardType      string `json:"card_type"`
	IssuerCountry string `json:"issuer_country"`
	IssuerName    string `json:"issuer_name"`
}

type PaymentMethod struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Saved bool   `json:"saved"`
	Card  Card   `json:"card"`
	Title string `json:"title"`
}

type AuthorizationDetails struct {
	Rrn          string `json:"rrn"`
	AuthCode     string `json:"auth_code"`
	ThreeDSecure struct {
		Applied bool `json:"applied"`
	} `json:"three_d_secure"`
}

type NotificationObject struct {
	ID                   string               `json:"id"`
	Status               string               `json:"status"`
	Paid                 bool                 `json:"paid"`
	Amount               Amount               `json:"amount"`
	AuthorizationDetails AuthorizationDetails `json:"authorization_details"`
	CreatedAt            time.Time            `json:"created_at"`
	Description          string               `json:"description"`
	ExpiresAt            time.Time            `json:"expires_at"`
	Metadata             map[string]string    `json:"metadata"`
	PaymentMethod        PaymentMethod        `json:"payment_method"`
	Refundable           bool                 `json:"refundable"`
	Test                 bool                 `json:"test"`
}

type NotificationEvent struct {
	Type   string             `json:"type"`
	Event  string             `json:"event"`
	Object NotificationObject `json:"object"`
}
