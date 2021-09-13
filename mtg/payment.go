package mtg

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/uuid"
	uuid2 "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/vmihailenco/msgpack"
)

var (
	threshold uint8 = 1
)

// 多签交易 Memo 规范
type Order struct {
	AssetID   uuid2.UUID `msgpack:"a"`
	Action    string     `msgpack:"c"`
	Amount    string     `msgpack:"m"`
	TimeLimit string     `msgpack:"t"`
}

// 将 Order 经过 mesgpack 打包，再 base64 加密
func Pack_memo(a, c, m, t string) string {
	packUuid, _ := uuid2.FromString(a)
	n := Order{
		AssetID:   packUuid,
		Action:    c,
		Amount:    m,
		TimeLimit: t,
	}
	pack, err := msgpack.Marshal(&n)
	if err != nil {
		log.Panicln(err)
	}
	memo := base64.StdEncoding.EncodeToString(pack)
	return memo
}

// Memo 解码，为 Pack_memo 逆过程
func Unpack_memo(memo string) Order {
	// 解码 memo
	parsedpack, _ := base64.RawURLEncoding.DecodeString(memo)
	fmt.Println(parsedpack)
	order_memo := Order{}
	err := msgpack.Unmarshal(parsedpack, &order_memo)
	if err != nil {
		fmt.Println(err)
	}
	// TODO: 判断 memo 是否有效
	// TODO: 如果有效则存入数据库

	return order_memo
}

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
