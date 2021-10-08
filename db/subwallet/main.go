// Generate a bunch of sub-wallet and store them to the db.
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"log"
	"os"
	"time"

	db "Solar-Space-Server/db"

	mixin "github.com/fox-one/mixin-sdk-go"
)

func main() {

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

	client, err := mixin.NewFromKeystore(&store)
	checkErr(err)
	ctx := context.Background()

	for i := 0; i <= 1000; i++ {
		Create_sub_wallet(ctx, *client)
	}

}

func Create_sub_wallet(ctx context.Context, client mixin.Client) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	sub, subStore, err := client.CreateUser(ctx, privateKey, "sub user")
	time.Sleep(1)
	if err != nil {
		log.Printf("CreateUser: %v", err)
		return
	}
	log.Println("create sub user", sub.UserID)
	log.Println("sub user pk", subStore.PrivateKey)
	// set pin
	subClient, _ := mixin.NewFromKeystore(subStore)
	if err := subClient.ModifyPin(ctx, "", "000000"); err != nil {
		log.Printf("ModifyPin: %v", err)
		return
	}

	time.Sleep(1)
	db.Insert_subWallet(subStore)
}

func checkErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
