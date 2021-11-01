package mongodb

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MONGO_URI      string = "mongodb://localhost:27017"
	MONGO_USERNAME string = "occitan"
	MONGO_PASSWORD string = "naticco"
	MONGO_DB_NAME  string = "oc"
)

func getMongoEnv() {
	mongoUri := os.Getenv("MONGO_URI")
	if mongoUri != "" {
		MONGO_URI = mongoUri
	}

	mongoUsername := os.Getenv("MONGO_USERNAME")
	if mongoUsername != "" {
		MONGO_USERNAME = mongoUsername
	}

	mongoPassword := os.Getenv("MONGO_PASSWORD")
	if mongoPassword != "" {
		MONGO_PASSWORD = mongoPassword
	}

	mongoDBname := os.Getenv("MONGO_DB_NAME")
	if mongoDBname != "" {
		MONGO_DB_NAME = mongoDBname
	}
}

func NewMongoClient() *mongo.Database {
	getMongoEnv()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(MONGO_URI).SetDirect(true).SetAuth(options.Credential{
		Username: MONGO_USERNAME,
		Password: MONGO_PASSWORD,
	}))
	if err != nil {
		log.Fatalln(err)
	}

	db := client.Database(MONGO_DB_NAME)
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
