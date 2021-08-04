package mtg

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/uuid"
	"github.com/shopspring/decimal"
)

var (
	members = []string{}
	threshold uint8 = 2
)

func Gen_multisig_payment(c *mixin.Client, client_id, assetID, amount, memo string) (string) {

	members = append(members, c.ClientID, client_id)

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

	payment, err := c.VerifyPayment(ctx, input)
	if err != nil {
		log.Panicln(err)
	}
	return payment.CodeID
}

func Sign_mtg_test(c *mixin.Client, access_token , assetID, memo, pin string) (string) {
	ctx := mixin.WithMixinNetHost(context.Background(), mixin.RandomMixinNetHost())
	// 读取用户
	user := mixin.NewFromAccessToken(access_token)

	var (
		utxo   *mixin.MultisigUTXO
		offset time.Time
	)
	const limit = 10
	for utxo == nil {
		outputs, err := user.ReadMultisigOutputs(ctx, members, threshold, offset, limit)
		if err != nil {
			log.Panicf("ReadMultisigOutputs: %v", err)
		}

		for _, output := range outputs {
			offset = output.UpdatedAt
			if hex.EncodeToString([]byte(output.AssetID)) == assetID {
				utxo = output
				break
			}
		}
		if len(outputs) < limit {
			break
		}
	}

	if utxo == nil {
		log.Panicln("No Unspent UTXO")
	}

	amount := utxo.Amount.Truncate(8)

	tx, err := c.MakeMultisigTransaction(ctx, &mixin.TransactionInput{
	Memo:   "multisig test",
	Inputs: []*mixin.MultisigUTXO{utxo},
	Outputs: []mixin.TransactionOutput{
		{
			Receivers: []string{user.ClientID}, // 用户收币
			Threshold: 1,
			Amount:    amount,
		},
	},
	Hint: uuid.New(),
	})

	if err != nil {
		log.Panicf("MakeMultisigTransaction: %v", err)
	}

	raw, err := tx.DumpTransaction()
	if err != nil {
		log.Panicf("DumpTransaction: %v", err)
	}

	// 机器人，生成签名请求
	req, err := c.CreateMultisig(ctx, mixin.MultisigActionSign, raw)
	if err != nil {
		log.Panicf("CreateMultisig: sign %v", err)
	}
	// 机器人，签名
	req, err = c.SignMultisig(ctx, req.RequestID, pin)
	if err != nil {
		log.Panicf("CreateMultisig: %v", err)
	}
	fmt.Println(req)

	re, err := c.CreateMultisig(ctx, mixin.MultisigActionSign, raw)
	if err != nil {
		log.Panicf("CreateMultisig: sign %v", err)
	}

	return re.CodeID
}