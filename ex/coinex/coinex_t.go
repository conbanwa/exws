package coinex

import (
	"errors"
	"qa3/wstrader/q"
	. "qa3/wstrader/web"
	"strings"
	"sync"

	"github.com/conbanwa/logs"
)

func (coinex *CoinEx) PairArray() (map[string]q.D, map[q.D]q.P, error) {
	return nil, nil, nil
}
func Sym2duo(pair string) q.D {
	parts := strings.Split(pair, "_")
	var res q.D
	if len(parts) == 2 {
		res = q.D{Base: parts[0], Quote: parts[1]}
	} else {
		logs.F("FATAL: DIV ERR!", pair)
	}
	return res
}
func (coinex *CoinEx) Fee() float64 {
	return 0.001
}
func (coinex *CoinEx) PlaceOrders(places [3]q.Order) (orders [3]q.Order, err error) {
	orders = places
	for _, p := range places {
		go func(p q.Order) {
		}(p)
	}
	return
}
func (coinex *CoinEx) Balances() (availables, frozens *sync.Map, err error) {
	availables, frozens = new(sync.Map), new(sync.Map)
	return
}
func (coinex *CoinEx) Test() bool {
	return true
}
func unify(local string) string {
	global := strings.ToUpper(local)
	return global
}
func (coinex *CoinEx) GetAttr() (a q.Attr) {
	a.MultiOrder = false
	return a
}

func (coinex *CoinEx) TradeFee() (map[string]q.TradeFee, error) {
	return nil, nil
}
func (coinex *CoinEx) WithdrawFee() (sf []q.NetworkWithdraw, err error) {

	return
}
func (coinex *CoinEx) OneTicker(d q.D) (ticker q.Bbo, err error) {
	return
}
func (coinex *CoinEx) AllTicker(SymPair map[string]q.D) (mdt *sync.Map, err error) {
	url := "CoinEx.baseurl" + "wait_market_tickers"
	respMap, err := HttpGet(nil, url)
	if err != nil {
		return
	}
	if respMap["status"].(string) == "error" {
		return mdt, errors.New(respMap["err-msg"].(string))
	}
	ticker := mdt
	//ticker.Pair = "currencyPair"
	return ticker, nil
}
