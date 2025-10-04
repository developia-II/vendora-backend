package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(router *gin.Engine, db *mongo.Database) {
	logrus.Info("Setting up routes...")

	router.GET("/", func(c *gin.Context) {
		logrus.Info("Handling request to /")
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is running!",
			"status":  "ok",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		logrus.Info("Handling request to /health")
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "vendora-backend",
		})
	})

	if db != nil {
		logrus.Info("Database connected - setting up database routes")
		router.POST("/api/v1/auth/register", NewAuthHandler(db).CreateUser)
		router.POST("/api/v1/auth/verify", NewAuthHandler(db).VerifyEmail)
		router.POST("/api/v1/auth/login", NewAuthHandler(db).LoginUser)
		router.POST("/api/v1/auth/forgot-password", NewAuthHandler(db).ForgotPassword)
		router.POST("/api/v1/auth/reset-password", NewAuthHandler(db).ResetPassword)

	} else {
		logrus.Warn("Database not connected - running with limited functionality")
	}
}
