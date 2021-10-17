package mtg

import (
	"context"
	// "encoding/json"
	// "sync"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pando/core"
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