package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kwamekyeimonies/restaurant_management_system_backend/controllers"
)

func MenuRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/menus", controllers.GetMenus())
	incomingRoutes.GET("/menus/:menu_id", controllers.GetMenu())
	incomingRoutes.POST("/menus", controllers.CreateMenu())
	incomingRoutes.PATCH("/menus/menu_id", controllers.UpdateMenu())
}
