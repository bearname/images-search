package main

import (
    "fmt"
    "github.com/stripe/stripe-go"
    "github.com/stripe/stripe-go/charge"
    "github.com/stripe/stripe-go/customer"
    "html/template"
    "net/http"
    "path/filepath"
)

func main() {
    publishableKey := "pk_test_RbEX8nfbG46rtdIFjveIM7SS00pLCspbDp"
    stripe.Key = "sk_test_KgVJC0IJtSVeZmmuuHu9aS0f005iVmCFmi"

    tmpls, _ := template.ParseFiles(filepath.Join("templates", "index.html"))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl := tmpls.Lookup("index.html")
        tmpl.Execute(w, map[string]string{"Key": publishableKey})
    })

    http.HandleFunc("/charge", func(w http.ResponseWriter, r *http.Request) {
        r.ParseForm()
        email := r.Form.Get("stripeEmail")

        customerParams := &stripe.CustomerParams{Email: &email}
        customerParams.SetSource(r.Form.Get("stripeToken"))
        newCustomer, err := customer.New(customerParams)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        var amount int64
        amount = 500
        value := "usd"
        id := newCustomer.ID
        disc := "Sample Charge"
        chargeParams := &stripe.ChargeParams{
            Amount:      &amount,
            Currency:    &value,
            Description: &disc,
            Customer:    &id,
        }

        if _, err := charge.New(chargeParams); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        fmt.Fprintf(w, "Charge completed successfully!")
    })
    log.Println("checking .....start")
    http.ListenAndServe(":4567", nil)
}