package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kwamekyeimonies/restaurant_management_system_backend/database"
	"github.com/kwamekyeimonies/restaurant_management_system_backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

type OrderItemPack struct {
	Table_id    *string
	Order_items []models.OrderItem
}

func GetOrderItems() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := orderItemCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error Occured Listing OrderedItem from DB..."})
			return
		}
		var allOrderedItems []bson.M
		if err = result.All(c, &allOrderedItems); err != nil {
			log.Fatal(err)
			return
		}
		ctx.JSON(http.StatusOK, allOrderedItems)
	}
}

func GetOrderItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		orderItemId := ctx.Param("order_item_id")
		var orderItem models.OrderItem

		err := orderCollection.FindOne(c, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
		defer cancel()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while listing Ordered Item"})
			return
		}
		ctx.JSON(http.StatusOK, orderItem)
	}
}

func ItemsByOrder(id string) (OrderItemPack []primitive.M, err error) {

}

func CreateOrderItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		//var orderItem models.OrderItem
		var orderItemPack OrderItemPack
		var order models.Order
		if err := ctx.BindJSON(&orderItemPack); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItemsToBeinserted := []interface{}{}
		order.Table_id = *orderItemPack.Table_id
		order_id := OrderItemOrderCreator(order)

		for _, orderItem := range orderItemPack.Order_items {
			orderItem.Order_id = order_id
			validationErr := validate.Struct(orderItem)

			if validationErr != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}

			orderItem.ID = primitive.NewObjectID()
			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Order_item_id = orderItem.ID.Hex()
			var num = ToFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num
			orderItemsToBeinserted = append(orderItemsToBeinserted, orderItem)
		}
		insertedOrderItems, err := orderItemCollection.InsertMany(c, orderItemsToBeinserted)

		if err != nil {
			log.Fatal(err)
		}
		defer cancel()
		ctx.JSON(http.StatusOK, insertedOrderItems)
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItem models.OrderItem
		orderItemId := ctx.Param("order_item_id")
		filter := bson.M{"order_item_id": orderItemId}
		var updateObj primitive.D

		if orderItem.Unit_price != nil {
			updateObj = append(updateObj, bson.E{"unit_price": *&orderItem.Unit_price})
		}

		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{"quantity": *orderItem.Quantity})
		}

		if orderItem.Food_id != nil {
			updateObj = append(updateObj, bson.E{"food_id": *orderItem.Food_id})
		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at": orderItem.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, updateErr := orderItemCollection.UpdateOne(
			c,
			filter,
			bson.D{
				{"$st", updateObj},
			},
			&opt,
		)

		if updateErr != nil {
			msg := fmt.Sprintf("Order Item update failed")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		ctx.JSON(http.StatusOK, result)
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var orderItemId = ctx.Param("order_item_id")

		allOrderedItems, err := ItemsByOrder(orderItemId)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error Occured while Listing Order Items by Order"})
		}

		ctx.JSON(http.StatusOK, allOrderedItems)

	}
}
