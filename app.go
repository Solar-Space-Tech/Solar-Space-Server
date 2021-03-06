package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"

	db "Solar-Space-Server/db"
	// "Solar-Space-Server/middlewares"
	"Solar-Space-Server/mtg"
	"Solar-Space-Server/util"

	mixin "github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/uuid"

	"net/http"
	"github.com/gin-gonic/gin"

	uuid2 "github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

var html = template.Must(template.New("https").Parse(`
<html>
<head>
  <title>Solar-Space-Tech API</title>
  <script src="/assets/app.js"></script>
</head>
<body>
  <h2 style="color:black;">Welcome to Solar-Space-Tech API Server</h2>
</body>
</html>
`))

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

	// Log Output
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	// 新建机器人实例
	client, err := mixin.NewFromKeystore(&store)
	checkErr(err)
	fmt.Println(client)

	ctx := context.Background()

	// 启动 gin http 服务器
	r := gin.Default()
	// r.Use(middlewares.Cors()) // 非 https 时需要
	r.Static("/assets", "./assets")
	r.SetHTMLTemplate(html)

	r.GET("/", func(c *gin.Context) {
		if pusher := c.Writer.Pusher(); pusher != nil {
			// use pusher.Push() to do server push
			if err := pusher.Push("/assets/app.js", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}
		c.HTML(200, "https", gin.H{
			"status": "success",
		})
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

		// Send user a msg when login successfully to hint user
		cid := mixin.UniqueConversationID(client.ClientID, user.UserID)
		id, _ := uuid.FromString(cid)

		msg := &mixin.MessageRequest{
			ConversationID: cid,
			RecipientID:    user.UserID,
			MessageID:      uuid2.NewV5(id, "login_successful").String(),
			Category:       mixin.MessageCategoryPlainText,
			Data:           base64.StdEncoding.EncodeToString([]byte("登陆成功")),
		}
		// Send the msg
		err = client.SendMessage(ctx, msg)
		checkErr(err)

		//跳转到 return_to,携带 access token
		c.Redirect(http.StatusMovedPermanently, "https://"+return_to+"/#/?access_token="+token)
	})

	r.GET("/api/test/query_uuid_by_phone", func(c *gin.Context) {
		phone := c.Query("phone")
		c.JSON(http.StatusOK, gin.H{
			"uuid": db.Query_uuid_by_phone(phone),
		})
	})

	r.POST("/api/test/deposit_to_multisign", func(c *gin.Context) {
		json := make(map[string]interface{})
		c.BindJSON(&json)
		// access_token := json["access_token"].(string)
		access_token, err := util.GetAccessToken(c)
		checkErr(err)
		var CNB = "965e5c6e-434c-3fa9-b780-c50f43cd955c"

		amount, _ := decimal.NewFromString("10")
		var threshold uint8 = 1
		user, _ := mixin.NewFromAccessToken(access_token).UserMe(ctx)
		members := []string{user.UserID, client.ClientID}
		fmt.Printf("members: %v\n", members)
		traceid := uuid.New()
		timeout := "1321354"
		code_id := mtg.TrustMTGPayment(client, CNB, traceid, timeout, amount, members, threshold)

		c.JSON(http.StatusOK, gin.H{
			"code_id": code_id,
		})
	})

	r.POST("/api/test/withdraw_from_multisign", func(c *gin.Context) {
		json := make(map[string]interface{})
		c.BindJSON(&json)
		// access_token := json["access_token"].(string)
		access_token, err := util.GetAccessToken(c)
		checkErr(err)
		var CNB = "965e5c6e-434c-3fa9-b780-c50f43cd955c"

		code_id := mtg.MTG_sign_test(client, access_token, CNB, "HI,MTG", pcs.Pin)

		c.JSON(http.StatusOK, gin.H{
			"code_id": code_id,
		})
	})

	r.POST("/api/test/encode_memo", func(c *gin.Context) {
		json := make(map[string]interface{})
		c.BindJSON(&json)
		packUuid, _ := uuid2.FromString(json["a"].(string))
		// actionType, _ := strconv.Atoi(c.PostForm("c"))
		encoded_memo := mtg.TrustAction(packUuid, json["t"].(string), json["m"].(string))
		c.JSON(http.StatusOK, gin.H{
			"memo": encoded_memo,
		})
	})

	r.POST("/api/test/decode_memo", func(c *gin.Context) {
		json := make(map[string]interface{})
		c.BindJSON(&json)
		memo := json["memo"].(string)
		decoded_memo := mtg.Unpack_memo(memo)
		fmt.Printf("%+v\n", decoded_memo)
		c.JSON(http.StatusOK, gin.H{
			"a": decoded_memo.AssetID,
			"c": decoded_memo.Type,
			"m": decoded_memo.Amount,
			"t": decoded_memo.Timeout,
		})
	})

	r.RunTLS(":443", "./6395448_api.leaper.one.pem", "./6395448_api.leaper.one.key")
	// r.Run(":8080")
}

func checkErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
