package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/matiuskm/go-commerce/models"
	xendit "github.com/xendit/xendit-go/v6"
	invoice "github.com/xendit/xendit-go/v6/invoice"
	"gorm.io/gorm"
)

func CreateXenditInvoice(tx *gorm.DB, order *models.Order, customerEmail string) error {
	log.Println("Creating Xendit Invoice")
	log.Printf("▶️  CreateXenditInvoice for OrderNum=%s, Amount=%d\n",
        order.OrderNum, order.Total)
	createInvoiceRequest := *invoice.NewCreateInvoiceRequest(order.OrderNum, float64(order.Total))
	createInvoiceRequest.SetCurrency("IDR")
	customer := invoice.NewCustomerObject()
	customer.SetEmail(customerEmail)
	createInvoiceRequest.SetCustomer(*customer)
	createInvoiceRequest.SetDescription(fmt.Sprintf("Order %s payment", order.OrderNum))
	createInvoiceRequest.SetSuccessRedirectUrl(fmt.Sprintf("%s?external_id=%s",os.Getenv("XENDIT_SUCCESS_URL"), order.OrderNum))
	createInvoiceRequest.SetFailureRedirectUrl(os.Getenv("XENDIT_FAILURE_URL"))
	createInvoiceRequest.SetPaymentMethods([]string{"BRI", "BNI", "MANDIRI", "PERMATA", "QRIS"})

	client := xendit.NewClient(os.Getenv("XENDIT_SECRET_KEY"))

	resp, r, err := client.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(createInvoiceRequest).
		Execute()
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InvoiceApi.CreateInvoice``: %v\n", err.Error())

        b, _ := json.Marshal(err.FullError())
        fmt.Fprintf(os.Stderr, "Full Error Struct: %v\n", string(b))

        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		return fmt.Errorf("xendit.CreateInvoice: %w", err)
	}
	log.Println(*resp.Id)
	log.Println(resp.InvoiceUrl)
	if resp.Id != nil {
		log.Println("Set XenditInvoice")
        order.XenditInvoice = *resp.Id
    }
	log.Println("Set XenditUrl")
	order.XenditUrl = resp.InvoiceUrl
	log.Println("Save Order")
    return tx.Save(order).Error
}