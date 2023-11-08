package controllers

import (
	"github.com/Micah-Shallom/modules/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func GetMenu() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func CreateMenu() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func UpdateMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		
	}
}