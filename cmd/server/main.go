package main

import (
	"os"

	"github.com/developia-II/ecommerce-backend/internal/database"
	"github.com/developia-II/ecommerce-backend/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Starting server...")
	gin.SetMode(gin.DebugMode)
	logrus.SetLevel(logrus.InfoLevel)

	logrus.Info("Attempting to connect to database...")
	db, err := database.ConnectToDB()
	if err != nil {
		logrus.WithError(err).Warn("Failed to connect to DB - running without database")
		db = nil
	} else {
		logrus.Info("Successfully connected to DB")
	}

	logrus.Info("Setting up Gin router...")
	router := gin.Default()
	logrus.Info("Calling handlers.SetupRoutes...")
	handlers.SetupRoutes(router, db)

	logrus.Info("Loading environment variables...")
	if err := godotenv.Load(); err != nil {
		logrus.WithError(err).Error("failed to load env")
	}
	PORT := os.Getenv("PORT")
	logrus.Info("Starting server on port: " + PORT)
	router.Run(PORT)
}
