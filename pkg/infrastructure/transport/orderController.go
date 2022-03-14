package transport

import (
	"encoding/json"
	"errors"
	"github.com/col3name/images-search/pkg/common/util"
	"github.com/col3name/images-search/pkg/domain/order"
	"github.com/col3name/images-search/pkg/domain/user"
	paySystem2 "github.com/col3name/images-search/pkg/infrastructure/paySystem"
	"github.com/col3name/images-search/pkg/infrastructure/yookassa"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"io/ioutil"
	"net/http"
)

type OrderController struct {
	BaseController
	orderService    order.Service
	userService     user.Service
	stripeService   *paySystem2.StripeService
	yookassaService *paySystem2.YookassaService
}

func NewOrderController(userService user.Service,
	service order.Service,
	paySystem *paySystem2.StripeService,
	yookassaService *paySystem2.YookassaService) *OrderController {
	c := new(OrderController)
	c.userService = userService
	c.orderService = service
	c.stripeService = paySystem
	c.yookassaService = yookassaService
	return c
}

func (c *OrderController) Pay() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*req).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		var createOrderDto order.CreateOrderDTO
		all, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(all, &createOrderDto)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		query := req.URL.Query()
		paySystemReq := query.Get("paySystem")

		var respCode int
		var orderID string
		switch paySystemReq {
		case "stripe":
			orderID, err = c.orderService.Buy(c.stripeService, createOrderDto)
			if err != nil {
				respCode = http.StatusInternalServerError
			}
		case "yookassa":
			orderID, err = c.orderService.Buy(c.yookassaService, createOrderDto)
			if err != nil {
				respCode = http.StatusInternalServerError
			}
		default:
			err = errors.New("paySystemReq := stripe|yookassa")
			respCode = http.StatusBadRequest
		}
		if err != nil {
			log.Println(err)
			w.WriteHeader(respCode)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success, orderId" + orderID))
	}
}

func (c *OrderController) OnEventYookassa() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*req).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		remoteIp, err := util.GetRemoteIp(req.RemoteAddr)
		if err != nil {
			log.Printf("Error reading request body: %v\n", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		const MaxBodyBytes = int64(65536)
		req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
		payload, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("Error reading request body: %v\n", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		event := yookassa.NotificationEvent{}

		if err = json.Unmarshal(payload, &event); err != nil {
			log.Printf("PayFailed to parse webhook body json: %v\n", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		code := http.StatusOK
		if _, err = c.yookassaService.OnHandleEvent(event, remoteIp); err != nil {
			code = http.StatusBadRequest
		}

		w.WriteHeader(code)
	}
}

func (c *OrderController) OnEventStripe() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*req).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		remoteIp, err := util.GetRemoteIp(req.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		const MaxBodyBytes = int64(65536)
		req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
		payload, err := ioutil.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		event := stripe.Event{}

		if err = json.Unmarshal(payload, &event); err != nil {
			log.Printf("PayFailed to parse webhook body json: %v\n", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = c.orderService.OnHandle(c.stripeService, event, remoteIp)
		if err != nil {
			log.Printf("PayFailed to parse webhook body json: %v\n", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (c *OrderController) GetOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*req).Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		vars := mux.Vars(req)
		orderId := vars["id"]

		if len(orderId) == 0 || !util.IsUUID(orderId) {
			msg := "Invalid 'id' query parameter. 'id' must be uuid"
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		username, ok := context.Get(req, "username").(string)
		if !ok {
			context.Clear(req)
			http.Error(w, "cannot check username", http.StatusBadRequest)
			return
		}
		userDto, err := c.userService.Find(username)
		if err != nil {
			log.Error(err)
			http.Error(w, "cannot check username", http.StatusBadRequest)
			return
		}
		returnOrderDto, err := c.orderService.GetOrder(&order.GetOrderDTO{
			OrderId: orderId,
			UserId:  userDto.Id,
		})
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c.WriteJsonResponse(w, returnOrderDto)
	}
}
