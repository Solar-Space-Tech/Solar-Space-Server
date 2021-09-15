package mtg

import (
	"reflect"
	"testing"

	"github.com/fox-one/mixin-sdk-go"
	uuid "github.com/satori/go.uuid"
)

func TestPack_memo(t *testing.T) {
	type args struct {
		a string
		c string
		m string
		t string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"cnb", args{"965e5c6e-434c-3fa9-b780-c50f43cd955c", "Trust", "0.5", "1629031344"}, "hKFhxBCWXlxuQ0w/qbeAxQ9DzZVcoWOlVHJ1c3ShbaMwLjWhdKoxNjI5MDMxMzQ0"},
		{"btc", args{"c6d0c728-2624-429b-8e0d-d9d19b6592fa", "Trust", "0.5", "1629031344"}, "hKFhxBDG0McoJiRCm44N2dGbZZL6oWOlVHJ1c3ShbaMwLjWhdKoxNjI5MDMxMzQ0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pack_memo(tt.args.a, tt.args.c, tt.args.m, tt.args.t); got != tt.want {
				t.Errorf("Pack_memo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnpack_memo(t *testing.T) {
	type args struct {
		memo string
	}
	packUuid_cnb, _ := uuid.FromString("965e5c6e-434c-3fa9-b780-c50f43cd955c")
	packUuid_btc, _ := uuid.FromString("c6d0c728-2624-429b-8e0d-d9d19b6592fa")
	tests := []struct {
		name string
		args args
		want Order
	}{
		// TODO: Add test cases.
		{
			"cnb",
			args{"hKFhxBCWXlxuQ0w/qbeAxQ9DzZVcoWOlVHJ1c3ShbaMwLjWhdKoxNjI5MDMxMzQ0"},
			Order{
				AssetID:   packUuid_cnb,
				Action:    "Trust",
				Amount:    "0.5",
				TimeLimit: "1629031344",
			},
		},
		{
			"btc",
			args{"hKFhxBDG0McoJiRCm44N2dGbZZL6oWOlVHJ1c3ShbaMwLjWhdKoxNjI5MDMxMzQ0"},
			Order{
				AssetID:   packUuid_btc,
				Action:    "Trust",
				Amount:    "0.5",
				TimeLimit: "1629031344",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unpack_memo(tt.args.memo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unpack_memo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMTG_payment_test(t *testing.T) {
	type args struct {
		c            *mixin.Client
		access_token string
		assetID      string
		amount       string
		memo         string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MTG_payment_test(tt.args.c, tt.args.access_token, tt.args.assetID, tt.args.amount, tt.args.memo); got != tt.want {
				t.Errorf("MTG_payment_test() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMTG_sign_test(t *testing.T) {
	type args struct {
		c            *mixin.Client
		access_token string
		assetID      string
		memo         string
		pin          string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MTG_sign_test(tt.args.c, tt.args.access_token, tt.args.assetID, tt.args.memo, tt.args.pin); got != tt.want {
				t.Errorf("MTG_sign_test() = %v, want %v", got, tt.want)
			}
		})
	}
}
