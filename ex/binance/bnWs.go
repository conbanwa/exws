package binance

import (
	"encoding/json"
	"fmt"
	"github.com/conbanwa/exws"
	"github.com/conbanwa/exws/config"
	"github.com/conbanwa/exws/cons"
	"github.com/conbanwa/exws/q"
	"github.com/conbanwa/exws/web"
	"github.com/conbanwa/num"
	"github.com/conbanwa/slice"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

const MaxSymbolChannels = 2
const MaxChannelSymbols = 358

type req struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     int      `json:"id"`
}
type resp struct {
	Stream string          `json:"stream"`
	Data   json.RawMessage `json:"data"`
}
type depthResp struct {
	LastUpdateId int     `json:"lastUpdateId"`
	Bids         [][]any `json:"bids"`
	Asks         [][]any `json:"asks"`
}
type SpotWs struct {
	c            *web.WsConn
	once         sync.Once
	wsBuilder    *web.WsBuilder
	reqId        int
	depthCallFn  func(depth *exws.Depth)
	tickerCallFn func(ticker *exws.Ticker)
	bboCallFn    func(ticker *q.Bbo)
	tradeCallFn  func(trade *q.Trade)
}

func NewSpotWs() *SpotWs {
	spotWs := &SpotWs{}
	proxyUrl := ""
	if config.UseProxy {
		log.Printf("proxy url: %s", os.Getenv("HTTPS_PROXY"))
		proxyUrl = os.Getenv("HTTPS_PROXY")
	}
	spotWs.wsBuilder = web.NewWsBuilder().
		WsUrl(TestnetSpotStreamBaseUrl + "?streams=depth/miniTicker/ticker/trade").
		ProxyUrl(proxyUrl).
		ProtoHandleFunc(spotWs.handle).AutoReconnect()
	return spotWs
}
func (s *SpotWs) connect() {
	s.once.Do(func() {
		s.c = s.wsBuilder.Build()
	})
}
func (s *SpotWs) DepthCallback(f func(depth *exws.Depth)) {
	s.depthCallFn = f
}
func (s *SpotWs) TickerCallback(f func(ticker *exws.Ticker)) {
	s.tickerCallFn = f
}
func (s *SpotWs) BBOCallback(f func(ticker *q.Bbo)) {
	s.bboCallFn = f
}
func (s *SpotWs) TradeCallback(f func(trade *q.Trade)) {
	s.tradeCallFn = f
}
func (s *SpotWs) SubscribeDepth(pair cons.CurrencyPair) error {
	defer func() {
		s.reqId++
	}()
	s.connect()
	return s.c.Subscribe(req{
		Method: "SUBSCRIBE",
		Params: []string{
			fmt.Sprintf("%s@depth10@100ms", pair.ToLower().ToSymbol("")),
		},
		Id: s.reqId,
	})
}
func (s *SpotWs) SubscribeTicker(pair cons.CurrencyPair) error {
	defer func() {
		s.reqId++
	}()
	s.connect()
	return s.c.Subscribe(req{
		Method: "SUBSCRIBE",
		Params: []string{pair.ToLower().ToSymbol("") + "@ticker"},
		Id:     s.reqId,
	})
}
func (s *SpotWs) SubscribeBBO(sm []string) (err error) {
	if len(sm) <= 0 {
		return fmt.Errorf("nothing to subscribe")
	}
	s.connect()
	n, pn := 0, MaxChannelSymbols
	params := make([][]string, len(sm)/pn+1)
	for _, k := range sm {
		n++
		params[n/pn] = append(params[n/pn], strings.ToLower(k)+"@bookTicker")
	}
	lp := len(params)
	if lp > MaxSymbolChannels {
		log.Error().Int("max", MaxSymbolChannels).Int("got", lp).Msg("too many symbol channels to subscribe")
		lp = MaxSymbolChannels
	}
	for i := 0; i < lp; i++ {
		s.reqId++
		err = s.c.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: params[i],
			Id:     s.reqId,
		})
		if err != nil {
			return
		}
	}
	return
}

