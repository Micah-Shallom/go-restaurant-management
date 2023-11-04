package routes

import (
	"github.com/Micah-Shallom/modules/controllers"
	"github.com/Micah-Shallom/modules/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	middleware.Authenticate()
	incomingRoutes.GET("/users/", controllers.GetUsers())
	incomingRoutes.GET("/users/:userid", controllers.GetUser())
}