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

var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc{
	return func(c *gin.Context) {
		// var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := orderCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing order items"})
			return
		}

		var allOrders []bson.M

		if err := result.All(ctx, &allOrders); err != nil {
			log.Fatal(err)
		}
		defer cancel()
		c.JSON(http.StatusOK, allOrders)
	}
}

func GetOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var order models.Order
		orderID := c.Param("order_id")

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderID}).Decode(&order)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while fetching the food item"})
			return
		}
		c.JSON(http.StatusOK, order)
	}
}

func CreateOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var order models.Order
		var table models.Table

		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return
		}

		validationErr := validate.Struct(&order)
		if validationErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
			return
		}

		if order.TableID != nil {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.TableID}).Decode(&table)
			if err != nil {
				msg := fmt.Sprintf("message: Table was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
		}

		order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.ID = primitive.NewObjectID()
		order.OrderID = order.ID.Hex()

		result, insertErr := orderCollection.InsertOne(ctx, order)
		defer cancel()

		if insertErr != nil {
			msg := fmt.Sprintf("Order Item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		var table models.Table
		var order models.Order

		var updateObj primitive.D

		orderID := c.Param("order_id")
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if order.TableID != nil {
			err := menuCollection.FindOne(ctx, bson.M{"table_id": order.TableID}).Decode(&table)
			defer cancel()
			if err != nil {
				msg := fmt.Sprintf("message: menu was found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}

			updateObj = append(updateObj, bson.E{"menu", order.TableID})
		}

		order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at",  order.UpdatedAt})

		upsert := true

		filter := bson.M{"order_id": orderID}

		opts := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opts,
		)	

		if err != nil {
			msg := fmt.Sprintf("order item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}

}

func OrderItemOrderCreator(order models.Order) string{
	order.CreatedAt, _ =  time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ =  time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.OrderID = order.ID.Hex()

	orderCollection.InsertOne(ctx, order)
	defer cancel()

	return order.OrderID
}