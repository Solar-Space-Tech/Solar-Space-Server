package util

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetAccessToken(c *gin.Context) (string, error) {
	access_token := c.GetHeader("Authorization")
	if access_token != "" {
		fmt.Printf("access_token: %v\n", access_token)
		return access_token, nil
	} else {
		return access_token, errors.New("SST: access_token is invaild.")
	}
}
