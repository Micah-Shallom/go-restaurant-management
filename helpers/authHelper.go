package helpers

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("usertype")
	err = nil 
	if userType != role {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	return err
}

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("usertype")
	uid := c.GetString("uid")
	err = nil

	if userType == "USER" && uid != userId {
		err = errors.New("unauthorized to access this resource")
		return err
	}

	err = CheckUserType(c, userType)
	return err
}