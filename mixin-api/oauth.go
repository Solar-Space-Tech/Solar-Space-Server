package oauth

import (
	"bytes"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Oauth(code string) *simplejson.Json {
	data := make(map[string]interface{})
	data["client_id"] = "64abce35-ad54-4828-9e87-b2f46148b0ad"
	data["code"] = code
	data["client_secret"] = "cf2dba3f96acf276a48cda3dce3e209a9d1b3a2f26a0729927c7e37ea718457b"
	bytesData, _ := json.Marshal(data)
	fmt.Println("data:", data)
	req, err := http.Post("https://mixin-api.zeromesh.net/oauth/token", "application/json", bytes.NewBuffer(bytesData))
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	js, _ := simplejson.NewJson([]byte(sb))
	return js.Get("data")
}