// TODO: test
func (s *SpotWs) SubscribeTrade(pair cons.CurrencyPair) error {
	defer func() {
		s.reqId++
	}()
	s.connect()
	return s.c.Subscribe(req{
		Method: "SUBSCRIBE",
		Params: []string{pair.ToLower().ToSymbol("") + "@aggTrade"},
		Id:     s.reqId,
	})
}
func (s *SpotWs) handle(data []byte) error {
	var r resp
	err := json.Unmarshal(data, &r)
	if err != nil {
		log.Error().Err(err).Bytes("response data", data).Msg("json unmarshal ws response error")
		return err
	}
	if strings.HasSuffix(r.Stream, "@bookTicker") {
		return s.bboHandle(r.Data, adaptStreamToCurrencyPair(r.Stream))
	}
	if strings.HasSuffix(r.Stream, "@depth10@100ms") {
		return s.depthHandle(r.Data, adaptStreamToCurrencyPair(r.Stream))
	}
	if strings.HasSuffix(r.Stream, "@ticker") {
		return s.tickerHandle(r.Data, adaptStreamToCurrencyPair(r.Stream))
	}
	if strings.HasSuffix(r.Stream, "@aggTrade") {
		return s.tradeHandle(r.Data, adaptStreamToCurrencyPair(r.Stream))
	}
	log.Warn().Bytes("handle", data).Msg("unknown ws response:")
	return nil
}
func (s *SpotWs) depthHandle(data json.RawMessage, pair cons.CurrencyPair) error {
	var (
		depthR depthResp
		dep    exws.Depth
		err    error
	)
	err = json.Unmarshal(data, &depthR)
	if err != nil {
		log.Error().Err(err).Bytes("response data", data).Msg("unmarshal depth response error")
		return err
	}
	dep.UTime = time.Now()
	dep.Pair = pair
	for _, bid := range depthR.Bids {
		dep.BidList = append(dep.BidList, exws.DepthRecord{
			Price:  num.ToFloat64(bid[0]),
			Amount: num.ToFloat64(bid[1]),
		})
	}
	for _, ask := range depthR.Asks {
		dep.AskList = append(dep.AskList, exws.DepthRecord{
			Price:  num.ToFloat64(ask[0]),
			Amount: num.ToFloat64(ask[1]),
		})
	}
	sort.Sort(sort.Reverse(dep.AskList))
	s.depthCallFn(&dep)
	return nil
}
func (s *SpotWs) tickerHandle(data json.RawMessage, pair cons.CurrencyPair) error {
	var (
		tickerData = make(map[string]any, 4)
		ticker     exws.Ticker
	)
	err := json.Unmarshal(data, &tickerData)
	if err != nil {
		log.Error().Err(err).Bytes("response data", data).Msg("unmarshal ticker response data error")
		return err
	}
	ticker.Pair = pair
	ticker.Vol = num.ToFloat64(tickerData["v"])
	ticker.Last = num.ToFloat64(tickerData["c"])
	ticker.Sell = num.ToFloat64(tickerData["a"])
	ticker.Buy = num.ToFloat64(tickerData["b"])
	ticker.High = num.ToFloat64(tickerData["h"])
	ticker.Low = num.ToFloat64(tickerData["l"])
	ticker.Date = num.ToInt[uint64](tickerData["E"])
	s.tickerCallFn(&ticker)
	return nil
}
func (s *SpotWs) bboHandle(data json.RawMessage, pair cons.CurrencyPair) error {
	if strings.Contains(slice.Bytes2String(data), "0.00000000") {
		return fmt.Errorf(cons.BINANCE + "0 in ask bid" + slice.Bytes2String(data))
	}
	var (
		tickerData = make(map[string]any, 4)
		ticker     q.Bbo
	)
	err := json.Unmarshal(data, &tickerData)
	if err != nil {
		log.Error().Err(err).Bytes("response data", data).Msg("unmarshal ticker response data error")
		return err
	}
	ticker.Pair = tickerData["s"].(string)          // symbol
	ticker.Bid = num.ToFloat64(tickerData["b"])     // best bid price
	ticker.BidSize = num.ToFloat64(tickerData["B"]) // best bid qty
	ticker.Ask = num.ToFloat64(tickerData["a"])     // best ask price
	ticker.AskSize = num.ToFloat64(tickerData["A"]) // best ask qty
	ticker.Updated = time.Now().UnixMilli()         // "u" order book updateId
	s.bboCallFn(&ticker)
	return nil
}
func (s *SpotWs) tradeHandle(data json.RawMessage, pair cons.CurrencyPair) error {
	var (
		tradeData = make(map[string]any, 4)
		trade     q.Trade
	)
	err := json.Unmarshal(data, &tradeData)
	if err != nil {
		log.Error().Err(err).Bytes("response data", data).Msg("unmarshal ticker response data error")
		return err
	}
	trade.Pair = pair                             //Symbol
	trade.Tid = num.ToInt[int64](tradeData["a"])  // Aggregate trade ID
	trade.Date = num.ToInt[int64](tradeData["E"]) // Event time
	trade.Amount = num.ToFloat64(tradeData["q"])  // Quantity
	trade.Price = num.ToFloat64(tradeData["p"])   // Price
	if tradeData["m"].(bool) {
		trade.Type = cons.BUY_MARKET
	} else {
		trade.Type = cons.SELL_MARKET
	}
	return nil
}
