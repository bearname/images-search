package main

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "photofinish/pkg/common/infrarstructure/server"
    "photofinish/pkg/common/util/uuid"
    "strconv"
    "time"

    "github.com/joho/godotenv"
    "github.com/stripe/stripe-go"
    "github.com/stripe/stripe-go/charge"
)

type CreateOrderDto struct {
    ReceiptEmail string `json:receiptEmail`
    StripeToken  StripeToken `json:stripeToken`
    Data         []struct {
        PictureID string `json:"pictureId"`
    } `json:"data"`
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

func cmain() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    router := mux.NewRouter()
    apiV1Route := router.PathPrefix("/api/v1").Subrouter()

    apiV1Route.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{
            "message": "Hello, World!",
        }`))
    }).Methods(http.MethodGet)
    router.HandleFunc("/webhook", func(w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        if (*req).Method == "OPTIONS" {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        const MaxBodyBytes = int64(65536)
        req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
        payload, err := ioutil.ReadAll(req.Body)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
            w.WriteHeader(http.StatusServiceUnavailable)
            return
        }

        event := stripe.Event{}

        if err := json.Unmarshal(payload, &event); err != nil {
            fmt.Fprintf(os.Stderr, "Failed to parse webhook body json: %v\n", err.Error())
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        // Unmarshal the event data into an appropriate struct depending on its Type
        switch event.Type {
        case "payment_intent.succeeded":
            var paymentIntent stripe.PaymentIntent
            err := json.Unmarshal(event.Data.Raw, &paymentIntent)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            fmt.Println("PaymentIntent was successful!")
        case "payment_method.attached":
            var paymentMethod stripe.PaymentMethod
            err := json.Unmarshal(event.Data.Raw, &paymentMethod)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            fmt.Println("PaymentMethod was attached to a Customer!")
        // ... handle other event types
        case "charge.succeeded", "charge.failed":
            if ok := !handlePayment(w, event); ok {
                fmt.Println(event.Type + " \npayment: " + strconv.FormatBool(ok))
            }
        default:
            fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
        }

        w.WriteHeader(http.StatusOK)
    }).Methods(http.MethodPost)

    apiV1Route.HandleFunc("/charges", func(w http.ResponseWriter, req *http.Request) {

        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        if (*req).Method == "OPTIONS" {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        var createOrderDto CreateOrderDto
        all, err := ioutil.ReadAll(req.Body)
        err = json.Unmarshal(all, &createOrderDto)
        if err != nil {
            fmt.Println(err.Error())
            w.Write([]byte(err.Error()))
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        stripe.Key = "sk_test_51KKrZ9Lmz0KHDy38ouJgOfCBDvnXMEOzJd4yh2EXwpFpOCUgCCYO9iNhSakqTxRgKDrgmfQyp8NpiFzUZXZy1wp200u37l9R7f"
        // Attempt to make the charge.
        // We are setting the charge response to _
        // as we are not using it.
        token := stripe.String("tok_visa")
        //token := paySystem.String(createOrderDto.StripeToken.ID)
        params := stripe.ChargeParams{
            Amount:       stripe.Int64(int64(len(createOrderDto.Data) * 100)),
            Currency:     stripe.String(string(stripe.CurrencyUSD)),
            Source:       &stripe.SourceParams{Token: token}, // this should come from clientside
            ReceiptEmail: stripe.String(createOrderDto.ReceiptEmail),
        }
        //TODO
        // type Order struct {
        //    id uuid
        //    pictures []Picture
        //    status OrderStatus
        //    userId uuid
        // }
        // createOrder(uuid)
        orderId := uuid.Generate()
        params.AddMetadata("order", string(all))
        params.AddMetadata("orderId", orderId.String())
        resp, err := charge.New(&params)
        if err != nil {
            // Handle any domainerror from attempt to charge
            var stripeErr StripeError
            err := json.Unmarshal([]byte(err.Error()), &stripeErr)
            if err != nil {
                fmt.Println(err)
            } else {
                fmt.Println(stripeErr)
            }
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        fmt.Println(resp.FailureMessage)
        fmt.Println(resp.FailureCode)

        fmt.Println(resp)

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("success"))
    }).Methods(http.MethodPost)

    port := "8000"

    log.Println("Start on port '" + port + " 'at " + time.Now().String())
    httpServer := server.HttpServer{}
    killSignalChan := httpServer.GetKillSignalChan()
    srv := httpServer.StartServer(port, router)
    httpServer.WaitForKillSignal(killSignalChan)
    err = srv.Shutdown(context.TODO())

    log.Println("Stop at " + time.Now().String())

    if err != nil {
        fmt.Println(err)
        return
    }
}

func handlePayment(w http.ResponseWriter, event stripe.Event) bool {
    var paymentIntent stripe.PaymentIntent
    err := json.Unmarshal(event.Data.Raw, &paymentIntent)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
        w.WriteHeader(http.StatusBadRequest)
        return false
    }
    order := paymentIntent.Metadata["order"]
    var chargeJson CreateOrderDto
    err = json.Unmarshal([]byte(order), &chargeJson)
    if err != nil {
        fmt.Println(err.Error())
        w.Write([]byte(err.Error()))
        w.WriteHeader(http.StatusBadRequest)
        return false
    }
    orderId := paymentIntent.Metadata["orderId"]

    var orderStatus int
    switch event.Type {

    case "charge.succeeded":
        orderStatus = 1
    case "charge.failed":
        orderStatus = 2
    default:
        return false
    }

    fmt.Println("OrderId: " + orderId)
    fmt.Println("OrderStatus: " + strconv.Itoa(orderStatus))
    fmt.Println(event.Type + " was successful!")
    //TODO
    // updateOrderStatus(orderId, status)
    // get /api/v1/users/{username}/orders
    // get /api/v1/users/{username}/orders/{orderId}
    return true
}
