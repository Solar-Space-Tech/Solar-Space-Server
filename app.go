package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	db "GG-server/db"

	"GG-server/middlewares"
	"GG-server/mtg"
	mixin "github.com/fox-one/mixin-sdk-go"
	"github.com/gin-gonic/gin"
	uuid2 "github.com/satori/go.uuid"
)

func main() {
	// 读取配置文件
	f_keystore, err := os.Open("./keystore.json")
	checkErr(err)
	f_pcs, err := os.Open("./pin&client_secret.json")
	checkErr(err)

	var (
		pcs struct {
			Pin           string `json:"pin"`
			Client_secret string `json:"client_secret"`
		}
		store mixin.Keystore
	)
	if err := json.NewDecoder(f_pcs).Decode(&pcs); err != nil {
		log.Panicln(err)
	}
	if err := json.NewDecoder(f_keystore).Decode(&store); err != nil {
		log.Panicln(err)
	}

	// 新建机器人实例
	client, err := mixin.NewFromKeystore(&store)
	checkErr(err)
	fmt.Println(client)

	ctx := context.Background()

	// 启动 gin http 服务器
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
		token, _, err := mixin.AuthorizeToken(ctx, store.ClientID, pcs.Client_secret, code, "")
		if err != nil {
			log.Printf("AuthorizeToken: %v", err)
		}
		ctx := context.Background()

		//获取用户信息
		user, err := mixin.UserMe(ctx, token)
		checkErr(err)
		fmt.Println("phone:", user.Phone)

		// 判断是否为新用户
		if !db.If_old_user(user.UserID, user.Phone) {
			db.Insert_mixin(user.Phone, user.UserID, user.FullName)
		}

		//跳转到 return_to,携带 access token
		c.Redirect(http.StatusMovedPermanently, "http://"+return_to+"/#/?access_token="+token)
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

		code_id := mtg.MTG_payment_test(client, access_token, CNB, "10", "HI,MTG")

		c.JSON(http.StatusOK, gin.H{
			"code_id": code_id,
		})
	})

	r.POST("/api/test/withdraw_from_multisign", func(c *gin.Context) {
		access_token := c.PostForm("access_token")
		var CNB = "965e5c6e-434c-3fa9-b780-c50f43cd955c"

		code_id := mtg.MTG_sign_test(client, access_token, CNB, "HI,MTG", pcs.Pin)

		c.JSON(http.StatusOK, gin.H{
			"code_id": code_id,
		})
	})

	r.POST("/api/test/encode_memo", func(c *gin.Context) {
		packUuid, _ := uuid2.FromString(c.PostForm("a"))
		order := mtg.Action{
			AssetID:   packUuid,
			Action:    c.PostForm("c"),
			Amount:    c.PostForm("m"),
			Timeout: c.PostForm("t"),
		}
		encoded_memo := order.Pack_memo()
		fmt.Printf("%+v\n", encoded_memo)
		c.JSON(http.StatusOK, gin.H{
			"memo": encoded_memo,
		})
	})

	r.POST("/api/test/decode_memo", func(c *gin.Context) {
		memo := c.PostForm("memo")
		decoded_memo := mtg.Unpack_memo(memo)
		fmt.Printf("%+v\n", decoded_memo)
		c.JSON(http.StatusOK, gin.H{
			"a": decoded_memo.AssetID,
			"c": decoded_memo.Action,
			"m": decoded_memo.Amount,
			"t": decoded_memo.Timeout,
		})
	})

	r.Run(":8080")
}

func checkErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
