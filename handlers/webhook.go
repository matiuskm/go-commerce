package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/matiuskm/go-commerce/db"
	"github.com/matiuskm/go-commerce/models"
)

// Xendit will POST JSON like:
//
//	{
//	  "id": "inv_123",
//	  "external_id": "ORD-20250519-...",
//	  "status": "PAID",
//	  ...other fields...
//	}
type XenditWebhookPayload struct {
	ID         string `json:"id"`
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
}

func XenditWebhookHandler(c *gin.Context) {
	token := c.GetHeader("X-Callback-Token")
	if token != os.Getenv("XENDIT_CALLBACK_TOKEN") {
		c.Status(http.StatusUnauthorized)
		return
	}

	var p XenditWebhookPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// find the order by the invoice ID Xendit sent
	var order models.Order
	if err := db.DB.
		Where("xendit_invoice = ?", p.ID).
		First(&order).Error; err != nil {
		// we don’t want to expose DB errors here—just return OK so Xendit doesn’t retry forever
		c.Status(http.StatusOK)
		return
	}

	// update your order’s status
	order.Status = strings.ToLower(p.Status) // e.g. "PAID" → "paid"
	db.DB.Save(&order)

	c.Status(http.StatusOK)
}
