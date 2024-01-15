package db

import (
	"context"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientInstance *mongo.Client

var mongoOnce sync.Once

var clientInstanceError error

func GetMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {

		if err := godotenv.Load(); err != nil {
			panic("Ошибка загрузки файла .env: " + err.Error())
		}
		mongoURL := os.Getenv("MONGODB_URI")

		if mongoURL == "" {
			mongoURL = "mongodb://localhost:27017"
		}


		clientOptions := options.Client().ApplyURI(mongoURL)

		client, err := mongo.Connect(context.Background(), clientOptions)

		clientInstance = client

		clientInstanceError = err
	})

	return clientInstance, clientInstanceError
}