package okex

import (
	. "github.com/conbanwa/exws/q"
	"github.com/conbanwa/logs"
	"strconv"
	"sync"
)

func (ok *OKEx) Fee() (f float64) {
	f = 0.003
	var response any
	err := ok.DoRequest("GET", "/api/spot/v3/trade_fee?instrument_id=BTC-USDT", "", &response)
	if err != nil {
		return
	}
	// logs.I(response)
	taker, o := response.(map[string]any)["taker"]
	if !o {
		logs.W("okex ", err, response.(map[string]any))
		logs.E("no taker string")
		return
	}
	f, err = strconv.ParseFloat(taker.(string), 64)
	if err != nil {
		logs.F(err, response.(map[string]any))
		return
	}
	return f
}

func (ok *OKEx) TradeFee() (map[string]TradeFee, error) {
	return nil, nil
}
func (ok *OKEx) WithdrawFee() (sf []NetworkWithdraw, err error) {

	return
}
func (ok *OKEx) Balances() (availables, frozens *sync.Map, err error) {
	availables, frozens = new(sync.Map), new(sync.Map)
	return
}
