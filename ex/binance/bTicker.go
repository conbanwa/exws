package binance

import (
	"encoding/json"
	"errors"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
	"github.com/conbanwa/wstrader/stat/zelo"
	"github.com/conbanwa/wstrader/web"
	"math"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/conbanwa/num"
)

var log = zelo.Writer.With().Str("ex", cons.BINANCE).Logger()

func (bn *Binance) Test() bool {
	return true
}
func (bn *Binance) PairArray() (map[string]q.D, map[q.D]q.P, error) {
	if bn.ExchangeInfo == nil {
		var err error
		bn.ExchangeInfo, err = bn.GetExchangeInfo()
		zelo.PanicOnErr(err).Send()
	}
	var sm = map[string]q.D{}
	var ps = map[q.D]q.P{}
	for _, v := range bn.ExchangeInfo.Symbols {
		if v.Status != "TRADING" || !v.IsSpotTradingAllowed {
			continue
		}
		sm[v.Symbol] = q.D{Base: v.BaseAsset, Quote: v.QuoteAsset}
		ps[sm[v.Symbol]] = q.P{
			Base:         v.GetBaseStep(),
			BaseLimit:    v.GetBaseStep(),
			Quote:        math.Pow10(-v.QuoteAssetPrecision),
			Price:        v.GetPriceStep(),
			MinBase:      v.GetMinBase(),
			MinBaseLimit: v.GetMinBase(),
			MinQuote:     v.GetMinQuote(),
			TakerNett:    1 - bn.Fee(),
		}
		if v.BaseAssetPrecision != v.BaseCommissionPrecision {
			panic(v)
		}
		if v.QuoteCommissionPrecision != v.QuotePrecision || v.QuoteCommissionPrecision != v.QuoteAssetPrecision {
			panic(v)
		}
	}
	zelo.Assert(len(sm) != 0, true).Msgf("%v", bn.ExchangeInfo.Symbols)
	return sm, ps, nil
}

type BookTicker struct {
	Symbol   string  `json:"symbol"`
	BidPrice float64 `json:"bidPrice,string"`
	BidQty   float64 `json:"bidQty,string"`
	AskPrice float64 `json:"askPrice,string"`
	AskQty   float64 `json:"askQty,string"`
}

func valid(v BookTicker) bool {
	return v.AskPrice > 0 && v.AskQty > 0 && v.BidPrice > 0 && v.BidQty > 0
}

/*
*
  - Level	30d Trade	BNB Balance	  	Maker / Taker	Maker / TakerBNB 25% off
    VIP 0	< 50 BTC	or	≥ 0 BNB 	0.1000% / 0.1000%	0.0750% / 0.0750%
    VIP 1	≥ 50 BTC	and	≥ 50 BNB	0.0900% / 0.1000%	0.0675% / 0.0750%
    VIP 2	≥ 500 BTC	and	≥ 200 BNB	0.0800% / 0.1000%	0.0600% / 0.0750%
    VIP 3	≥ 1500 BTC	and	≥ 500 BNB	0.0700% / 0.1000%	0.0525% / 0.0750%
*/
func (bn *Binance) Fee() float64 {
	return 0.00075
}
func (bn *Binance) PlaceOrders(places [3]q.Order) (orders [3]q.Order, err error) {
	var wg sync.WaitGroup
	for i, p := range places {
		if p.Amount <= 0 {
			continue
		}
		wg.Add(1)
		go func(p q.Order, i int) {
			defer wg.Done()
			orders[i] = bn.place(p)
			if orders[i].Err == nil {
				return
			}
			log.Error().Err(orders[i].Err).Send()
			log.Error().Float64("price", p.Price).Str("symbol", p.Symbol).Float64("amount:", p.Amount).Send()
			if orders[i].HasErrPrefix("EOF") {
				for j := 0; j < 5; j++ {
					orders[i] = bn.place(p)
					if orders[i].Err == nil || !orders[i].HasErrPrefix("EOF") {
						break
					}
				}
			}
		}(p, i)
	}
	wg.Wait()
	return
}
func (bn *Binance) place(p q.Order) q.Order {
	test := ""
	if !p.Real {
		test = "/test"
	}
	orderSide := "BUY"
	if p.Sell {
		orderSide = "SELL"
	}
	orderType := "MARKET"
	if p.Limit {
		orderType = "LIMIT"
	}
	path := bn.apiV3 + OrderUri + test
	params := url.Values{}
	params.Set("symbol", p.Symbol)
	params.Set("side", orderSide)
	params.Set("type", orderType)
	// params.Set("newOrderRespType", "FULL")
	params.Set("quantity", p.SAmount)
	if p.Limit {
		params.Set("timeInForce", "IOC") //FOK")GTC
		params.Set("price", strconv.FormatFloat(p.Price, 'f', -1, 64))
		// } else {
		// params.Set("quoteOrderQty", (p.SAmount))
		// params.Set("newOrderRespType", "FULL") //RESULT ACK
	}
	bn.buildParamsSigned(&params)
	respMap := make(map[string]any)
	if resp, err := web.HttpPostForm2(bn.httpClient, path, params, bn.header()); err != nil {
		p.Err = err
		return p
	} else {
		p.Err = json.Unmarshal(resp, &respMap)
	}
	if len(respMap) == 0 {
		if !p.Real {
			p.DealAmount = p.Amount
			return p
		}
		p.Err = errors.New("order err: empty map")
	}
	if p.Err != nil {
		log.Error().Err(p.Err).Any("resp", respMap).Msgf("order err: %+v", p)
		return p
	}
	log.Println(respMap)

	orderId := respMap["orderId"].(int)
	dealAmount := num.ToFloat64(respMap["executedQty"])
	if cummulativeQuoteQty := num.ToFloat64(respMap["cummulativeQuoteQty"]); cummulativeQuoteQty > 0 && dealAmount > 0 {
		p.AvgPrice = cummulativeQuoteQty / dealAmount
	}
	if !p.Real {
		dealAmount = p.Amount
	}
	p.OrderID = orderId
	p.OrderID2 = strconv.Itoa(orderId)
	p.DealAmount = dealAmount
	p.OrderTime = num.ToInt[int](respMap["transactTime"])
	// p.Status =    bn.parseOrderStatus(respMap["status"]),
	return p
}
func unify(local string) string {
	global := strings.ToUpper(local)
	return global
}
func (bn *Binance) GetAttr() (a q.Attr) {
	a.MultiOrder = false
	return a
}

func (bn *Binance) OneTicker(d q.D) (ticker q.Bbo, err error) {
	return
}
func (bn *Binance) AllTicker(SymPair map[string]q.D) (*sync.Map, error) {
	resp, err := web.HttpGet5(bn.httpClient, ApiV3+"ticker/bookTicker", nil)
	if err != nil {
		log.Error().Err(err).Bytes("resp", resp).Send()
		return nil, err
	}
	var bodyDataMap []BookTicker
	err = json.Unmarshal(resp, &bodyDataMap)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}
	var mdt sync.Map
	for _, v := range bodyDataMap {
		if !valid(v) {
			continue
		}
		pair, ok := SymPair[v.Symbol]
		if !ok {
			log.Debug().Msg(v.Symbol + "not exit")
			continue
		}
		var ticker = q.Bbo{
			Pair:    v.Symbol,
			Bid:     v.BidPrice,
			BidSize: v.BidQty,
			Ask:     v.AskPrice,
			AskSize: v.AskQty,
		}
		if ticker.Valid() {
			mdt.Store(pair, ticker)
		}
	}
	return &mdt, nil
}
