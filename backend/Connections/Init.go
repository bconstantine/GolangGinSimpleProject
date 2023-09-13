package routes

import (
	"context"
	"fmt" // for printing to console
	"log" // for logging errors
	"os" // for getting environment variables
	"time" // for setting timeouts

	"github.com/joho/godotenv" // for loading environment variables from .env file
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DBinstance func
//Returning a mongo.Client instance
func DBinstance() *mongo.Client {
	err := godotenv.Load(".env") // load .env file as early as possible

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MongoDb := os.Getenv("MONGODB_URL") // get the mongodb url from the .env file

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		//Immediately terminate the function
		log.Fatal(err)
	}

	//Context.Background() returns a non-nil, empty Context. It is never canceled, has no values, and has no deadline.
	//Context.WithTimeout returns a copy of parent with the timeout adjusted to be no longer than timeout.
	//cancel terminates the context, releasing any resources associated with it.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	//defer is used to ensure that a function call is performed later in a program’s execution, usually for purposes of cleanup.
	defer cancel()

	//Connect to MongoDB
	//placing context for Mongodb Connect, auto cancels when timeout is reached
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	return client
}

//Client Database instance
var Client *mongo.Client = DBinstance()

//OpenCollection is a  function makes a connection with a collection in the database
func OpenCollection(databaseName string, client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database(databaseName).Collection(collectionName)
	return collection
}