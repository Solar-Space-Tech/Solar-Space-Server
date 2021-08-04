package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"encoding/json"
	"net/http"

	db "GG-server/db"

	"GG-server/middlewares"
	"GG-server/mtg"
	"github.com/gin-gonic/gin"
	mixin "github.com/fox-one/mixin-sdk-go"
)


func main()  {
	// 读取配置文件
	f, err := os.Open("./keystore.json")
	if err != nil {
		log.Panicln(err)
	}

	var (
		pin = "" // 配置 PIN
		client_secret  = "" // 配置应用密钥
		store mixin.Keystore
	)
	if err := json.NewDecoder(f).Decode(&store); err != nil {
		log.Panicln(err)
	}

	client, err := mixin.NewFromKeystore(&store)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(client)

	ctx := context.Background()

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
		// body := oauth.Oauth(code)
		token, _, err := mixin.AuthorizeToken(ctx, store.ClientID, client_secret, code, "")
		if err != nil {
			log.Printf("AuthorizeToken: %v", err)
		}
		ctx := context.Background()

		//获取用户信息
		user, err := mixin.UserMe(ctx, token)
		if err != nil {
			log.Panicln("err:", err)
		}
		fmt.Println("phone:", user.Phone)

		//TODO: 判断是否为新用户
		db.Insert_mixin(user.Phone, user.UserID)

		//跳转到 return_to,携带 access token
		c.Redirect(http.StatusMovedPermanently, "http://"+return_to+"?access_token="+token)
	})

	r.GET("/api/test/query_uuid_by_phone", func(c *gin.Context) {
		phone := c.Query("phone")
		c.JSON(http.StatusOK, gin.H{
			"uuid": db.Query_uuid_by_phone(phone),
		})
	})

	r.POST("/api/test/deposit_to_multisign", func(c *gin.Context) {
		access_token := c.PostForm("access_token")
		var CNB = "965e5c6e-434c-3fa9-b780-c50f43cd955c"

		code_id := mtg.Gen_multisig_payment(client, access_token, CNB, "10", "HI,MTG")

		c.JSON(http.StatusOK, gin.H{
			"code_id": code_id,
		})
	})

	r.POST("/api/test/withdraw_from_multisign", func(c *gin.Context) {
		access_token := c.PostForm("access_token")
		var CNB = "965e5c6e-434c-3fa9-b780-c50f43cd955c"

		code_id := mtg.Sign_mtg_test(client, access_token, CNB, "HI,MTG", pin)

		c.JSON(http.StatusOK, gin.H{
			"code_id": code_id,
		})
	})
	r.Run(":8080")
}