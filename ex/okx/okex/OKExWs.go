package okex

import (
	"encoding/json"
	"fmt"
	"os"
	. "github.com/conbanwa/wstrader/cons"
	. "github.com/conbanwa/wstrader/util"
	. "github.com/conbanwa/wstrader/web"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/conbanwa/logs"
)

type wsResp struct {
	Event     string `json:"event"`
	Channel   string `json:"channel"`
	Table     string `json:"table"`
	Data      json.RawMessage
	Success   bool `json:"success"`
	ErrorCode any  `json:"errorCode"`
}
type Ws struct {
	base *OKEx
	*WsBuilder
	once       *sync.Once
	WsConn     *WsConn
	respHandle func(channel string, data json.RawMessage) error
}

func NewOKExV3Ws(base *OKEx, handle func(channel string, data json.RawMessage) error) *Ws {
	okV3Ws := &Ws{
		once:       new(sync.Once),
		base:       base,
		respHandle: handle,
	}
	okV3Ws.WsBuilder = NewWsBuilder().
		ProxyUrl(os.Getenv("HTTPS_PROXY")).
		WsUrl("wss://ws.okx.com:8443/ws/v3").
		ReconnectInterval(time.Second).
		AutoReconnect().
		Heartbeat(func() []byte { return []byte("ping") }, 28*time.Second).
		DecompressFunc(FlateDecompress).ProtoHandleFunc(okV3Ws.handle)
	return okV3Ws
}
func (okV3Ws *Ws) clearChan(c chan wsResp) {
	for {
		if len(c) > 0 {
			<-c
		} else {
			break
		}
	}
}
func (okV3Ws *Ws) getTablePrefix(currencyPair CurrencyPair, contractType string) string {
	if contractType == SWAP_CONTRACT {
		return "swap"
	}
	return "futures"
}
func (okV3Ws *Ws) ConnectWs() {
	okV3Ws.once.Do(func() {
		okV3Ws.WsConn = okV3Ws.WsBuilder.Build()
	})
}
func (okV3Ws *Ws) parseChannel(channel string) (string, error) {
	metas := strings.Split(channel, "/")
	if len(metas) != 2 {
		return "", fmt.Errorf("unknown channel: %s", channel)
	}
	return metas[1], nil
}
func (okV3Ws *Ws) getKlinePeriodFormChannel(channel string) int {
	metas := strings.Split(channel, ":")
	if len(metas) != 2 {
		return 0
	}
	i, _ := strconv.ParseInt(metas[1], 10, 64)
	return int(i)
}
func (okV3Ws *Ws) handle(msg []byte) error {
	logs.D("[ws] [response] ", string(msg))
	if string(msg) == "pong" {
		return nil
	}
	var wsResp wsResp
	err := json.Unmarshal(msg, &wsResp)
	if err != nil {
		logs.E(err)
		return err
	}
	if wsResp.ErrorCode != nil {
		logs.E(string(msg))
		return fmt.Errorf("%s", string(msg))
	}
	if wsResp.Event != "" {
		switch wsResp.Event {
		case "subscribe":
			logs.I("subscribed:", wsResp.Channel)
			return nil
		case "error":
			log.Error().Msgf(string(msg))
		default:
			logs.I(string(msg))
		}
		return fmt.Errorf("unknown websocket message: %v", wsResp)
	}
	if wsResp.Table != "" {
		err = okV3Ws.respHandle(wsResp.Table, wsResp.Data)
		if err != nil {
			logs.E("handle ws data error:", err)
		}
		return err
	}
	return fmt.Errorf("unknown websocket message: %v", wsResp)
}
func (okV3Ws *Ws) Subscribe(sub map[string]any) error {
	okV3Ws.ConnectWs()
	return okV3Ws.WsConn.Subscribe(sub)
}
