package controllers

import (
	"github.com/Micah-Shallom/modules/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func GetTable() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func CreateTable() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func UpdateTable() gin.HandlerFunc{
	return func(c *gin.Context){
		
	}
}