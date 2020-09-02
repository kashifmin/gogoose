package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kashifmin/gogoose/examples/types"
	"github.com/kashifmin/gogoose/gen"
	"go.mongodb.org/mongo-driver/bson"
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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	db := NewMongoClient()
	userModel := gen.NewUserModel(db.Collection("kuser"))
	users, err := userModel.Find(context.Background(), bson.M{})
	check(err)
	for _, usr := range users {
		fmt.Println(usr.GetRaw())
		usr.GetRaw().Age = 27
		err = usr.Save(context.Background())
		check(err)
	}
	oid := primitive.NewObjectID()
	doc := userModel.New(&types.User{Name: "Kashif", Age: 28, ID: &oid})
	// save a document
	err = doc.Save(context.Background())
	if err != nil {
		panic(err)
	}
}
