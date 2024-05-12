package bitfinex

import (
	"errors"
	"github.com/conbanwa/num"
	"github.com/conbanwa/wstrader/q"
	. "github.com/conbanwa/wstrader/web"
	"strings"
	"sync"

	"github.com/conbanwa/logs"
)

func (bfx *Bitfinex) PairArray() (sm map[string]q.D, pairs map[q.D]q.P, err error) {
	sm = map[string]q.D{}
	pairs = map[q.D]q.P{}
	const turl = "https://api-pub.bitfinex.com/v2/tickers?symbols=ALL"
	respMap, err := HttpGet3(bfx.httpClient, turl, nil)
	if err != nil {
		return nil, nil, err
	}
	for _, v := range respMap {
		u, ok := v.([]any)
		if !ok {
			panic(v)
			return nil, nil, errors.New("convert")
		}
		parts, ok := u[0].(string)
		if !ok {
			panic(u)
		}
		d, fok := Sym2duo(parts)
		if parts[0] == 'f' {
		} else if fok != "" {
			sm[parts] = d
			pairs[d] = q.P{}
		} else {
			panic(u[0])
			return nil, nil, errors.New("bitfinex not string")
		}
	}
	return
}
func Sym2duo(parts string) (d q.D, ok string) {
	if parts[0] == 't' {
		colon := strings.Index(parts, ":")
		if colon != -1 {
			d = q.D{Base: unify(parts[1:colon]), Quote: unify(parts[colon+1:])}
			ok = parts[1:colon] + parts[colon+1:]
		} else if len(parts) == 7 { //LUNA tALBT:USD
			d = q.D{Base: unify(parts[1:4]), Quote: unify(parts[4:])}
			ok = parts[1:]
		} else {
			logs.W(d)
		}
	}
	return
}
func (bfx *Bitfinex) Fee() float64 {
	var respMap map[string]any
	err := bfx.doAuthenticatedRequest("POST", "summary", map[string]any{}, &respMap)
	if err != nil {
		logs.E(err, respMap)
		return 0.001
	}
	// logs.D("Bitfinex ", respMap)
	return respMap["taker_fee"].(float64)
}
func (bfx *Bitfinex) PlaceOrders(places [3]q.Order) (orders [3]q.Order, err error) {
	orders = places
	for _, p := range places {
		go func(p q.Order) {
		}(p)
	}
	return orders, nil
}
func (bfx *Bitfinex) Balances() (availables, frozens *sync.Map, err error) {
	availables, frozens = new(sync.Map), new(sync.Map)
	var respMap []any
	err = bfx.doAuthenticatedRequest("GET", "balances", map[string]any{}, &respMap)
	if err != nil {
		return
	}
	for _, v := range respMap {
		subacc := v.(map[string]any)
		typeStr := subacc["type"].(string)
		currency := unify(subacc["currency"].(string))
		amount := num.ToFloat64(subacc["amount"])
		available := num.ToFloat64(subacc["available"])
		if typeStr == "exchange" {
			availables.Store(currency, available)
			frozens.Store(currency, amount-available)
		}
		// account := balancemap[typeStr]
		// if account == nil {
		// 	account = new(Account)
		// 	account.SubAccounts = make(map[Currency]SubAccount, 6)
		// }
		// account.NetAsset = amount
		// account.Asset = amount
		// account.SubAccounts[currency] = SubAccount{
		// 	Currency:     currency,
		// 	Amount:       available,
		// 	ForzenAmount: amount - available,
		// 	LoanAmount:   0}
		// balancemap[typeStr] = account
	}
	return
}
func (bfx *Bitfinex) Test() bool {
	return true
}
func unify(local string) string {
	global := strings.ToUpper(local)
	switch global {
	case "ALG":
		return "ALGO"
	case "ATO":
		return "ATOM"
	case "DAT":
		return "DATA"
	case "DOG":
		return "MDOGE"
	case "DSH":
		return "DASH"
	case "EDO":
		return "PNT"
	case "GNT":
		return "GLM"
	case "IOT":
		return "IOTA"
	case "LBT":
		return "LBTC"
	case "MNA":
		return "MANA"
	case "QTM":
		return "QTUM"
	case "RBT":
		return "RBTC"
	case "SNG":
		return "SNGLS"
	case "STJ":
		return "STORJ"
	case "TSD":
		return "TUSD"
	case "UDC":
		return "USDC"
	case "USK":
		return "USDK"
	case "UST":
		return "USDT"
	case "WBT":
		return "WBTC"
	case "YYW":
		return "YOYOW"
	}
	return global
}
func (bfx *Bitfinex) GetAttr() (a q.Attr) {
	a.MultiOrder = false
	return a
}

func (bfx *Bitfinex) TradeFee() (map[string]q.TradeFee, error) {
	return nil, nil
}
func (bfx *Bitfinex) WithdrawFee() (sf []q.NetworkWithdraw, err error) {

	return
}
func (bfx *Bitfinex) OneTicker(d q.D) (ticker q.Bbo, err error) {
	return
}
func (bfx *Bitfinex) AllTicker(SymPair map[string]q.D) (mdt *sync.Map, err error) {
	var PairSym = map[q.D]string{}
	const turl = "https://api-pub.bitfinex.com/v2/tickers?symbols=ALL"
	respMap, err := HttpGet3(bfx.httpClient, turl, nil)
	if err != nil {
		return
	}
	for _, v := range respMap {
		u, ok := v.([]any)
		if !ok {
			panic(v)
			return mdt, errors.New("convert")
		}
		var f [11]float64
		parts, ok := u[0].(string)
		if !ok {
			panic(u)
		}
		d, fok := Sym2duo(parts)
		if parts[0] == 'f' {
		} else if fok != "" {
			var ticker q.Bbo
			for i := 1; i < 5; i++ {
				f[i], ok = u[i].(float64)
				if !ok {
					panic(u[i])
					return mdt, errors.New("not float64")
					// } else if f[i] == 0 {
					// 	logs.E(v, u)
				}
			}
			ticker.Bid = f[1]
			ticker.BidSize = f[2]
			ticker.Ask = f[3]
			ticker.AskSize = f[4]
			if ticker.Valid() {
				mdt.Store(d, ticker)
				PairSym[d] = fok
				// } else {
				// 	logs.E(ticker)
			}
		} else {
			panic(u[0])
			return mdt, errors.New("bitfinex not string")
		}
	}
	// if mdt[Duo{"TERRAUST", "USD"}].Ask == 0 {
	// logs.I(respMap)
	// }
	return
}
