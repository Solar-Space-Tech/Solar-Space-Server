package mtg

import (
	"context"
	// "encoding/json"
	// "sync"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
	// "github.com/jinzhu/gorm"
	// "github.com/jmoiron/sqlx/types"
)

type walletService struct {
	client    *mixin.Client
	members   []string
	threshold uint8
	pin       string
}

func (s *walletService) Pull(ctx context.Context, offset time.Time, limit int) ([]*core.Output, error) {
	outputs, err := s.client.ReadMultisigOutputs(ctx, s.members, s.threshold, offset, limit)
	if err != nil {
		return nil, err
	}

	results := make([]*core.Output, 0, len(outputs))
	for _, output := range outputs {
		result := convertToOutput(output)
		results = append(results, result)
	}

	return results, nil
}

// type walletStore struct {
// 	db   *gorm.DB
// 	once sync.Once
// }

// func save(db *gorm.DB, output *core.Output, ack bool) error {
// 	tx := db.Update().Model(output).Where("trace_id = ?", output.TraceID).Updates(map[string]interface{}{
// 		"state":     output.State,
// 		"signed_tx": output.SignedTx,
// 		"version":   gorm.Expr("version + 1"),
// 	})

// 	if tx.Error != nil {
// 		return tx.Error
// 	}

// 	if tx.RowsAffected == 0 {
// 		if ack {
// 			return db.Update().Create(output).Error
// 		}

// 		return saveRawOutput(db, output)
// 	}

// 	return nil
// }

// type RawOutput struct {
// 	ID        int64          `sql:"PRIMARY_KEY" json:"id"`
// 	CreatedAt int64          `json:"created_at"`
// 	TraceID   string         `sql:"size:36" json:"trace_id"`
// 	Version   int64          `sql:"not null" json:"version"`
// 	Ack       types.BitBool  `sql:"type:bit(1)" json:"ack"`
// 	Data      types.JSONText `sql:"type:TEXT" json:"data"`
// }

// func saveRawOutput(db *gorm.DB, output *core.Output) error {
// 	data, _ := json.Marshal(output)

// 	raw := &RawOutput{
// 		CreatedAt: output.CreatedAt.UnixNano(),
// 		TraceID:   output.TraceID,
// 		Version:   1,
// 		Data:      data,
// 	}

// 	tx := db.Update().Model(raw).
// 		Where("trace_id = ?", raw.TraceID).
// 		Updates(map[string]interface{}{
// 			"data":    raw.Data,
// 			"version": gorm.Expr("version + 1"),
// 		})

// 	if tx.Error != nil {
// 		return tx.Error
// 	}

// 	if tx.RowsAffected == 0 {
// 		return db.Update().Create(raw).Error
// 	}

// 	return nil
// }

func convertToOutput(utxo *mixin.MultisigUTXO) *core.Output {
	return &core.Output{
		CreatedAt:       utxo.CreatedAt,
		UpdatedAt:       utxo.UpdatedAt,
		TraceID:         utxo.UTXOID,
		AssetID:         utxo.AssetID,
		Amount:          utxo.Amount,
		Sender:          utxo.Sender,
		Memo:            utxo.Memo,
		State:           utxo.State,
		TransactionHash: utxo.TransactionHash.String(),
		OutputIndex:     utxo.OutputIndex,
		SignedTx:        utxo.SignedTx,
	}
}