package main

import (
	"context"
	"fmt"
	"log"
	"flag"
	"os"
	"encoding/json"
	"net/http"

	db "GG-server/db"

	"GG-server/middlewares"
	"GG-server/mixin-api"
	"github.com/gin-gonic/gin"
	mixin "github.com/fox-one/mixin-sdk-go"
)


func main()  {
	flag.Parse()

	f, err := os.Open("./keystore.json")
	if err != nil {
		log.Panicln(err)
	}

	var store mixin.Keystore
	if err := json.NewDecoder(f).Decode(&store); err != nil {
		log.Panicln(err)
	}

	client, err := mixin.NewFromKeystore(&store)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(client)

	r := gin.Default()
	r.Use(middlewares.Cors())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
		fmt.Println(client)
	})

	// 接收验证码并且跳转到相应网址
	r.GET("/me", func(c *gin.Context) {
		code := c.Query("code")
		return_to := c.Query("return_to")
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

		//跳转到 return_to,携带 access token
		c.Redirect(http.StatusMovedPermanently, "http://"+return_to+"?access_token="+body.Get("access_token").MustString())
	})

	r.GET("/api/test/query_uuid_by_phone", func(c *gin.Context) {
		phone := c.Query("phone")
		c.JSON(http.StatusOK, gin.H{
			"uuid": db.Query_uuid_by_phone(phone),
		})
	})
	r.Run(":80")
}