package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Micah-Shallom/modules/database"
	"github.com/Micah-Shallom/modules/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := orderCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing tables"})
			return
		}

		var allTables []bson.M
		if err := result.All(ctx, &allTables); err != nil {
			log.Fatal(err)
			return
		}
		defer cancel()
		
	}
}

func GetTable() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table
		tableID := c.Param("table_id")

		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableID}).Decode(&table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while fetching table"})
			return
		}
		c.JSON(http.StatusOK, table)
		defer cancel()
	}
}

func CreateTable() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table

		if err := c.ShouldBindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(table)
		if validationErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
			return
		}

		table.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.ID = primitive.NewObjectID() 
		table.TableID = table.ID.Hex()

		result, err := tableCollection.InsertOne(ctx, table)
		defer cancel()
		if err != nil {
			msg := "Table item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func UpdateTable() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table
		tableID := c.Param("table_id")

		if err := c.ShouldBindJSON(&table); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if table.NumberOfGuests != nil {
			updateObj = append(updateObj, bson.E{"number_of_guests", table.NumberOfGuests})
		}
		if table.TableNumber != nil {
			updateObj = append(updateObj, bson.E{"table_number",table.TableNumber})
		}
	
		table.UpdatedAt, _ = time.Parse(time.RFC3339, string(time.Now().Format(time.RFC3339)))

		upsert := true
		filter := bson.M{"table_id": tableID}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := tableCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"&set", updateObj},
			},
			&opt,
		)
		defer cancel()
		if err != nil {
			msg := fmt.Sprintf("Table Item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}