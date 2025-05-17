package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/matiuskm/go-commerce/db"
	"github.com/matiuskm/go-commerce/models"
)

// GET /my/addresses
func ListAddresses(c *gin.Context) {
	uid := c.MustGet("user_id").(uint)
	var addrs []models.Address
	db.DB.Where("user_id = ?", uid).Find(&addrs)

	response := []models.AddressResponse{}
	for _, addr := range addrs {
	  response = append(response, models.AddressResponse{
		ID: addr.ID,
		Label: addr.Label,
		Phone: addr.Phone,
		Street: addr.Street,
		RecipientName: addr.RecipientName,
	  })
	}

	c.JSON(200, gin.H{"addresses": response})
  }
  
  // POST /my/addresses
  func CreateAddress(c *gin.Context) {
	uid := c.MustGet("user_id").(uint)
	var in models.Address
	if err := c.ShouldBindJSON(&in); err != nil {
	  c.JSON(400, gin.H{"error": err.Error()}); return
	}
	in.UserID = uid
	db.DB.Create(&in)

	response := models.AddressResponse{
	  ID: in.ID,
	  Label: in.Label,
	  Phone: in.Phone,
	  Street: in.Street,
	  RecipientName: in.RecipientName,
	}

	c.JSON(201, gin.H{"address": response})
  }
  
  // PUT /my/addresses/:id
  func UpdateAddress(c *gin.Context) {
	uid := c.MustGet("user_id").(uint)
	id := c.Param("id")
	var addr models.Address
	if err := db.DB.First(&addr, "id = ? AND user_id = ?", id, uid).Error; err != nil {
	  c.JSON(404, gin.H{"error": "not found"}); return
	}
	if err := c.ShouldBindJSON(&addr); err != nil {
	  c.JSON(400, gin.H{"error": err.Error()}); return
	}
	db.DB.Save(&addr)

	response := models.AddressResponse{
	  ID: addr.ID,
	  Label: addr.Label,
	  Phone: addr.Phone,
	  Street: addr.Street,
	  RecipientName: addr.RecipientName,
	}

	c.JSON(200, gin.H{"address": response})
  }
  
  // DELETE /my/addresses/:id
  func DeleteAddress(c *gin.Context) {
	uid := c.MustGet("user_id").(uint)
	id := c.Param("id")
	if err := db.DB.Delete(&models.Address{}, "id = ? AND user_id = ?", id, uid).Error; err != nil {
	  c.JSON(500, gin.H{"error": "delete failed"}); return
	}

	c.Status(204)
  }
  