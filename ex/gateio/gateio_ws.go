package gateio

import (
	"encoding/json"
	"fmt"
	"github.com/conbanwa/num"
	. "github.com/conbanwa/wstrader"
	. "github.com/conbanwa/wstrader/cons"
	. "github.com/conbanwa/wstrader/q"
	. "github.com/conbanwa/wstrader/web"
	"os"
	"sync"
	"time"

	"github.com/conbanwa/logs"
)

const subscribe = "subscribe"
const subscribed = "subscribed"
const trades = "trades"
const candles = "candles"

type Ws struct {
	*WsBuilder
	sync.Once
	wsConn         *WsConn
	eventMap       map[int64]SubscribeEvent
	tickerCallback func(*Ticker)
	bboCallback    func(*Bbo)
	depthCallback  func(*Depth)
	tradeCallback  func(*Trade)
	candleCallback func(*Kline)
}
type SubscribeEvent struct {
	Event     string `json:"event"`
	SubID     string `json:"subId"`
	Channel   string `json:"channel"`
	ChanID    int64  `json:"chanId"`
	Symbol    string `json:"symbol"`
	Precision string `json:"prec,omitempty"`
	Frequency string `json:"freq,omitempty"`
	Key       string `json:"key,omitempty"`
	Len       string `json:"len,omitempty"`
	Pair      string `json:"pair"`
}
type EventMap map[int64]SubscribeEvent

func NewWs() *Ws {
	ws := &Ws{WsBuilder: NewWsBuilder(), eventMap: make(map[int64]SubscribeEvent)}
	ws.WsBuilder = ws.WsBuilder.ProxyUrl(os.Getenv("HTTPS_PROXY")).
		WsUrl("wss://ws.gateio.io/v3/").
		AutoReconnect().
		// DisableEnableCompression().
		ProtoHandleFunc(ws.handle)
	return ws
}
func (ws *Ws) SubscribeBBO(sm []string) (err error) {
	if ws.bboCallback == nil {
		return fmt.Errorf("please set bbo callback func")
	}
	for sym := range sm {
		if err = ws.subscribe(map[string]any{
			"event":   subscribe,
			"channel": "ticker",
			"symbol":  sym,
		}); err != nil {
			return err
		}
		time.Sleep(60 * time.Second)
	}
	return
}
func (ws *Ws) SubscribeTicker(pair CurrencyPair) error {
	panic("not implement")
}
func (ws *Ws) SubscribeDepth(pair CurrencyPair) error {
	panic("not implement")
}
func (ws *Ws) SubscribeTrade(pair CurrencyPair) error {
	panic("not implement")
}
func (ws *Ws) subscribe(sub map[string]any) error {
	ws.connectWs()
	return ws.wsConn.Subscribe(sub)
}
func (ws *Ws) connectWs() {
	ws.Do(func() {
		ws.wsConn = ws.WsBuilder.Build()
	})
}
func (ws *Ws) handle(msg []byte) error {
	var event SubscribeEvent
	if err := json.Unmarshal(msg, &event); err == nil {
		switch event.Event {
		case subscribed:
			ws.eventMap[event.ChanID] = event
			logs.I(event)
			return nil
		case "unsubscribed":
			logs.I(event)
		case "error":
			logs.E(string(msg))
		default:
			logs.E(event)
		}
	}
	var resp []any
	if err := json.Unmarshal(msg, &resp); err == nil {
		channelID := num.ToInt[int64](resp[0])
		event, ok := ws.eventMap[channelID]
		if !ok {
			return nil
		}
		switch event.Channel {
		case "ticker":
			if raw, ok := resp[1].([]any); ok {
				t := ws.bboFromRaw(event.Symbol, raw)
				ws.bboCallback(t)
				return nil
			}
		default:
			logs.E(event)
		}
	}
	return nil
}
func (ws *Ws) tickerFromRaw(pair CurrencyPair, raw []any) *Ticker {
	return &Ticker{
		Pair: pair,
		Buy:  num.ToFloat64(raw[0]),
		Sell: num.ToFloat64(raw[2]),
		Last: num.ToFloat64(raw[6]),
		Vol:  num.ToFloat64(raw[7]),
		High: num.ToFloat64(raw[8]),
		Low:  num.ToFloat64(raw[9]),
		Date: uint64(time.Now().UnixNano() / int64(time.Millisecond)),
	}
}
func (ws *Ws) bboFromRaw(pair string, raw []any) *Bbo {
	return &Bbo{
		Pair:    pair,
		Bid:     num.ToFloat64(raw[0]),
		BidSize: num.ToFloat64(raw[1]),
		Ask:     num.ToFloat64(raw[2]),
		AskSize: num.ToFloat64(raw[3]),
	}
}
func (ws *Ws) TickerCallback(tickerCallback func(*Ticker)) {
	ws.tickerCallback = tickerCallback
}
func (ws *Ws) BBOCallback(bboCallback func(*Bbo)) {
	ws.bboCallback = bboCallback
}
func (ws *Ws) DepthCallback(depthCallback func(*Depth)) {
	ws.depthCallback = depthCallback
}
func (ws *Ws) TradeCallback(tradeCallback func(*Trade)) {
	ws.tradeCallback = tradeCallback
}
