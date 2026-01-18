package main

import (
	"log"
	"net/http"
	"os"

	"slink-backend/internal/api"
	"slink-backend/internal/config"
	"slink-backend/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	gin.SetMode(gin.DebugMode)

	cfg := config.Load()

	supabaseClient, err := database.ConnectSupabase(cfg.SupabaseURL, cfg.SupabaseKey, cfg.SupabaseProjectRef)
	if err != nil {
		log.Fatal("Failed to connect to Supabase:", err)
	}

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	apiHandler := api.NewHandler(supabaseClient, cfg)

	apiGroup := router.Group("/api")
	{
		apiGroup.POST("/register", apiHandler.Register)
		apiGroup.POST("/login", apiHandler.Login)

		protected := apiGroup.Group("/")
		protected.Use(apiHandler.AuthMiddleware())
		{
			protected.GET("/profile", apiHandler.GetProfile)
			protected.POST("/shorten", apiHandler.ShortenURL)
			protected.GET("/links", apiHandler.GetLinksByUser)
		}

		apiGroup.GET("/qr/:shortCode", apiHandler.GenerateQR)
	}

	router.GET("/:shortCode", apiHandler.RedirectURL)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "slink-backend",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
