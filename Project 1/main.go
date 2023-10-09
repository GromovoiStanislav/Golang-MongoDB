package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Person struct {
	Name  string
	Age   int
	Email string
}

func main() {

	// Устанавливаем параметры подключения к MongoDB.
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Проверяем соединение с сервером MongoDB.
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// Получаем коллекцию в базе данных "mydb".
	collection := client.Database("mydb").Collection("people")

	// Создание записи (Create)
	person := Person{Name: "John Doe", Age: 30, Email: "john.doe@example.com"}
	insertResult, err := collection.InsertOne(context.Background(), person)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)

	id := insertResult.InsertedID

	// Чтение записи (Read)
	var result Person
	filter := bson.M{"name": "John Doe"}
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found document: %+v\n", result)

	// Чтение записи (Read)
	filter = bson.M{"_id": id}
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found document: %+v\n", result)

	// Создание объекта ObjectID по строковому значению
	objectID, err := primitive.ObjectIDFromHex("6523e18868778894b7187f70")
	if err != nil {
		log.Fatal(err)
	}
	// Чтение записи (Read)
	filter = bson.M{"_id": objectID}
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found document by ID: %+v\n", result)

	// Обновление записи (Update)
	update := bson.M{"$set": bson.M{"age": 31}}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and modified %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// Удаление записи (Delete)
	//deleteResult, err := collection.DeleteOne(context.Background(), filter)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Deleted %v documents.\n", deleteResult.DeletedCount)

	/// поиск ВСЕХ записей
	filter = bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем срез для хранения результатов.
	var results []Person

	// Используем cursor.All, чтобы извлечь все результаты в срез.
	if err := cursor.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}

	// Выводим все результаты.
	for _, result := range results {
		fmt.Printf("Found document: %+v\n", result)
	}

	//// Получаем НЕСКОЛЬКО записей
	// Создаем фильтр, который выберет все документы, соответствующие вашему критерию.
	filter = bson.M{"age": bson.M{"$gt": 30}} // Пример: выбор всех людей старше 30 лет.

	// Выполняем запрос и получаем курсор на результаты.
	cursor, err = collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем срез для хранения результатов.
	//var results []Person

	// Используем cursor.All, чтобы извлечь все результаты в срез.
	if err := cursor.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}

	// Выводим все найденные документы.
	for _, result := range results {
		fmt.Printf("Found document 2: %+v\n", result)
	}

	// Закрываем соединение с MongoDB.
	client.Disconnect(context.Background())
}
