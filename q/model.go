package q

import (
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/stat/zelo"
)

var log = zelo.Writer

type Trade struct {
	Tid    int64             `json:"tid"`
	Type   cons.TradeSide    `json:"type"`
	Amount float64           `json:"amount,string"`
	Price  float64           `json:"price,string"`
	Date   int64             `json:"date_ms"`
	Pair   cons.CurrencyPair `json:"omitempty"`
}

type TradeFee struct {
	MakerFee, TakerFee float64
}

type NetworkWithdraw struct {
	AddressRegex            string `json:"addressRegex"`
	Coin                    string `json:"coin"`
	DepositDesc             string `json:"depositDesc"` // shown only when "depositEnable" is false.
	DepositEnable           bool   `json:"depositEnable"`
	FeeDollar               float64
	IsDefault               bool    `json:"isDefault"`
	MemoRegex               string  `json:"memoRegex"`
	MinConfirm              float64 `json:"minConfirm"` // min number for balance confirmation
	Name                    string  `json:"name"`
	Network                 string  `json:"network"`
	ResetAddressStatus      bool    `json:"resetAddressStatus"`
	SpecialTips             string  `json:"specialTips"`
	UnLockConfirm           float64 `json:"unLockConfirm"` // confirmation number for balance unlock
	WithdrawDesc            string  `json:"withdrawDesc"`  // shown only when "withdrawEnable" is false.
	WithdrawEnable          bool    `json:"withdrawEnable"`
	WithdrawIntegerMultiple float64 `json:"withdrawIntegerMultiple"`
	ExFee                   float64
	ExMin                   float64
	Fee                     float64 `json:"withdrawFee"`
	Max                     float64 `json:"withdrawMax"`
	Min                     float64 `json:"withdrawMin"`
	SameAddress             bool    `json:"sameAddress"` // If the coin needs to provide memo to withdraw
	EstimatedArrivalTime    float64 `json:"estimatedArrivalTime"`
	Busy                    bool    `json:"busy"`
}
