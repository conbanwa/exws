package q

import (
	"fmt"
	"qa3/wstrader/cons"
	"strings"
)

type Order struct {
	Price        float64
	Amount       float64
	SAmount      string
	OrderID2     string
	Symbol       string
	Limit        bool
	Sell         bool
	Real         bool
	Err          error
	OrderID      int //deprecated
	Side         cons.TradeSide
	Type         string //limit / market
	OrderType    int    //0:default,1:gtc,2:fok,3:ioc
	OrderTime    int    // create  timestamp
	FinishedTime int64  //finished timestamp
	Currency     cons.CurrencyPair
	Cid          string //客户端自定义ID
	I            int
	Status       cons.TradeStatus
	AvgPrice     float64
	DealAmount   float64
	Fee          float64
}

type Place struct {
	Price, Amount   float64
	SAmount, Symbol string
	Sell            bool
	Err             error
}

func (o Order) Summarize() Place {
	return Place{
		Price:   o.Price,
		Amount:  o.Amount,
		SAmount: o.SAmount,
		Symbol:  o.Symbol,
		Sell:    o.Sell,
		Err:     o.Err,
	}
}

func (o Order) Speak() string {
	if o.Amount <= 0 {
		return "\n"
	}
	return fmt.Sprintf("%+v\n", o)
}

func (o Order) HasErrPrefix(suffix string) bool {
	return strings.HasPrefix(o.Err.Error(), suffix)
}
