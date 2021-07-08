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
	data["client_id"] = "b31dc198-3375-4f53-bad4-59abb4f30dd9"
	data["code"] = code
	data["client_secret"] = "02481ebfe76922d2ebc800b19738a99f0a59d0fe1d517ab45957d6ba030343be"
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
