package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToDB() (*mongo.Database, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	var MONGO_URI = os.Getenv("MONGO_URI")
	fmt.Println(MONGO_URI)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := mongo.Connect(ctx, options.Client().
		SetMaxPoolSize(10).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(30*time.Second).
		SetServerSelectionTimeout(30*time.Second).
		SetConnectTimeout(30*time.Second).
		ApplyURI(MONGO_URI))

	if err != nil {
		return nil, err
	}

	if err := res.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return res.Database("vendora"), nil
}
