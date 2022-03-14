package yookassa

import (
	"encoding/base64"
	"encoding/json"
	"github.com/col3name/images-search/pkg/common/util/uuid"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type SDK struct {
	shopId        int
	apiKey        string
	apiBaseUrl    string
	appUrl        string
	serviceDomain string
}

func NewYookassaSDK(shopId int, apiKey string, appUrl string, serviceDomain string) *SDK {
	s := new(SDK)
	s.shopId = shopId
	s.apiKey = apiKey
	s.apiBaseUrl = "https://api.yookassa.ru/v3"
	s.serviceDomain = serviceDomain
	s.appUrl = appUrl
	return s
}

func (s *SDK) Pay(payDTO *CreatePaymentDTO) (*PaymentResp, error) {
	paymentReq := PaymentReq{
		Amount: Amount{Value: strconv.Itoa(int(payDTO.Value)), Currency: string(payDTO.Currency)},
		PaymentMethod: PaymentMethodData{
			Type: "bank_card",
		},
		Confirmation: Confirmation{
			Type:      "redirect",
			ReturnURL: payDTO.ReturnUrl,
		},
		Metadata: payDTO.Meta,
	}
	marshal, err := json.Marshal(paymentReq)
	if err != nil {
		return nil, err
	}

	payload := strings.NewReader(string(marshal))
	res, err := s.doRequest("/payments", payload, s.getBasicAuthToken())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var paymentResp PaymentResp
	err = json.Unmarshal(body, &paymentResp)
	return &paymentResp, err
}

func (s *SDK) GetWebhookList() (*WebhookListResp, error) {
	req, err := http.NewRequest(http.MethodGet, s.apiBaseUrl+"/webhooks", nil)
	if err != nil {
		return nil, err
	}
	s.setAuthorization(req, s.getOAuth2ApiToken())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var webhooksResp WebhookListResp
	err = json.Unmarshal(all, &webhooksResp)
	if err != nil {
		return nil, err
	}
	return &webhooksResp, nil
}

func (s *SDK) CreateWebhook(webhook *BaseWebhook) (*Webhook, error) {
	data, err := json.Marshal(webhook)
	if err != nil {
		return nil, err
	}
	payload := strings.NewReader(string(data))

	resp, err := s.doRequest("/payments", payload, s.getOAuth2ApiToken())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Webhook
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *SDK) doRequest(action string, payload *strings.Reader, authorization string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, s.apiBaseUrl+action, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	idempotenceKey := uuid.Generate().String()
	req.Header.Add("Idempotence-Key", idempotenceKey)
	s.setAuthorization(req, authorization)
	res, err := http.DefaultClient.Do(req)
	return res, err
}

func (s *SDK) setAuthorization(req *http.Request, value string) {
	req.Header.Add("Authorization", value)
}

func (s *SDK) getOAuth2ApiToken() string {
	return "Bearer " + s.apiKey
}

func (s *SDK) getBasicAuthToken() string {
	return "Basic " + s.basicAuth(strconv.Itoa(s.shopId), s.apiKey)
}

func (s *SDK) basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
