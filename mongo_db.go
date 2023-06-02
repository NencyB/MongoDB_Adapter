package main

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// var dbName = "company"
// var colName = "employees"

// var collection *mongo.Collection

// type insertemp struct {
// 	EmployeesName string `bson:"employees_name"`
// 	Age           int    `bson:"age"`
// }

// func main() {
// 	connect()
// }

// func connect() {
// 	// Use the SetServerAPIOptions() method to set the Stable API version to 1
// 	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
// 	opts := options.Client().ApplyURI("mongodb+srv://nencybatada:dhvani@cluster0.ihgut8k.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)

// 	// Create a new client and connect to the server
// 	client, err := mongo.Connect(context.TODO(), opts)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer func() {
// 		if err = client.Disconnect(context.TODO()); err != nil {
// 			panic(err)
// 		}
// 	}()

// 	// Send a ping to confirm a successful connection
// 	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

// 	collection = client.Database(dbName).Collection(colName)

// 	inserted, err := collection.Find(context.Background(), bson.D{{}})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var employee []insertemp

// 	if err = inserted.All(context.TODO(), &employee); err != nil {
// 		fmt.Println(err)
// 	}
// 	defer inserted.Close(context.Background())
// 	for _, i := range employee {
// 		fmt.Println(i)
// 	}
// }
