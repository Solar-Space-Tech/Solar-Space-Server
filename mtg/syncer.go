package mtg

import (
	"Solar-Space-Server/db"
	"context"
	"errors"
	"fmt"
	"time"

	// "github.com/fox-one/pando/core"
)

const checkpointKey = "sync_checkpoint"

func New(
	// wallets core.WalletStore,
	walletz walletService,
	property db.Property,
) *Syncer {
	return &Syncer{
		// wallets: wallets,
		walletz: walletz,
		property: property,
	}
}

type Syncer struct {
	// wallets  core.WalletStore
	walletz  walletService
	property db.Property
}

func (w *Syncer) Run(ctx context.Context) error {

	dur := time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.run(ctx); err == nil {
				dur = 100 * time.Millisecond
			} else {
				dur = 500 * time.Millisecond
			}
		}
	}
}

func (w *Syncer) run(ctx context.Context) error {

	v, err := w.property.Get_offset(checkpointKey)
	checkErr(err)

	offset := v.Time()

	const limit = 500
	outputs, err := w.walletz.Pull(ctx, offset, limit)
	checkErr(err)

	if len(outputs) == 0 {
		return errors.New("EOF")
	}

	// log.Debugln("walletz.Pull", len(outputs), "outputs")

	nextOffset := outputs[len(outputs)-1].UpdatedAt
	end := len(outputs) < limit

	fmt.Printf("nextOffset: %v\n", nextOffset)
	fmt.Printf("end: %v\n", end)

	return nil
}
