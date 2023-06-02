package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/spanner"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/iterator"
)

var dbName = "company"
var colName = "employees"

var collection *mongo.Collection

type insertemp struct {
	EmployeesName string `bson:"employees_name"`
	Age           int64  `bson:"age"`
}

type SpannerData struct {
	EmployeesName string `spanner:"employees_name"`
	Age           int64  `spanner:"age"`
}

func main() {
	connectMongoDB()
}

func connectMongoDB() {
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb+srv://nencybatada:dhvani@cluster0.ihgut8k.mongodb.net/?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.Background())

	// Ping MongoDB to verify the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")

	// Fetch data from MongoDB
	collection = client.Database(dbName).Collection(colName)
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())
	fmt.Println("Data fetched")
	fmt.Println("Database name: ", dbName, "Column Name : ", colName)

	if err := os.Setenv("SPANNER_EMULATOR_HOST", "localhost:9010"); err != nil {
		fmt.Println("Error while setting environment variable: ")
	}

	// Connect to Cloud Spanner
	spannerClient, err := spanner.NewClient(context.Background(), "projects/spanner-project/instances/spanner-instance/databases/spanner-database")
	if err != nil {
		panic(err)
	}
	defer spannerClient.Close()
	fmt.Println("Connected to spanner")

	// Insert MongoDB data into Cloud Spanner
	mutations := []*spanner.Mutation{}
	for cursor.Next(context.Background()) {
		var result insertemp
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}

		spannerData := SpannerData{
			EmployeesName: result.EmployeesName,
			Age:           result.Age,
		}

		mutations = append(mutations, spanner.InsertOrUpdate("employees", []string{"employees_name", "age"}, []interface{}{spannerData.EmployeesName, spannerData.Age}))
	}

	_, err = spannerClient.Apply(context.Background(), mutations)
	if err != nil {
		panic(err)
	}
	fmt.Println("Data inserted into Cloud Spanner!")

	// // Read data from Cloud Spanner
	stmt := spanner.Statement{SQL: "SELECT employees_name, age FROM employees"}
	iter := spannerClient.Single().Query(context.Background(), stmt)
	defer iter.Stop()

	// Process Cloud Spanner query results
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		var spannerData SpannerData
		if err := row.ToStruct(&spannerData); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Name: %s, Age: %d\n", spannerData.EmployeesName, spannerData.Age)
	}
}
