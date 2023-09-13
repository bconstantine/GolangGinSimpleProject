package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


// Order struct
//In format of:
// Field Name, Field Type, Field Description/Tag,
	//bson tag means it is saved in the database as that name
	//json tag means it is sent to the frontend as that name
type Order struct {
	ID         	primitive.ObjectID 	`bson:"_id"`
	Dish       	*string            	`json:"dish" validate:"required,min=2,max=100"`
	Price   	*float64           	`json:"price" validate:"required,min=0,max=100"`
	Server		*string			`json:"server" validate:"min=2,max=100"`
	Table		*string			`json:"table" validate:"max=2"`
}