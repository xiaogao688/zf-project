package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"go-project/learn/gin/gin-jwt/gin-jwt/controllers"
	"go-project/learn/gin/gin-jwt/gin-jwt/middlewares"
	"go-project/learn/gin/gin-jwt/gin-jwt/models"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file. %v\n", err)
	}
}

func main() {
	models.ConnectDatabase()
	r := gin.Default()
	public := r.Group("/api")
	{
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)
	}

	protected := r.Group("/api/admin")
	{
		protected.Use(middlewares.JwtAuthMiddleware()) // 在路由组中使用中间件
		protected.GET("/user", controllers.CurrentUser)
	}

	r.Run("0.0.0.0:8000")
}
