package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"


	"mongo-example/internal/db"
)

const  (
	Database = "products-api"
	ProductsCollection = "products"
)

type Product struct {
	ID        primitive.ObjectID
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
}

func main() {
	testConnection()
	createProduct()
	getAllProducts()
}

func testConnection() {
	
	client, err := db.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

}

func createProduct() {

	client, err := db.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database(Database).Collection(ProductsCollection)

	product := Product{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Title: "Product #1",
	}
	
	_, err = collection.InsertOne(context.Background(), product)
	if err != nil {
		log.Println(err) 
	}

}

func getAllProducts() {

	client, err := db.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database(Database).Collection(ProductsCollection)

	var products []Product

	cur, err := collection.Find(context.Background(), bson.D{
		primitive.E{},
	})
	if err != nil {
		log.Println(err) 
	}

	for cur.Next(context.Background()) {
		var p Product

		err := cur.Decode(&p)
		if err != nil {
			log.Println(err) 
		}

		products = append(products, p)

		//fmt.Println(p)
	}

	fmt.Println(products)

	// for _, product := range products {
	// 	fmt.Println(product)
	// }
	
}