package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matiuskm/go-commerce/config"
	"github.com/matiuskm/go-commerce/db"

	"github.com/matiuskm/go-commerce/handlers"
	"github.com/matiuskm/go-commerce/middlewares"
)

func main() {
	config.LoadEnv()

	db.Init()

	r := gin.Default()

	// Public routes
	r.GET("/", handlers.HomeHandler)
	r.POST("/auth/register", handlers.RegisterHandler)
	r.POST("/auth/login", handlers.LoginHandler)

	r.GET("/products", handlers.GetAllProductsHandler)
	r.GET("/products/:id", handlers.GetProductByIDHandler)

	// Protected routes
	auth := r.Group("/")
	auth.Use(middlewares.AuthMiddleware()) 
	{
		// User routes
		auth.GET("/my/profile", handlers.GetUserProfileHandler)
		auth.PATCH("/my/profile", handlers.UpdateUserProfileHandler)
		auth.POST("/my/cart", handlers.SaveCartHandler)
		auth.GET("/my/cart", handlers.GetSavedCartHandler)
		auth.GET("/my/orders", handlers.GetOrderHistoryHandler)
		auth.GET("my/orders/:id", handlers.GetOrderDetailHandler)
		auth.POST("/checkout", handlers.CheckoutHandler)

		// Admin routes
		admin := auth.Group("/admin")
		admin.Use(middlewares.AdminMiddleware())
		{
			admin.GET("/users", handlers.AdminGetUsersHandler)
			admin.GET("/users/:id", handlers.AdminGetUsersHandler)
			admin.DELETE("/users/:id", handlers.AdminGetUsersHandler)

			admin.GET("/products", handlers.AdminListProductHandler)
			admin.POST("/products", handlers.AdminCreateProductHandler)
			admin.PATCH("/products/:id", handlers.AdminUpdateProductHandler)
			admin.DELETE("/products/:id", handlers.AdminDeleteProductHandler)
			
			admin.GET("/orders", handlers.AdminListOrdersHandler)
			admin.GET("/orders/:id", handlers.AdminGetOrder)
			admin.PUT("/orders/:id/status", handlers.AdminUpdateOrderStatus)
		}
	}

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}