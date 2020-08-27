package main

import (
	"context"
	"os"
	"time"

	"github.com/kashifmin/gogoose"
	"github.com/kashifmin/gogoose/gen"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoClient creates a new mongodb client
func NewMongoClient() *mongo.Database {
	uri := os.Getenv("MONGO_URI")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.
		Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client.Database("test")
}

func main() {
	db := NewMongoClient()
	userModel := gen.NewUserModel(db.Collection("kuser"))
	oid := primitive.NewObjectID()
	doc := userModel.New(&gogoose.User{Name: "Kashif", Age: 23, ID: &oid})
	err := doc.Save(context.Background())
	if err != nil {
		panic(err)
	}
}
