package ftx

import (
	"encoding/json"
	"fmt"
	"os"
	. "github.com/conbanwa/wstrader"
	. "github.com/conbanwa/wstrader/cons"
	. "github.com/conbanwa/wstrader/q"
	. "github.com/conbanwa/wstrader/web"
	"sync"
	"time"

	"github.com/conbanwa/logs"
)

const subscribe = "subscribe"
const subscribed = "subscribed"
const ticker = "ticker"
const trades = "trades"

type Ws struct {
	*WsBuilder
	sync.Once
	wsConn         *WsConn
	eventMap       map[int64]SubscribeEvent
	tickerCallback func(*Ticker)
	bboCallback    func(*Bbo)
	depthCallback  func(*Depth)
	tradeCallback  func(*Trade)
}
type SubscribeEvent struct {
	Channel string `json:"channel"`
	Symbol  string `json:"market"`
	Type    string `json:"type"`
	Data    struct {
		Bid     float64 `json:"bid"`
		Ask     float64 `json:"ask"`
		BidSize float64 `json:"bidSize"`
		AskSize float64 `json:"askSize"`
		Last    float64 `json:"last"`
		Time    float64 `json:"time"`
	} `json:"data"`
}

//	type SubscribeEvent struct {
//		Channel string `json:"channel"`
//		Symbol  string `json:"market"`
//		Type    string `json:"type"`
//
// Event     string `json:"op"`
// SubID     string `json:"subId"`
// ChanID    int64  `json:"chanId"`
// Precision string `json:"prec,omitempty"`
// Frequency string `json:"freq,omitempty"`
// Key       string `json:"key,omitempty"`
// Len       string `json:"len,omitempty"`
// Pair      string `json:"pair"`
// }
type EventMap map[int64]SubscribeEvent

func NewWs() *Ws {
	fws := &Ws{WsBuilder: NewWsBuilder(), eventMap: make(map[int64]SubscribeEvent)}
	fws.WsBuilder = fws.WsBuilder.ProxyUrl(os.Getenv("HTTPS_PROXY")).
		WsUrl("wss://ftx.com/ws/").
		AutoReconnect().
		DisableEnableCompression().
		ProtoHandleFunc(fws.handle)
	return fws
}
func (fws *Ws) TickerCallback(tickerCallback func(*Ticker)) {
	fws.tickerCallback = tickerCallback
}
func (fws *Ws) BBOCallback(bboCallback func(*Bbo)) {
	fws.bboCallback = bboCallback
}
func (fws *Ws) DepthCallback(depthCallback func(*Depth)) {
	fws.depthCallback = depthCallback
}
func (fws *Ws) TradeCallback(tradeCallback func(*Trade)) {
	fws.tradeCallback = tradeCallback
}
func (fws *Ws) SubscribeTicker(pair CurrencyPair) error {
	if fws.tickerCallback == nil {
		return fmt.Errorf("please set ticker callback func")
	}
	return fws.subscribe(map[string]any{
		"op":      subscribe,
		"channel": ticker,
		// "market":  convertPairToBitfinexSymbol("t", pair)
	})
}
func (fws *Ws) SubscribeBBO(sm []string) (err error) {
	if fws.bboCallback == nil {
		return fmt.Errorf("please set bbo callback func")
	}
	for sym := range sm {
		if err = fws.subscribe(map[string]any{
			"op":      subscribe,
			"channel": ticker,
			"market":  sym,
		}); err != nil {
			return err
		}
		time.Sleep(60 * time.Millisecond)
	}
	return
}
func (fws *Ws) SubscribeDepth(pair CurrencyPair) error {
	return nil
}
func (fws *Ws) SubscribeTrade(pair CurrencyPair) error {
	if fws.tradeCallback == nil {
		return fmt.Errorf("please set trade callback func")
	}
	return fws.subscribe(map[string]any{
		"op":      subscribe,
		"channel": trades,
		// "market":  convertPairToBitfinexSymbol("t", pair)
	},
	)
}
func (fws *Ws) subscribe(sub map[string]any) error {
	fws.connectWs()
	return fws.wsConn.Subscribe(sub)
}
func (fws *Ws) connectWs() {
	fws.Do(func() {
		fws.wsConn = fws.WsBuilder.Build()
	})
}
func (fws *Ws) handle(msg []byte) error {
	var event SubscribeEvent
	if err := json.Unmarshal(msg, &event); err == nil {
		switch event.Type {
		case subscribed:
			// fws.eventMap[event.ChanID] = event
			// logs.I(event)
			return nil
		case "unsubscribed":
			logs.I(event)
		case "error":
			logs.E(string(msg))
		default:
			// logs.E(event)
		}
	}
	switch event.Channel {
	case ticker:
		fws.bboCallback(&Bbo{
			Pair:    event.Symbol,
			Bid:     event.Data.Bid,
			BidSize: event.Data.BidSize,
			Ask:     event.Data.Ask,
			AskSize: event.Data.AskSize,
		})
		return nil
	default:
		logs.E(event)
	}
	logs.E(string(msg))
	return nil
}
