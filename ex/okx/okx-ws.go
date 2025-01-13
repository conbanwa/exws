package okx

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
	"sort"
	"strings"
	"sync"
	"time"
)

const MaxSymbolChannels = 2
const MaxChannelSymbols = 358

type req struct {
	Op   string `json:"op"`
	Args []Arg  `json:"args"`
}
type Arg struct {
	Channel    string `json:"channel"`
	InstId     string `json:"instId,omitempty"`
	InstType   string `json:"instType,omitempty"`
	InstFamily string `json:"instFamily,omitempty"`
}

func toReq(pair ...string) req {
	args := make([]Arg, len(pair))
	for i, v := range pair {
		args[i] = Arg{
			Channel: "bbo-tbt",
			InstId:  v,
		}
	}
	return req{
		Op:   "subscribe",
		Args: args,
	}
}

type resp struct {
	Arg    Arg             `json:"arg"`
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
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
	spotWs.wsBuilder = web.NewWsBuilder().
		WsUrl(wsPublicUrl + "").
		ProxyUrl(config.GetProxy()).
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
	return s.c.Subscribe(toReq("BTC-USDT"))
}
func (s *SpotWs) SubscribeTicker(pair cons.CurrencyPair) error {
	defer func() {
		s.reqId++
	}()
	s.connect()
	return s.c.Subscribe(toReq("BTC-USDT"))
}
func (s *SpotWs) SubscribeBBO(sm []string) (err error) {
	if len(sm) <= 0 {
		return fmt.Errorf("nothing to subscribe")
	}
	s.connect()
	// n, pn := 0, MaxChannelSymbols
	// params := make([][]string, len(sm)/pn+1)
	// for _, k := range sm {
	// 	n++
	// 	params[n/pn] = append(params[n/pn], strings.ToLower(k)+"@bookTicker")
	// }
	// lp := len(params)
	// if lp > MaxSymbolChannels {
	// 	log.Error().Int("max", MaxSymbolChannels).Int("got", lp).Msg("too many symbol channels to subscribe")
	// 	lp = MaxSymbolChannels
	// }
	// for i := 0; i < lp; i++ {
	s.reqId++
	err = s.c.Subscribe(toReq(sm...))
	if err != nil {
		return
	}
	// }
	return
}

// TODO: test
func (s *SpotWs) SubscribeTrade(pair cons.CurrencyPair) error {
	defer func() {
		s.reqId++
	}()
	s.connect()
	return s.c.Subscribe(toReq("BTC-USDT"))
}
func (s *SpotWs) handle(data []byte) error {
	var r resp
	err := json.Unmarshal(data, &r)
	if err != nil {
		log.Error().Err(err).Bytes("response data", data).Msg("json unmarshal ws response error")
		return err
	}
	if len(r.Data) == 0 {
		log.Warn().Err(err).Bytes("response data", data).Msg("len(r.Data) == 0")
		return nil
	}
	if r.Arg.Channel == "bbo-tbt" {
		return s.bboHandle(r.Data, r.Arg.InstId)
	}
	if strings.HasPrefix(r.Arg.Channel, "books") {
		return s.bboHandle(r.Data, r.Arg.InstId)
	}
	// if strings.HasPrefix(r.Stream, "books") {
	// 	return s.depthHandle(r.Data, adaptStreamToCurrencyPair(r.Stream))
	// }
	// if strings.HasPrefix(r.Stream, "ticker") {
	// 	return s.tickerHandle(r.Data, adaptStreamToCurrencyPair(r.Stream))
	// }
	// if strings.HasPrefix(r.Stream, "trade") {
	// 	return s.tradeHandle(r.Data, adaptStreamToCurrencyPair(r.Stream))
	// }
	log.Warn().Bytes("handle", data).Msg("unknown ws response:")
	return nil
}

type bboResp []struct {
	Asks      [][]string `json:"asks"`
	Bids      [][]string `json:"bids"`
	Ts        string     `json:"ts"`
	Checksum  int        `json:"checksum"`
	PrevSeqID int        `json:"prevSeqId"`
	SeqID     int        `json:"seqId"`
}

func (s *SpotWs) bboHandle(data json.RawMessage, InstId string) error {
	if strings.Contains(slice.Bytes2String(data), "0.00000000") {
		return fmt.Errorf(cons.BINANCE + "0 in ask bid" + slice.Bytes2String(data))
	}
	var (
		tickerData bboResp
		ticker     q.Bbo
	)
	err := json.Unmarshal(data, &tickerData)
	if err != nil {
		log.Error().Err(err).Int("len", len(data)).Bytes("response data", data).Str("InstId", InstId).Msg("unmarshal bbo error")
		return err
	}
	ticker.Pair = InstId                                     // symbol
	ticker.Bid = num.ToFloat64(tickerData[0].Bids[0][0])     // best bid price
	ticker.BidSize = num.ToFloat64(tickerData[0].Bids[0][1]) // best bid qty
	ticker.Ask = num.ToFloat64(tickerData[0].Asks[0][0])     // best ask price
	ticker.AskSize = num.ToFloat64(tickerData[0].Asks[0][1]) // best ask qty
	// ticker.Updated = time.Now().UnixMilli()         // order book updateId
	s.bboCallFn(&ticker)
	return nil
}

type depthResp struct {
	LastUpdateId int     `json:"lastUpdateId"`
	Bids         [][]any `json:"bids"`
	Asks         [][]any `json:"asks"`
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
