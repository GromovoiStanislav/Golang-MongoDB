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

func readDocumentByID(collection *mongo.Collection, id interface{}) Person {
	// Преобразуем id в строку
	idStr, ok := id.(primitive.ObjectID)
	if !ok {
		log.Fatal("Failed to convert ID to string")
	}

	var result Person
	filter := bson.M{"_id": idStr}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func getAllDocuments(collection *mongo.Collection) []Person {
	/// поиск ВСЕХ записей
	filter := bson.M{}
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

	return results
}

func getManyDocuments(collection *mongo.Collection, filter bson.M) []Person {
	/// Получаем НЕСКОЛЬКО записей
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

	return results
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

	// Создание записи асинхронно (Create)
	insertChannel := make(chan *mongo.InsertOneResult)
	go func() {
		person := Person{Name: "Tom Sawyer", Age: 30, Email: "tom.sawyer@example.com"}
		insertResult, err := collection.InsertOne(context.Background(), person)
		if err != nil {
			log.Fatal(err)
		}
		insertChannel <- insertResult
	}()

	// Чтение записи асинхронно (Read)
	readChannel := make(chan Person)
	go func() {
		var result Person
		filter := bson.M{"name": "Tom Sawyer"}
		err := collection.FindOne(context.Background(), filter).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		readChannel <- result
	}()

	// Обновление записи асинхронно (Update)
	updateChannel := make(chan *mongo.UpdateResult)
	go func() {
		filter := bson.M{"name": "Tom Sawyer"}
		update := bson.M{"$set": bson.M{"age": 31}}
		updateResult, err := collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			log.Fatal(err)
		}
		updateChannel <- updateResult
	}()

	// Удаление записи асинхронно (Delete)
	deleteChannel := make(chan *mongo.DeleteResult)
	go func() {
		filter := bson.M{"name": "Tom Sawyer"}
		deleteResult, err := collection.DeleteOne(context.Background(), filter)
		if err != nil {
			log.Fatal(err)
		}
		deleteChannel <- deleteResult
	}()

	// Ожидаем завершения всех асинхронных операций.
	insertResult := <-insertChannel
	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)

	// Чтение записи асинхронно (Read)
	readChannel2 := make(chan Person)
	go func(id interface{}) {
		var result Person
		filter := bson.M{"_id": id}
		err := collection.FindOne(context.Background(), filter).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		readChannel2 <- result
	}(insertResult.InsertedID)
	// После этого можно извлечь результат из канала readChannel2
	result := <-readChannel2
	fmt.Printf("Found document by ID: %+v\n", result)

	// Чтение записи асинхронно (Read)
	readChannel3 := make(chan Person)
	go func(id interface{}) {
		result := readDocumentByID(collection, id)
		readChannel3 <- result
	}(insertResult.InsertedID)
	// После этого можно извлечь результат из канала readChannel3
	result = <-readChannel3
	fmt.Printf("Found document by ID: %+v\n", result)

	readResult := <-readChannel
	fmt.Printf("Found document: %+v\n", readResult)

	updateResult := <-updateChannel
	fmt.Printf("Matched %v documents and modified %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// Чтение записи асинхронно (Read)
	readChannel4 := make(chan []Person)
	go func() {
		results := getAllDocuments(collection)
		readChannel4 <- results
	}()
	// Выводим все результаты.
	results := <-readChannel4
	close(readChannel4)
	fmt.Printf("Found ALL documents: %+v\n", results)
	for _, result := range results {
		fmt.Printf("One of ALL documents: %+v\n", result)
	}

	// Чтение записи асинхронно (Read)
	readChannel5 := make(chan []Person)
	go func() {
		// Создаем фильтр, который выберет все документы, соответствующие вашему критерию.
		filter := bson.M{"age": bson.M{"$gt": 30}} // Пример: выбор всех людей старше 30 лет.
		results := getManyDocuments(collection, filter)
		readChannel5 <- results
	}()
	// Выводим все результаты.
	results = <-readChannel5
	close(readChannel5)
	fmt.Printf("Found Many documents: %+v\n", results)
	for _, result := range results {
		fmt.Printf("One of Many documents: %+v\n", result)
	}

	deleteResult := <-deleteChannel
	fmt.Printf("Deleted %v documents.\n", deleteResult.DeletedCount)

	// Закрываем соединение с MongoDB.
	client.Disconnect(context.Background())
}
