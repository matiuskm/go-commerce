package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/matiuskm/go-commerce/config"
	"github.com/matiuskm/go-commerce/db"

	"github.com/matiuskm/go-commerce/handlers"
	"github.com/matiuskm/go-commerce/middlewares"
)

func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; object-src 'none';")
		// Consider adding other headers like:
		// c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		// c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Next()
	}
}

func main() {
	config.LoadEnv()

	db.Init()

	r := gin.New()

	// Logging middleware
	r.Use(gin.Logger())

	// Recovery middleware
	r.Use(gin.Recovery())

	// Custom security headers middleware
	r.Use(SecurityHeadersMiddleware())

	// set cors
	originEnv := os.Getenv("CORS_ORIGINS")
	log.Println("Loaded CORS_ORIGINS:", originEnv)

	if originEnv == "" {
		log.Fatal("CORS_ORIGINS env var is required")
	}

	corsOrigins := strings.Split(originEnv, ",")
	r.Use(cors.New(cors.Config{
		AllowOrigins:  corsOrigins,
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("[%s] %s - %d - %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start))
	})

	r.Static("/uploads", "./uploads")
	r.StaticFile("/robots.txt", "./robots.txt")

	// Public routes
	r.POST("/auth/register", handlers.RegisterHandler)
	r.POST("/auth/login", handlers.LoginHandler)

	r.GET("/products", handlers.GetAllProductsHandler)
	r.GET("/products/:id", handlers.GetProductByIDHandler)

	r.POST("/webhooks/xendit", handlers.XenditWebhookHandler)

	// Protected routes
	auth := r.Group("/")
	auth.Use(middlewares.AuthMiddleware())
	{
		// User routes
		auth.GET("/my/profile", handlers.GetUserProfileHandler)
		auth.PATCH("/my/profile", handlers.UpdateUserProfileHandler)
		auth.POST("/my/cart", handlers.SaveCartHandler)
		auth.GET("/my/cart", handlers.GetSavedCartHandler)
		auth.DELETE("/my/cart/:productID", handlers.RemoveCartItemHandler)
		auth.GET("/my/orders", handlers.GetOrderHistoryHandler)
		auth.GET("my/orders/:id", handlers.GetOrderDetailHandler)
		auth.GET("/my/addresses", handlers.ListAddresses)
		auth.POST("/my/addresses", handlers.CreateAddress)
		auth.PUT("/my/addresses/:id", handlers.UpdateAddress)
		auth.DELETE("/my/addresses/:id", handlers.DeleteAddress)
		auth.POST("/checkout", handlers.CheckoutHandler)

		// Admin routes
		admin := auth.Group("/admin")
		admin.Use(middlewares.AdminMiddleware())
		{
			admin.GET("/users", handlers.AdminListUsersHandler)
			admin.GET("/users/:id", handlers.AdminGetUserHandler)
			admin.PUT("/users/:id", handlers.AdminUpdateUserHandler)
			admin.PATCH("/users/:id/role", handlers.AdminUpdateUserRoleHandler)
			admin.DELETE("/users/:id", handlers.AdminDeleteUserHandler)

			admin.GET("/products", handlers.AdminListProductHandler)
			admin.GET("/products/:id", handlers.AdminGetProductHandler)
			admin.POST("/products", handlers.AdminCreateProductHandler)
			admin.PATCH("/products/:id", handlers.AdminUpdateProductHandler)
			admin.DELETE("/products/:id", handlers.AdminDeleteProductHandler)
			admin.POST("/products/:id/image", handlers.AdminUploadImageHandler)


			admin.GET("/orders", handlers.AdminListOrdersHandler)
			admin.GET("/orders/:id", handlers.AdminGetOrder)
			admin.PUT("/orders/:id/status", handlers.AdminUpdateOrderStatus)
		}
	}
 
	log.Println("Server started on :8080")
	port := os.Getenv("PORT")
	// http.ListenAndServe(":"+port, r)
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	log.Println("Server started on :" + port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
