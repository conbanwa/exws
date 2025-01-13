package okx

import (
	// "github.com/amir-the-h/okex"
	// "github.com/amir-the-h/okex/api"
	// "github.com/amir-the-h/okex/events"
	// "github.com/amir-the-h/okex/events/public"
	// ws_public_requests "github.com/amir-the-h/okex/requests/ws/public"
	// "os"
	// "context"
	"encoding/json"
	"fmt"
	"github.com/conbanwa/num"
	"github.com/conbanwa/slice"
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/config"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
	"github.com/conbanwa/wstrader/web"
	"sort"
	"strings"
	"sync"
	"time"
)

const MaxSymbolChannels = 2
const MaxChannelSymbols = 358

type req struct {
	Op string `json:"op"`
	Args   []Arg  `json:"args"`
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
			Channel:  "books",
			InstId: v,
		}
	}
	return req{
		Op:   "subscribe",
		Args: args,
	}
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
	depthCallFn  func(depth *wstrader.Depth)
	tickerCallFn func(ticker *wstrader.Ticker)
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

// func main() {
// 	apiKey := "YOUR-API-KEY"
// 	secretKey := "YOUR-SECRET-KEY"
// 	passphrase := "YOUR-PASS-PHRASE"
// 	dest := okex.NormalServer // The main API server
// 	ctx := context.Background()
// 	client, err := api.NewClient(ctx, apiKey, secretKey, passphrase, &dest)
// 	if err != nil {
// 		panic(err)
// 	}

// 	log.Println("Starting")
// 	errChan := make(chan *events.Error)
// 	subChan := make(chan *events.Subscribe)
// 	uSubChan := make(chan *events.Unsubscribe)
// 	logChan := make(chan *events.Login)
// 	sucChan := make(chan *events.Success)
// 	client.Ws.SetChannels(errChan, subChan, uSubChan, logChan, sucChan)

// 	obCh := make(chan *public.OrderBook)
// 	err = client.Ws.Public.OrderBook(ws_public_requests.OrderBook{
// 		InstID:  "BTC-USD-SWAP",
// 		Channel: "books",
// 	}, obCh)
// 	if err != nil {
// 		panic(err)
// 	}

//		for {
//			select {
//			case <-logChan:
//				log.Print("[Authorized]")
//			case success := <-sucChan:
//				log.Printf("[SUCCESS]\t%+v", success)
//			case sub := <-subChan:
//				channel, _ := sub.Arg.Get("channel")
//				log.Printf("[Subscribed]\t%s", channel)
//			case uSub := <-uSubChan:
//				channel, _ := uSub.Arg.Get("channel")
//				log.Printf("[Unsubscribed]\t%s", channel)
//			case err := <-client.Ws.ErrChan:
//				log.Printf("[Error]\t%+v", err)
//				for _, datum := range err.Data {
//					log.Printf("[Error]\t\t%+v", datum)
//				}
//			case i := <-obCh:
//				ch, _ := i.Arg.Get("channel")
//				log.Printf("[Event]\t%s", ch)
//				for _, p := range i.Books {
//					for i := len(p.Asks) - 1; i >= 0; i-- {
//						log.Printf("\t\tAsk\t%+v", p.Asks[i])
//					}
//					for _, bid := range p.Bids {
//						log.Printf("\t\tBid\t%+v", bid)
//					}
//				}
//			case b := <-client.Ws.DoneChan:
//				log.Printf("[End]:\t%v", b)
//				return
//			}
//		}
//	}
func (s *SpotWs) DepthCallback(f func(depth *wstrader.Depth)) {
	s.depthCallFn = f
}
func (s *SpotWs) TickerCallback(f func(ticker *wstrader.Ticker)) {
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
	// if strings.HasSuffix(r.Stream, "@bookTicker") {
	// 	return s.bboHandle(r.Data, adaptStreamToCurrencyPair(r.Stream))
	// }
	// if strings.HasSuffix(r.Stream, "@depth10@100ms") {
	// 	return s.depthHandle(r.Data, adaptStreamToCurrencyPair(r.Stream))
	// }
	// if strings.HasSuffix(r.Stream, "@ticker") {
	// 	return s.tickerHandle(r.Data, adaptStreamToCurrencyPair(r.Stream))
	// }
	// if strings.HasSuffix(r.Stream, "@aggTrade") {
	// 	return s.tradeHandle(r.Data, adaptStreamToCurrencyPair(r.Stream))
	// }
	log.Warn().Bytes("handle", data).Msg("unknown ws response:")
	return nil
}
func (s *SpotWs) depthHandle(data json.RawMessage, pair cons.CurrencyPair) error {
	var (
		depthR depthResp
		dep    wstrader.Depth
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
		dep.BidList = append(dep.BidList, wstrader.DepthRecord{
			Price:  num.ToFloat64(bid[0]),
			Amount: num.ToFloat64(bid[1]),
		})
	}
	for _, ask := range depthR.Asks {
		dep.AskList = append(dep.AskList, wstrader.DepthRecord{
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
		ticker     wstrader.Ticker
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
