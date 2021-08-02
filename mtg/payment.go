package mtg

import (
	"context"
	"log"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/uuid"
	"github.com/shopspring/decimal"
)

var (
	members = []string{}
	threshold uint8 = 2
)

func Gen_multisig_payment(store mixin.Keystore, client_id, assetID, amount, memo string) string{

	members = append(members, store.ClientID, client_id)

	client, err := mixin.NewFromKeystore(&store)
	if err != nil {
		log.Panicln(err)
	}
	
	ctx := mixin.WithMixinNetHost(context.Background(), mixin.RandomMixinNetHost())

	amount_decimal, _ := decimal.NewFromString(amount)
	input := mixin.TransferInput{
		AssetID: assetID, 
		// AssetID: "965e5c6e-434c-3fa9-b780-c50f43cd955c",
		Amount:  amount_decimal, 
		TraceID: uuid.New(),
		Memo:    memo,
		OpponentMultisig: struct {
			Receivers []string `json:"receivers,omitempty"`
			Threshold uint8    `json:"threshold,omitempty"`
		}{
			Receivers: members,
			Threshold: threshold,
		},
	}

	payment, err := client.VerifyPayment(ctx, input)
	if err != nil {
		log.Panicln(err)
	}
	return payment.CodeID
}
