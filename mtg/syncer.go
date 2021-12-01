package mtg

import (
	"Solar-Space-Server/db"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	// "github.com/fox-one/pando/core"
)

const checkpointKey = "sync_checkpoint"

func New(
	client *mixin.Client,
	mermbers []string,
	threshold uint8,
) *Syncer {
	return &Syncer{
		Client:    client,
		Mermbers:  mermbers,
		Threshold: threshold,
	}
}

type Syncer struct {
	Client    *mixin.Client
	Mermbers  []string
	Threshold uint8
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

	v, err := db.Get_offset(checkpointKey)
	checkErr(err)

	offset := v.Time()

	const limit = 500
	outputs, err := PullUTXOs(ctx, w.Client, w.Mermbers, w.Threshold, offset, limit)
	checkErr(err)

	if len(outputs) == 0 {
		return errors.New("EOF")
	}

	// log.Debugln("walletz.Pull", len(outputs), "outputs")

	nextOffset := outputs[len(outputs)-1].UpdatedAt
	end := len(outputs) < limit

	fmt.Printf("nextOffset: %v\n", nextOffset)
	fmt.Printf("end: %v\n", end)
	// TODO: Store UTXO locally via wallets
	
	//TODO: Store nextOffset via property
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
