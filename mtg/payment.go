package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

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

	const limit = 10
	for utxo == nil {
		outputs, err := client.ReadMultisigOutputs(ctx, members, threshold, offset, limit)
		if err != nil {
			log.Panicf("ReadMultisigOutputs: %v", err)
		}
		for _, output := range outputs {
			offset = output.UpdatedAt
			if hex.EncodeToString(output.TransactionHash[:]) == h.TransactionHash {
				utxo = output
				if(strings.Contains(utxo.asset_id, "CNB") == false || utxo.state == "signed"){
					continue
				}
				break
			}
		}
		if len(outputs) < limit {
			break
		}
	}

	payment, err := client.VerifyPayment(ctx, input)
	if err != nil {
		log.Panicln(err)
	}


	return payment.CodeID
}