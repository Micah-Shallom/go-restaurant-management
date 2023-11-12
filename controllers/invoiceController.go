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

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func GetInvoices() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := invoiceCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing invoice items"})
			return
		}
		var allInvoices []bson.M
		if err = result.All(ctx, &allInvoices); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allInvoices)
	}	
}

func GetInvoice() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var invoice models.Invoice
		invoiceID := c.Param("invoice_id")
		filter := bson.M{"invoice_id":invoiceID}
		err := invoiceCollection.FindOne(ctx,filter).Decode(&invoice)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while getting invoice item"})
			return
		}

		var invoiceView models.InvoiceViewFormat
		allOrderItems, err := ItemsByOrder(invoice.OrderID)
		invoiceView.OrderID = invoice.OrderID
		invoiceView.PaymentDueDate = invoice.PaymentDueDate

		invoiceView.PaymentMethod = "null"
		if invoice.PaymentMethod != nil {
			invoiceView.PaymentMethod = *invoice.PaymentMethod
		}

		invoiceView.InvoiceID = invoice.InvoiceID
		invoiceView.PaymentStatus = invoice.PaymentStatus
		invoiceView.PaymentDue = allOrderItems[0]["payment_due"]
		invoiceView.TableNumber = allOrderItems[0]["table_number"]
		invoiceView.OrderDetails = allOrderItems[0]["order_items"]

		c.JSON(http.StatusOK, invoiceView)
	}
}

func CreateInvoice() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var invoice models.Invoice
		var order models.Order

		if err := c.ShouldBindJSON(&invoice); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.OrderID}).Decode(&order)
		if err != nil {
			msg := fmt.Sprintf("Message: Order was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}
		status := "PENDING"
		if invoice.PaymentStatus == nil {
			invoice.PaymentStatus = &status
		}

		invoice.PaymentDueDate, _ = time.Parse(time.RFC3339, time.Now().AddDate(0,0,1).Format(time.RFC3339))
		invoice.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.ID = primitive.NewObjectID()
		invoice.InvoiceID = invoice.ID.Hex()

		validationErr := validate.Struct(invoice)
		if validationErr != nil {
			c.JSON(http.StatusInternalServerError, validationErr.Error())
			return
		}

		result, insertErr := invoiceCollection.InsertOne(ctx, invoice)
		if insertErr != nil {
			msg := fmt.Sprintf("Invoice item was not created")
			c.JSON(http.StatusInternalServerError, msg)
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, result)
	}
}

func UpdateInvoice() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		var invoice models.Invoice
		invoiceID := c.Param("invoice_id")

		if err := c.ShouldBindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"invoice_id":invoiceID}

		var updateObj primitive.D

		if invoice.PaymentMethod != nil {
			updateObj = append(updateObj, bson.E{"payment_method", invoice.PaymentMethod})
		}
		if invoice.PaymentStatus != nil {
			updateObj = append(updateObj, bson.E{"payment_status", invoice.PaymentStatus})
		}
		invoice.UpdatedAt, _ =  time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", invoice.UpdatedAt})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		status := "PENDING"
		if invoice.PaymentStatus == nil {
			invoice.PaymentStatus = &status
		}

		result, err := invoiceCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)
		defer cancel()
		if err != nil {
			msg := fmt.Sprintf("invoice item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}