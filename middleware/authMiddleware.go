package middleware

import (
	"fmt"
	"net/http"

	"github.com/Micah-Shallom/modules/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context){
		clientToken := c.Request.Header.Get("token")
		if clientToken == ""{
			c.JSON(http.StatusInternalServerError, gin.H{"error":fmt.Sprintf("No Authorization Header Provided")})
			c.Abort()
			return
		}
		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error":err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("firstname", claims.FirstName)
		c.Set("lastname", claims.LastName)
		c.Set("uid", claims.Uid)
		c.Set("userType", claims.UserType)
		c.Next()
	}
}