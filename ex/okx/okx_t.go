package okx

import (
	"fmt"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
	"github.com/conbanwa/wstrader/stat/zelo"
	. "github.com/conbanwa/wstrader/web"
	"strings"
	"sync"
)

var log = zelo.Writer.With().Str("ex", cons.OKEX).Logger()

type TickerV5Response struct {
	Code int        `json:"code,string"`
	Msg  string     `json:"msg"`
	Data []TickerV5 `json:"data"`
}

func (ok *OKX) PairArray() (sm map[string]q.D, ps map[q.D]q.P, err error) {
	sm = map[string]q.D{}
	ps = map[q.D]q.P{}
	urlPath := ok.config.Endpoint + "/api/v5/public/instruments?instType=SPOT"
	type instrumentsResponse struct {
		Code int    `json:"code,string"`
		Msg  string `json:"msg"`
		Data []struct {
			Alias        string  `json:"alias"`
			BaseCcy      string  `json:"baseCcy"`
			Category     string  `json:"category"`
			CtMult       string  `json:"ctMult"`
			CtType       string  `json:"ctType"`
			CtVal        string  `json:"ctVal"`
			CtValCcy     string  `json:"ctValCcy"`
			ExpTime      string  `json:"expTime"`
			InstId       string  `json:"instId"`
			InstType     string  `json:"instType"`
			Lever        string  `json:"lever"`
			ListTime     float64 `json:"listTime,string"`
			LotSz        float64 `json:"lotSz,string"`
			MaxIcebergSz float64 `json:"maxIcebergSz,string"`
			MaxLmtSz     float64 `json:"maxLmtSz,string"`
			MaxMktSz     float64 `json:"maxMktSz,string"`
			MaxStopSz    float64 `json:"maxStopSz,string"`
			MaxTriggerSz float64 `json:"maxTriggerSz,string"`
			MaxTwapSz    float64 `json:"maxTwapSz,string"`
			MinSz        float64 `json:"minSz,string"`
			OptType      string  `json:"optType"`
			QuoteCcy     string  `json:"quoteCcy"`
			SettleCcy    string  `json:"settleCcy"`
			State        string  `json:"state"`
			Stk          string  `json:"stk"`
			TickSz       float64 `json:"tickSz,string"`
			Uly          string  `json:"uly"`
		} `json:"data"`
	}
	var response instrumentsResponse
	err = HttpGet4(ok.config.HttpClient, urlPath, nil, &response)
	if err != nil {
		return
	}
	if response.Code != 0 {
		err = fmt.Errorf("GetTickerV5 error:%s", response.Msg)
		return
	}
	for _, v := range response.Data {
		if v.State == "live" {
			sm[v.InstId] = Sym2duo(v.InstId)
			ps[Sym2duo(v.InstId)] = q.P{
				Base:      v.LotSz,
				BaseLimit: v.LotSz,
				Quote:     v.TickSz,
				Price:     v.TickSz,
				// MinQuote:     e["baseCurrency
				MinBaseLimit: v.MinSz,
				MinBase:      v.MinSz,
				TakerNett:    0.998}
		}
	}
	return
}
func Sym2duo(pair string) (res q.D) {
	parts := strings.Split(pair, "-")
	if len(parts) != 2 {
		panic("FATAL: SPLIT ERR! " + pair)
	}
	if parts[0] == parts[1] {
		panic("FATAL: SAME CURRENCY ERR! " + pair)
	}
	return q.D{Base: unify(parts[0]), Quote: unify(parts[1])}
}
func (ok *OKX) PlaceOrders(places [3]q.Order) (orders [3]q.Order, err error) {
	orders = places
	for _, p := range places {
		go func(p q.Order) {
		}(p)
	}
	return
}
func (ok *OKX) Test() bool {
	return true
}
func unify(local string) string {
	global := strings.ToUpper(local)
	return global
}
func (ok *OKX) GetAttr() (a q.Attr) {
	a.MultiOrder = false
	return a
}

func (ok *OKX) TradeFee() (map[string]q.TradeFee, error) {
	return nil, nil
}
func (ok *OKX) OneTicker(d q.D) (ticker q.Bbo, err error) {
	return
}
func (ok *OKX) AllTicker(SymPair map[string]q.D) (mdt *sync.Map, err error) {
	mdt = &sync.Map{}
	urlPath := ok.config.Endpoint + "/api/v5/market/tickers?instType=SPOT"
	var response TickerV5Response
	err = HttpGet4(ok.config.HttpClient, urlPath, nil, &response, `(\w+Px":)""`)
	if err != nil {
		return
	}
	if response.Code != 0 {
		return mdt, fmt.Errorf("GetTickerV5 error:%s", response.Msg)
	}
	for _, v := range response.Data {
		mdt.Store(
			Sym2duo(v.InstId),
			q.Bbo{
				Bid:     v.BuyPrice,
				BidSize: v.BuySize,
				Ask:     v.SellPrice,
				AskSize: v.SellSize,
				Pair:    v.InstId,
			})
	}
	return
}
