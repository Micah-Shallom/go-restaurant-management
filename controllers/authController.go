package controllers

import "github.com/gin-gonic/gin"

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func HashPassword(password string) string {
	return ""
}

func VerifyPassword(userPassword string, providePassword string)(bool, string) {
	return true, ""
}