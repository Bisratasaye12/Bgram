package config

import (
	models "BChat/Domain/Models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB(env *models.Env) *mongo.Database {

	clientOptions := options.Client().ApplyURI(env.DBURL)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(env.DBNAME)
	return db
}
