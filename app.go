package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	db "GG-server/db"

	"GG-server/middlewares"
	"GG-server/mixin-api"
	"github.com/gin-gonic/gin"
	mixin "github.com/fox-one/mixin-sdk-go"
)

func main()  {
	r := gin.Default()
	r.Use(middlewares.Cors())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})

	// 接收验证码并且跳转到相应网址
	r.POST("/me", func(c *gin.Context) {
		code := c.PostForm("code")
		return_to := c.PostForm("return_to")
		body := oauth.Oauth(code)

		ctx := context.Background()

		//获取用户信息
		user, err := mixin.UserMe(ctx, body.Get("access_token").MustString())
		if err != nil {
			log.Panicln("err:", err)
		}
		fmt.Println("phone:", user.Phone)

		//TODO: 判断是否为新用户
		db.Insert_mixin(user.Phone, user.UserID)

		c.Redirect(http.StatusMovedPermanently, "http://"+return_to)
	})

	r.GET("/api/test/query_uuid_by_phone", func(c *gin.Context) {
		phone := c.Query("phone")
		c.JSON(http.StatusOK, gin.H{
			"uuid": db.Query_uuid_by_phone(phone),
		})
	})
	r.Run(":80")
}