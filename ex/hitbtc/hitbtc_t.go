package hitbtc

import (
	"github.com/conbanwa/wstrader/q"
	"strings"
	"sync"

	"github.com/conbanwa/logs"
)

func (hitbtc *Hitbtc) PairArray() (sm map[string]q.D, prec map[q.D]q.P, err error) {
	sm = map[string]q.D{}
	prec = map[q.D]q.P{}
	var resp []map[string]any
	err = hitbtc.doRequest("GET", SYMBOLS_URI, &resp)
	if err != nil {
		return
	}
	for _, e := range resp {
		d := q.D{Base: e["baseCurrency"].(string), Quote: e["quoteCurrency"].(string)}
		sm[e["id"].(string)] = d
		prec[d] = q.P{
			// Base:         e["baseCurrency"].(float64),
			// Quote:        e["tick_size"].(float64),
			// Price:        e["baseCurrency"].(float64),
			MinQuote:     e["baseCurrency"].(float64),
			MinBaseLimit: e["baseCurrency"].(float64),
			// MinBase:      minbase,
			// BaseLimit:    e["baseCurrency"].(float64),
			TakerNett: 0.998}
	}
	return
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
func (hitbtc *Hitbtc) Fee() float64 {
	return 0.001
}
func (hitbtc *Hitbtc) PlaceOrders(places [3]q.Order) (orders [3]q.Order, err error) {
	orders = places
	for _, p := range places {
		go func(p q.Order) {
		}(p)
	}
	return
}
func (hitbtc *Hitbtc) Balances() (availables, frozens *sync.Map, err error) {
	availables, frozens = new(sync.Map), new(sync.Map)
	return
}
func (hitbtc *Hitbtc) Test() bool {
	return true
}
func unify(local string) string {
	global := strings.ToUpper(local)
	return global
}
func (hitbtc *Hitbtc) GetAttr() (a q.Attr) {
	a.MultiOrder = false
	return a
}

func (hitbtc *Hitbtc) TradeFee() (map[string]q.TradeFee, error) {
	return nil, nil
}
func (hitbtc *Hitbtc) WithdrawFee() (sf []q.NetworkWithdraw, err error) {

	return
}
func (hitbtc *Hitbtc) OneTicker(d q.D) (ticker q.Bbo, err error) {
	return
}
func (hitbtc *Hitbtc) AllTicker(SymPair map[string]q.D) (mdt *sync.Map, err error) {
	logs.F(hitbtc.GetSymbols())
	return
}
