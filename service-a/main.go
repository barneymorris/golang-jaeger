package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	connectMongo()
	setupWebServer()
}

func connectMongo() {
	opts := options.Client()
	opts.ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(context.Background(), opts)
	//Seed the database with todo's
	docs := []interface{}{
		bson.D{{"id", "1"}, {"title", "Buy groceries"}},
		bson.D{{"id", "2"}, {"title", "install Aspecto.io"}},
		bson.D{{"id", "3"}, {"title", "Buy dogz.io domain"}},
	}
	client.Database("todo").Collection("todos").InsertMany(context.Background(), docs)
}

func setupWebServer() {
	r := gin.Default()
	r.GET("/todo", func(c *gin.Context) {
		collection := client.Database("todo").Collection("todos")
		//Important: Make sure to pass c.Request.Context() as the context and not c itself - TBD
		cur, findErr := collection.Find(c.Request.Context(), bson.D{})
		if findErr != nil {
			c.AbortWithError(500, findErr)
			return
		}
		results := make([]interface{}, 0)
		curErr := cur.All(c, &results)
		if curErr != nil {
			c.AbortWithError(500, curErr)
			return
		}
		c.JSON(http.StatusOK, results)
	})
	_ = r.Run(":8080")
}
