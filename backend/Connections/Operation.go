package routes

import (
	"context"
	"fmt"
	"net/http" // for sending http responses
	"time" // for setting timeouts

	"server/Models" 
	"github.com/go-playground/validator/v10" // for validating struct fields
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive" // for generating a new ObjectID
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gin-gonic/gin"
)

var validate = validator.New()

var orderCollection *mongo.Collection = OpenCollection("GolangGinSimpleProject", Client, "orders")

//Create an order
func CreateOrder(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var order models.Order

	//BindJSON: Parse JSON from the request body.
	if err := c.BindJSON(&order); err != nil {
		//If the JSON is invalid, return an "Bad Request" error.
		//gin.H is a shortcut for map[string]interface{}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	//Validate the order
	validationErr := validate.Struct(order)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}
	order.ID = primitive.NewObjectID()

	//Insert the order into the database
	result, insertErr := orderCollection.InsertOne(ctx, order)
	if insertErr != nil {
		msg := fmt.Sprintf("order item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

//get all orders
func GetOrders(c *gin.Context){
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var orders []bson.M //array of bson.M, which is a map[string]interface{}

	cursor, err := orderCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	
	//Iterate the cursor and decode each item into a bson.M (map[string]interface{})
	if err = cursor.All(ctx, &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	fmt.Println(orders)

	c.JSON(http.StatusOK, orders)
}

//get all orders by the server's name
func GetOrdersByServer(c *gin.Context){

	server := c.Params.ByName("server")
	
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var orders []bson.M

	cursor, err := orderCollection.Find(ctx, bson.M{"server": server})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	if err = cursor.All(ctx, &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	fmt.Println(orders)

	c.JSON(http.StatusOK, orders)
}

//get an order by its id
func GetOrderById(c *gin.Context){

	orderID := c.Params.ByName("id")
	//convert the string id to a primitive.ObjectID type
	//the second return value is an error, which we ignore using _
	docID, err := primitive.ObjectIDFromHex(orderID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var order bson.M

	if err := orderCollection.FindOne(ctx, bson.M{"_id": docID}).Decode(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	fmt.Println(order)

	c.JSON(http.StatusOK, order)
}

//update a server's name for an order
func UpdateServer(c *gin.Context){

	orderID := c.Params.ByName("id")
	docID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	type Server struct {
		Server		*string				`json:"server"`
	}

	var server Server

	if err := c.BindJSON(&server); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	//bson.D is an ordered representation of a BSON document which is a Map
	result, err := orderCollection.UpdateOne(ctx, bson.M{"_id": docID}, 
		bson.D{
			{"$set", bson.D{{"server", server.Server}}},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, result.ModifiedCount)

}

//update the order
func UpdateOrder(c *gin.Context){

	orderID := c.Params.ByName("id")
	docID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var order models.Order

	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	validationErr := validate.Struct(order)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}

	result, err := orderCollection.ReplaceOne(
		ctx,
		bson.M{"_id": docID},
		bson.M{
			"dish":  order.Dish,
			"price": order.Price,
			"server": order.Server,
			"table": order.Table,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, result.ModifiedCount)
}

//delete an order given the id
func DeleteOrder(c * gin.Context){
	orderID := c.Params.ByName("id")
	docID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := orderCollection.DeleteOne(ctx, bson.M{"_id": docID})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, result.DeletedCount)
}