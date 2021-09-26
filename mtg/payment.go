package mtg

import (
	"context"
	"log"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/uuid"
	"github.com/shopspring/decimal"
)

var (
	threshold uint8 = 1
)

func MTG_payment_test(c *mixin.Client, access_token, assetID, amount, memo string) string {
	ctx := mixin.WithMixinNetHost(context.Background(), mixin.RandomMixinNetHost())
	user, err := mixin.UserMe(ctx, access_token) // 新建机器人实例
	if err != nil {
		log.Panicln("err:", err, access_token)
	}

	members := []string{c.ClientID, user.UserID} // 门限签名的“分母”名单
	amount_decimal, _ := decimal.NewFromString(amount)
	input := mixin.TransferInput{
		AssetID: assetID,
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
	return payment.CodeID // CodeID 可以组成 mixin://codes/[CodeID] 格式的 scheme url 用以唤醒支付页面
}

func MTG_sign_test(c *mixin.Client, access_token, assetID, memo, pin string) string {
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

		// TODO: 整理筛选 outputs 的方法
		for _, output := range outputs {
			offset = output.UpdatedAt
			if (output.AssetID == assetID) && (output.State == mixin.UTXOStateUnspent) { // 判断币种是否匹配
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
