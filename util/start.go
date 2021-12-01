package util

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fox-one/mixin-sdk-go"
)

type MixinBot struct{
	Store	mixin.Keystore `json:"store,omitempty"`
	Pin     string `json:"pin,omitempty"`
	Client_secret string `json:"client_secret,omitempty"`
}

func StartMixin() (MixinBot, error) {
	var r MixinBot
	// 读取配置文件
	f_keystore, err := os.Open("./keystore.json")
	CheckErr(err)
	if err != nil {
		return r, err
	}
	f_pcs, err := os.Open("./pin&client_secret.json")
	CheckErr(err)
	if err != nil {
		return r, err
	}
	if err := json.NewDecoder(f_pcs).Decode(&r); err != nil {
		log.Panicln(err)
		return r, err
	}
	if err := json.NewDecoder(f_keystore).Decode(&r.Store); err != nil {
		log.Panicln(err)
		return r, err
	}

	return r, nil
}