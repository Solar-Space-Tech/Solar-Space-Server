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

type Payment struct {
	AssetID    string          `json:"asset_id,omitempty"`
	Amount     decimal.Decimal `json:"amount,omitempty"`
	ActionType string          `json:"type,omitempty"`
	Receivers  []string        `json:"receivers,omitempty"`
	Threshold  uint8           `json:"threshold,omitempty"`
	TraceID    string          `json:"trace_id,omitempty"`
	Timeout    string          `json:"time_out,omitempty"`
}

func (p Payment) MTG_payment(c *mixin.Client) string {
	ctx := mixin.WithMixinNetHost(context.Background(), mixin.RandomMixinNetHost())
	var memo string
	assetid, _ := uuid.FromString(p.AssetID)
	switch p.ActionType {
	case "Trust":
		memo = TrustAction(assetid, p.Timeout, p.Amount.String())
		// TODO: case...
	}
	input := mixin.TransferInput{
		AssetID: p.AssetID,
		Amount:  p.Amount,
		TraceID: p.TraceID,
		Memo:    memo,
		OpponentMultisig: struct {
			Receivers []string `json:"receivers,omitempty"`
			Threshold uint8    `json:"threshold,omitempty"`
		}{
			Receivers: p.Receivers,
			Threshold: p.Threshold,
		},
	}

	payment, err := c.VerifyPayment(ctx, input)
	if err != nil {
		log.Panicln(err)
	}
	return payment.CodeID
}

func TrustMTGPayment(c *mixin.Client, asset_id, trace_id, time_out string, amount_decimal decimal.Decimal, receivers []string, threshold uint8) string {
	payment := Payment{
		AssetID:    asset_id,
		Amount:     amount_decimal,
		ActionType: "Trust",
		Receivers:  receivers,
		Threshold:  threshold,
		TraceID:    trace_id,
		Timeout:    time_out,
	}
	code_id := payment.MTG_payment(c)
	return code_id
}

func MTG_payment_test(c *mixin.Client, access_token, assetID, amount, memo string) string {
	ctx := mixin.WithMixinNetHost(context.Background(), mixin.RandomMixinNetHost())
	user, err := mixin.UserMe(ctx, access_token) // 新建用户实例
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
