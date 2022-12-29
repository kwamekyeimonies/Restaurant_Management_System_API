package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kwamekyeimonies/restaurant_management_system_backend/controllers"
)

func OrderItemRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/orderitems", controllers.GetOrderItems())
	incomingRoutes.GET("/orderitems/:orderitem_id", controllers.GetOrderItem())
	incomingRoutes.POST("/orderitems", controllers.CreateOrderItem())
	incomingRoutes.GET("/oderitems-order/:order_id", controllers.GetOrderItemsByOrder())
	incomingRoutes.PATCH("/orderitems/:orderitem_id", controllers.UpdateOrderItem())
}
