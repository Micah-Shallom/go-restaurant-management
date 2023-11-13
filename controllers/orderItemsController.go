package controllers

import (
	"context"
	"log"
	"net/http"
	"options"
	"time"

	"github.com/Micah-Shallom/modules/database"
	"github.com/Micah-Shallom/modules/helpers"
	"github.com/Micah-Shallom/modules/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func CreateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItemPack models.OrderItemPack
		var order models.Order

		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order.OrderDate, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItemsToBeInserted := []interface{}{}
		order.TableID = orderItemPack.TableID
		orderID := OrderItemOrderCreator(order)

		for _, orderItem := range orderItemPack.OrderItems {
			orderItem.OrderID = orderID

			validationErr := validate.Struct(orderItem)
			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}

			orderItem.ID = primitive.NewObjectID()
			orderItem.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.OrderItemID = orderItem.ID.Hex()
			var num = helpers.ToFixed(*orderItem.UnitPrice, 2)
			orderItem.UnitPrice = &num
			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}

		insertedOrderItems, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)

		if err != nil {
			log.Fatal(err)
		}
		defer cancel()
		c.JSON(http.StatusOK, insertedOrderItems)
	}
}

func GetOrderItems() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel =  context.WithTimeout(context.Background(), 100*time.Second)
		result, err := orderItemCollection.Find(context.TODO(), bson.M{})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occurred while listing ordered items"})
			return
		}

		var allOrderedItems []bson.M
		if err = result.All(ctx, &allOrderedItems); err != nil {{
			log.Fatal(err)
			return
		}}
		c.JSON(http.StatusOK, result)
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		orderID := c.Param("order_id")
		allOrderedItems, err := ItemsByOrder(orderID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing order items by order ID"})
			return
		}
		c.JSON(http.StatusOK, allOrderedItems)
	}
}

func ItemsByOrder(id string)(OrderItems []primitive.M, err error){
	
}

func GetOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		orderItemID := c.Param("order_item_id")

		var orderItem models.OrderItem
		
		err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemID}).Decode(&orderItem)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while getting order item"})
			return
		}
		c.JSON(http.StatusOK, orderItem)
	}
}


func UpdateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItem models.OrderItem

		if err := c.ShouldBindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}

		orderItemID := c.Param("order_item_id")

		filter := bson.M{"order_item_id": orderItemID}

		var updateObj primitive.D

		if orderItem.UnitPrice != nil {
			updateObj = append(updateObj, bson.E{"unit_price", orderItem.UnitPrice})
		}
		if orderItem.Quantity != nil {
			updateObj =  append(updateObj, bson.E{"quantity", *orderItem.Quantity})
		}
		if orderItem.FoodID != nil {
			updateObj = append(updateObj, bson.E{"food_id", *orderItem.FoodID})
		}

		orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", orderItem.UpdatedAt})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)
		defer cancel()
		if err != nil {
			msg := "Order Item Update Failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}