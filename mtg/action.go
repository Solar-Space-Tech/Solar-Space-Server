package mtg

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"github.com/vmihailenco/msgpack"
)

// 多签交易 Memo 规范
type Action struct {
	// action type
	Type string `msgpack:"c,omitemnty"`
	// AssetID is pair quote asset id if base asset will be paid, otherwise this is base asset id.
	AssetID uuid.UUID `msgpack:"a,omitemnty"`
	Amount string `msgpack:"m,omitemnty"`
	Timeout string `msgpack:"t,omitemnty"`
}

func TrustAction(assetID uuid.UUID, timeout, amount string) Action {
	return Action{
		Type:    "Trust",
		AssetID: assetID,
		Timeout: timeout,
		Amount:  amount,
	}
}

// 将 Order 经过 mesgpack 打包，再 base64 加密
func (A Action) Pack_memo() string {
	pack, err := msgpack.Marshal(&A)
	if err != nil {
		log.Panicln(err)
	}
	memo := base64.StdEncoding.EncodeToString(pack)
	return memo
}

// Memo 解码，为 Pack_memo 逆过程
func Unpack_memo(memo string) Action {
	// 解码 memo
	parsedpack, _ := base64.StdEncoding.DecodeString(memo)
	fmt.Println(parsedpack)
	order_memo := Action{}
	err := msgpack.Unmarshal(parsedpack, &order_memo)
	if err != nil {
		fmt.Println(err)
	}
	// TODO: 判断 memo 是否有效
	// TODO: 如果有效则存入数据库

	return order_memo
}
