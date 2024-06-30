package binance

import (
	"encoding/json"
	"github.com/conbanwa/num"
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
	"github.com/conbanwa/wstrader/web"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type FuturesWs struct {
	base         *Futures
	fOnce        sync.Once
	dOnce        sync.Once
	wsBuilder    *web.WsBuilder
	f            *web.WsConn
	d            *web.WsConn
	depthCallFn  func(depth *wstrader.Depth)
	tickerCallFn func(ticker *wstrader.FutureTicker)
	tradeCalFn   func(trade *q.Trade, contract string)
}

func NewFuturesWs() *FuturesWs {
	futuresWs := new(FuturesWs)
	futuresWs.wsBuilder = web.NewWsBuilder().
		ProxyUrl(os.Getenv("HTTPS_PROXY")).
		ProtoHandleFunc(futuresWs.handle).AutoReconnect()
	httpCli := &http.Client{
		Timeout: 10 * time.Second,
	}
	if os.Getenv("HTTPS_PROXY") != "" {
		httpCli = &http.Client{
			Transport: &http.Transport{
				Proxy: func(r *http.Request) (*url.URL, error) {
					return url.Parse(os.Getenv("HTTPS_PROXY"))
				},
			},
			Timeout: 10 * time.Second,
		}
	}
	futuresWs.base = NewBinanceFutures(&wstrader.APIConfig{
		HttpClient: httpCli,
	})
	return futuresWs
}
func (s *FuturesWs) connectUsdtFutures() {
	s.fOnce.Do(func() {
		s.f = s.wsBuilder.WsUrl(TestnetFutureUsdWsBaseUrl).Build()
	})
}
func (s *FuturesWs) connectFutures() {
	s.dOnce.Do(func() {
		s.d = s.wsBuilder.WsUrl(TestnetFutureCoinWsBaseUrl).Build()
	})
}
func (s *FuturesWs) DepthCallback(f func(depth *wstrader.Depth)) {
	s.depthCallFn = f
}
func (s *FuturesWs) TickerCallback(f func(ticker *wstrader.FutureTicker)) {
	s.tickerCallFn = f
}
func (s *FuturesWs) TradeCallback(f func(trade *q.Trade, contract string)) {
	s.tradeCalFn = f
}
func (s *FuturesWs) SubscribeDepth(pair cons.CurrencyPair, contractType string) error {
	switch contractType {
	case cons.SWAP_USDT_CONTRACT:
		s.connectUsdtFutures()
		return s.f.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: []string{pair.AdaptUsdToUsdt().ToLower().ToSymbol("") + "@depth10@100ms"},
			Id:     1,
		})
	default:
		s.connectFutures()
		sym, _ := s.base.adaptToSymbol(pair.AdaptUsdtToUsd(), contractType)
		return s.d.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: []string{strings.ToLower(sym) + "@depth10@100ms"},
			Id:     2,
		})
	}
}
func (s *FuturesWs) SubscribeTicker(pair cons.CurrencyPair, contractType string) error {
	switch contractType {
	case cons.SWAP_USDT_CONTRACT:
		s.connectUsdtFutures()
		return s.f.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: []string{pair.AdaptUsdToUsdt().ToLower().ToSymbol("") + "@ticker"},
			Id:     1,
		})
	default:
		s.connectFutures()
		sym, _ := s.base.adaptToSymbol(pair.AdaptUsdtToUsd(), contractType)
		return s.d.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: []string{strings.ToLower(sym) + "@ticker"},
			Id:     2,
		})
	}
}
func (s *FuturesWs) SubscribeTrade(pair cons.CurrencyPair, contractType string) error {
	switch contractType {
	case cons.SWAP_USDT_CONTRACT:
		s.connectUsdtFutures()
		return s.f.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: []string{pair.AdaptUsdToUsdt().ToLower().ToSymbol("") + "@aggTrade"},
			Id:     1,
		})
	default:
		s.connectFutures()
		sym, _ := s.base.adaptToSymbol(pair.AdaptUsdtToUsd(), contractType)
		return s.d.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: []string{strings.ToLower(sym) + "@aggTrade"},
			Id:     1,
		})
	}
}
func (s *FuturesWs) handle(data []byte) error {
	var tickers = make(map[string]any, 4)
	err := json.Unmarshal(data, &tickers)
	if err != nil {
		return err
	}
	if e, ok := tickers["e"].(string); ok && e == "depthUpdate" {
		dep := s.depthHandle(tickers["b"].([]any), tickers["a"].([]any))
		dep.ContractType = tickers["s"].(string)
		symbol, ok := tickers["ps"].(string)
		if ok {
			dep.Pair = adaptSymbolToCurrencyPair(symbol)
		} else {
			dep.Pair = adaptSymbolToCurrencyPair(dep.ContractType) //usdt swap
		}
		dep.UTime = time.Unix(0, num.ToInt[int64](tickers["T"])*int64(time.Millisecond))
		s.depthCallFn(dep)
		return nil
	}
	if e, ok := tickers["e"].(string); ok && e == "24hrTicker" {
		s.tickerCallFn(s.tickerHandle(tickers))
		return nil
	}
	if e, ok := tickers["e"].(string); ok && e == "aggTrade" {
		contractType := tickers["s"].(string)
		s.tradeCalFn(s.tradeHandle(tickers), contractType)
		return nil
	}
	log.Info().Bytes("handle", data).Msg("unknown ws response:")
	return nil
}
func (s *FuturesWs) depthHandle(bids []any, asks []any) *wstrader.Depth {
	var dep wstrader.Depth
	for _, item := range bids {
		bid := item.([]any)
		dep.BidList = append(dep.BidList,
			wstrader.DepthRecord{
				Price:  num.ToFloat64(bid[0]),
				Amount: num.ToFloat64(bid[1]),
			})
	}
	for _, item := range asks {
		ask := item.([]any)
		dep.AskList = append(dep.AskList, wstrader.DepthRecord{
			Price:  num.ToFloat64(ask[0]),
			Amount: num.ToFloat64(ask[1]),
		})
	}
	sort.Sort(sort.Reverse(dep.AskList))
	return &dep
}
func (s *FuturesWs) tickerHandle(tickers map[string]any) *wstrader.FutureTicker {
	var ticker wstrader.FutureTicker
	ticker.Ticker = new(wstrader.Ticker)
	symbol, ok := tickers["ps"].(string)
	if ok {
		ticker.Pair = adaptSymbolToCurrencyPair(symbol)
	} else {
		ticker.Pair = adaptSymbolToCurrencyPair(tickers["s"].(string)) //usdt swap
	}
	ticker.ContractType = tickers["s"].(string)
	ticker.Date = num.ToInt[uint64](tickers["E"])
	ticker.High = num.ToFloat64(tickers["h"])
	ticker.Low = num.ToFloat64(tickers["l"])
	ticker.Last = num.ToFloat64(tickers["c"])
	ticker.Vol = num.ToFloat64(tickers["v"])
	return &ticker
}
func (s *FuturesWs) tradeHandle(tickers map[string]any) *q.Trade {
	var trade q.Trade
	symbol, ok := tickers["s"].(string) // Symbol
	if ok {
		trade.Pair = adaptSymbolToCurrencyPair(symbol) //usdt swap
	}
	trade.Tid = num.ToInt[int64](tickers["a"])  // Aggregate trade ID
	trade.Date = num.ToInt[int64](tickers["E"]) // Event time
	trade.Amount = num.ToFloat64(tickers["q"])  // Quantity
	trade.Price = num.ToFloat64(tickers["p"])   // Price
	if tickers["m"].(bool) {
		trade.Type = cons.BUY_MARKET
	} else {
		trade.Type = cons.SELL_MARKET
	}
	return &trade
}
