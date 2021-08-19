package mtg

import (
	"context"
	"log"
	"time"
	"encoding/base64"
	uuid2 "github.com/satori/go.uuid"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/vmihailenco/msgpack"
	"github.com/fox-one/pkg/uuid"
	"github.com/shopspring/decimal"
)

var (
	threshold uint8 = 1
)

type Order struct {
	A uuid2.UUID `json:"a,omitempty" msgpack:"a,omitempty`
	C string `json:"c,omitempty" msgpack:"c,omitempty`
	M string `json:"m,omitempty" msgpack:"m,omitempty`
	T string `json:"t,omitempty" msgpack:"t,omitempty`

}

func Pack_memo(a, c, m, t string) string {
	packUuid, _ := uuid2.FromString(a)
	pack, _ := msgpack.Marshal(Order{A: packUuid, C: c, M: m, T: t,})
	memo := base64.StdEncoding.EncodeToString(pack)
	return memo
}
func Unpack_memo(memo string) Order {
	// 解码 memo
	parsedpack, _ := base64.StdEncoding.DecodeString(memo)
	order_memo := Order{}
	err := msgpack.Unmarshal(parsedpack, &order_memo)
	if err != nil {
		log.Panicln(err)
	}
	// TODO: 判断 memo 是否有效
	// TODO: 如果有效则存入数据库

	return order_memo
}

func MTG_payment_test(c *mixin.Client, access_token, assetID, amount, memo string) (string) {
	// log.Panicf("-|-|access_token|-|-/n<%s>\n",access_token)
	ctx := mixin.WithMixinNetHost(context.Background(), mixin.RandomMixinNetHost())
	user, err := mixin.UserMe(ctx, access_token)
	if err != nil {
		log.Panicln("err:", err, access_token)
	}

	members := []string{c.ClientID, user.UserID}

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


func MTG_sing_test(c *mixin.Client, access_token , assetID, memo, pin string) (string) {
	// log.Panicf("-|-|access_token|-|-/n<%s>\n",access_token)
	ctx := mixin.WithMixinNetHost(context.Background(), mixin.RandomMixinNetHost())
	// 读取用户
	client := mixin.NewFromAccessToken(access_token)

	user, err := client.UserMe(ctx)
	if err != nil {
		log.Panicln("err:", err)
	}

	members := []string{c.ClientID, user.UserID}

	var (
		utxo   *mixin.MultisigUTXO
		offset time.Time
	)
	const limit = 10
	for utxo == nil {
		outputs, err := client.ReadMultisigOutputs(ctx, members, threshold, offset, limit)
		if err != nil {
			log.Panicf("ReadMultisigOutputs: %v", err)
		}

		for _, output := range outputs {
			offset = output.UpdatedAt
			if (output.AssetID == assetID) && (output.State == mixin.UTXOStateUnspent) { // 判断币种
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

	amount := utxo.Amount.Div(decimal.NewFromFloat(2)).Truncate(8)

	tx, err := c.MakeMultisigTransaction(ctx, &mixin.TransactionInput{
	Memo:   "multisig test",
	Inputs: []*mixin.MultisigUTXO{utxo},
	Outputs: []mixin.TransactionOutput{
		{
			Receivers: []string{user.UserID}, // 用户收币
			Threshold: threshold,
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

	//用户，生成签名请求
	re, err := client.CreateMultisig(ctx, mixin.MultisigActionSign, raw)
	if err != nil {
		log.Panicf("CreateMultisig: sign %v", err)
	}

	return re.CodeID
}
