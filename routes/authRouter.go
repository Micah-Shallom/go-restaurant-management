package routes

import (
	"github.com/Micah-Shallom/modules/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("/signup/", controllers.SignUp())
	incomingRoutes.POST("/login/", controllers.Login())
}