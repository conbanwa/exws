package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/conbanwa/exws/config"
	"github.com/conbanwa/exws/stat/zelo"
	"github.com/gorilla/websocket"
)

type WsConfig struct {
	WsUrl                          string
	ProxyUrl                       string
	ReqHeaders                     map[string][]string //连接的时候加入的头部信息
	HeartbeatIntervalTime          time.Duration       //
	HeartbeatData                  func() []byte       //心跳数据2
	IsAutoReconnect                bool
	ProtoHandleFunc                func([]byte) error           //协议处理函数
	DecompressFunc                 func([]byte) ([]byte, error) //解压函数
	ErrorHandleFunc                func(err error)
	ConnectSuccessAfterSendMessage func() []byte //for reconnect
	IsDump                         bool
	DisableEnableCompression       bool
	readDeadLineTime               time.Duration
	reconnectInterval              time.Duration
}

var dialer = &websocket.Dialer{
	Proxy:             http.ProxyFromEnvironment,
	HandshakeTimeout:  30 * time.Second,
	EnableCompression: true,
}

type WsConn struct {
	c *websocket.Conn
	WsConfig
	writeBufferChan        chan []byte
	pingMessageBufferChan  chan []byte
	pongMessageBufferChan  chan []byte
	closeMessageBufferChan chan []byte
	subs                   atomic.Value
	close                  chan bool
	*sync.Mutex
}
type WsBuilder struct {
	wsConfig *WsConfig
}

func NewWsBuilder() *WsBuilder {
	return &WsBuilder{&WsConfig{
		ReqHeaders:        make(map[string][]string, 1),
		reconnectInterval: time.Second * 10,
		ProxyUrl:          config.GetProxy(true),
	}}
}
func (b *WsBuilder) WsUrl(wsUrl string) *WsBuilder {
	b.wsConfig.WsUrl = wsUrl
	return b
}
func (b *WsBuilder) ProxyUrl(proxyUrl string) *WsBuilder {
	b.wsConfig.ProxyUrl = proxyUrl
	return b
}
func (b *WsBuilder) ReqHeader(key, value string) *WsBuilder {
	b.wsConfig.ReqHeaders[key] = append(b.wsConfig.ReqHeaders[key], value)
	return b
}
func (b *WsBuilder) AutoReconnect() *WsBuilder {
	b.wsConfig.IsAutoReconnect = true
	return b
}
func (b *WsBuilder) Dump() *WsBuilder {
	b.wsConfig.IsDump = true
	return b
}
func (b *WsBuilder) Heartbeat(heartbeat func() []byte, t time.Duration) *WsBuilder {
	b.wsConfig.HeartbeatIntervalTime = t
	b.wsConfig.HeartbeatData = heartbeat
	return b
}
func (b *WsBuilder) ReconnectInterval(t time.Duration) *WsBuilder {
	b.wsConfig.reconnectInterval = t
	return b
}
func (b *WsBuilder) ProtoHandleFunc(f func([]byte) error) *WsBuilder {
	b.wsConfig.ProtoHandleFunc = f
	return b
}
func (b *WsBuilder) DisableEnableCompression() *WsBuilder {
	b.wsConfig.DisableEnableCompression = true
	return b
}
func (b *WsBuilder) DecompressFunc(f func([]byte) ([]byte, error)) *WsBuilder {
	b.wsConfig.DecompressFunc = f
	return b
}
func (b *WsBuilder) ErrorHandleFunc(f func(err error)) *WsBuilder {
	b.wsConfig.ErrorHandleFunc = f
	return b
}
func (b *WsBuilder) ConnectSuccessAfterSendMessage(msg func() []byte) *WsBuilder {
	b.wsConfig.ConnectSuccessAfterSendMessage = msg
	return b
}
func (b *WsBuilder) Build() *WsConn {
	wsConn := &WsConn{WsConfig: *b.wsConfig}
	return wsConn.NewWs()
}
func (ws *WsConn) NewWs() *WsConn {
	if ws.HeartbeatIntervalTime == 0 {
		ws.readDeadLineTime = time.Minute
	} else {
		ws.readDeadLineTime = ws.HeartbeatIntervalTime * 2
	}
	if err := ws.connect(); err != nil {
		panic(fmt.Errorf("[%s] %s", ws.WsUrl, err.Error()))
	}
	ws.close = make(chan bool, 1)
	ws.pingMessageBufferChan = make(chan []byte, 10)
	ws.pongMessageBufferChan = make(chan []byte, 10)
	ws.closeMessageBufferChan = make(chan []byte, 10)
	ws.writeBufferChan = make(chan []byte, 10)
	ws.Mutex = new(sync.Mutex)
	go ws.writeRequest()
	go ws.receiveMessage()
	if ws.ConnectSuccessAfterSendMessage != nil {
		msg := ws.ConnectSuccessAfterSendMessage()
		ws.SendMessage(msg)
		log.Error().Bytes("msg data", msg).Str("url", ws.WsUrl).Msg("[ws] execute the connect success after send")
	}
	return ws
}
func (ws *WsConn) connect() error {
	if ws.ProxyUrl != "" {
		proxy, err := url.Parse(ws.ProxyUrl)
		if err == nil {
			log.Info().Msgf("[ws][%s] proxy url:%s", ws.WsUrl, proxy)
			dialer.Proxy = http.ProxyURL(proxy)
		} else {
			log.Error().Msgf("[ws][%s]parse proxy url [%s] err %s  ", ws.WsUrl, ws.ProxyUrl, err.Error())
		}
	}
	if ws.DisableEnableCompression {
		dialer.EnableCompression = false
	}
	wsConn, resp, err := dialer.Dial(ws.WsUrl, ws.ReqHeaders)
	if err != nil {
		log.Error().Msgf("[ws][%s] %s", ws.WsUrl, err.Error())
		if ws.IsDump && resp != nil {
			dumpData, _ := httputil.DumpResponse(resp, true)
			log.Debug().Bytes("dumpData", dumpData).Str("url", ws.WsUrl).Msg("DumpResponse")
		}
		return err
	}
	wsConn.SetReadDeadline(time.Now().Add(ws.readDeadLineTime))
	if ws.IsDump {
		dumpData, _ := httputil.DumpResponse(resp, true)
		log.Debug().Bytes("dumpData", dumpData).Str("url", ws.WsUrl).Msg("DumpResponse")
	}
	log.Info().Msgf("[ws][%s] connected", ws.WsUrl)
	ws.c = wsConn
	return nil
}
func (ws *WsConn) reconnect() {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	ws.c.Close() //主动关闭一次
	var err error
	for retry := 1; retry <= 100; retry++ {
		time.Sleep(ws.WsConfig.reconnectInterval * time.Duration(retry))
		err = ws.connect()
		if err != nil {
			log.Error().Msgf("[ws] [%s] websocket reconnect fail , %s", ws.WsUrl, err.Error())
		} else {
			break
		}
	}
	if err != nil {
		log.Error().Msgf("[ws] [%s] retry connect 100 count fail , begin exiting. ", ws.WsUrl)
		ws.CloseWs()
		if ws.ErrorHandleFunc != nil {
			ws.ErrorHandleFunc(errors.New("retry reconnect fail"))
		}
	} else {
		//re subscribe
		if ws.ConnectSuccessAfterSendMessage != nil {
			msg := ws.ConnectSuccessAfterSendMessage()
			ws.SendMessage(msg)
			log.Error().Bytes("msg data", msg).Str("url", ws.WsUrl).Msg("[ws] execute the connect success after send")
			time.Sleep(time.Second) //wait response
		}

		if subbed := ws.subs.Load(); subbed != nil {
			for _, sub := range subbed.([][]byte) {
				log.Debug().Bytes("sub", sub).Str("url", ws.WsUrl).Msg("[ws] re subscribe")
				ws.SendMessage(sub)
			}
		}
	}
}
func (ws *WsConn) writeRequest() {
	var (
		heartTimer *time.Timer
		err        error
	)
	if ws.HeartbeatIntervalTime == 0 {
		heartTimer = time.NewTimer(time.Hour)
	} else {
		heartTimer = time.NewTimer(ws.HeartbeatIntervalTime)
	}
	for {
		select {
		case <-ws.close:
			log.Info().Msgf("[ws][%s] close websocket , exiting write message goroutine.", ws.WsUrl)
			return
		case d := <-ws.writeBufferChan:
			err = ws.c.WriteMessage(websocket.TextMessage, d)
		case d := <-ws.pingMessageBufferChan:
			err = ws.c.WriteMessage(websocket.PingMessage, d)
		case d := <-ws.pongMessageBufferChan:
			err = ws.c.WriteMessage(websocket.PongMessage, d)
		case d := <-ws.closeMessageBufferChan:
			err = ws.c.WriteMessage(websocket.CloseMessage, d)
		case <-heartTimer.C:
			if ws.HeartbeatIntervalTime > 0 {
				err = ws.c.WriteMessage(websocket.TextMessage, ws.HeartbeatData())
				heartTimer.Reset(ws.HeartbeatIntervalTime)
			}
		}
		if err != nil {
			log.Error().Msgf("[ws][%s] write message %s", ws.WsUrl, err.Error())
			//time.Sleep(time.Second)
		}
	}
}
func (ws *WsConn) Subscribe(subEvent any) error {
	data, err := json.Marshal(subEvent)
	if err != nil {
		log.Error().Msgf("[ws][%s] json encode error , %s", ws.WsUrl, err)
		return err
	}
	//logs.D("subscribed url:", slice.Bytes2String(data))
	ws.writeBufferChan <- data
	if subbed := ws.subs.Load(); subbed != nil {
		ws.subs.Store(append(subbed.([][]byte), data))
	}
	return nil
}
func (ws *WsConn) SendMessage(msg []byte) {
	ws.writeBufferChan <- msg
}
func (ws *WsConn) SendPingMessage(msg []byte) {
	ws.pingMessageBufferChan <- msg
}
func (ws *WsConn) SendPongMessage(msg []byte) {
	ws.pongMessageBufferChan <- msg
}
func (ws *WsConn) SendCloseMessage(msg []byte) {
	ws.closeMessageBufferChan <- msg
}
func (ws *WsConn) SendJsonMessage(m any) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	ws.writeBufferChan <- data
	return nil
}
func (ws *WsConn) receiveMessage() {
	//exit
	ws.c.SetCloseHandler(func(code int, text string) error {
		log.Warn().Msgf("[ws][%s] websocket exiting [code=%d , text=%s]", ws.WsUrl, code, text)
		// ws.CloseWs()
		return nil
	})
	ws.c.SetPongHandler(func(pong string) error {
		log.Debug().Msgf("[%s] received [pong] %s", ws.WsUrl, pong)
		ws.c.SetReadDeadline(time.Now().Add(ws.readDeadLineTime))
		return nil
	})
	ws.c.SetPingHandler(func(ping string) error {
		log.Debug().Msgf("[%s] received [ping] %s", ws.WsUrl, ping)
		ws.SendPongMessage([]byte(ping))
		ws.c.SetReadDeadline(time.Now().Add(ws.readDeadLineTime))
		return nil
	})
	for {
		select {
		case <-ws.close:
			log.Info().Msgf("[ws][%s] close websocket , exiting receive message goroutine.", ws.WsUrl)
			return
		default:
			t, msg, err := ws.c.ReadMessage()
			if err != nil {
				log.Error().Msgf("[ws][%s] %s", ws.WsUrl, err.Error())
				if ws.IsAutoReconnect {
					log.Info().Msgf("[ws][%s] Unexpected Closed , Begin Retry Connect.", ws.WsUrl)
					ws.reconnect()
					continue
				}
				if ws.ErrorHandleFunc != nil {
					ws.ErrorHandleFunc(err)
				}
				return
			}
			//			logs.D(slice.Bytes2String(msg))
			ws.c.SetReadDeadline(time.Now().Add(ws.readDeadLineTime))
			switch t {
			case websocket.TextMessage:
				ws.ProtoHandleFunc(msg)
			case websocket.BinaryMessage:
				if ws.DecompressFunc == nil {
					ws.ProtoHandleFunc(msg)
				} else {
					if msg2, err := ws.DecompressFunc(msg); err != nil {
						log.Error().Msgf("[ws][%s] decompress error %s", ws.WsUrl, err.Error())
					} else {
						ws.ProtoHandleFunc(msg2)
					}
				}
				//	case websocket.CloseMessage:
				//	ws.CloseWs()
			default:
				log.Error().Bytes("msg data", msg).Str("url", ws.WsUrl).Msg("[ws] websocket message type")
			}
		}
	}
}
func (ws *WsConn) CloseWs() {
	//ws.close <- true
	close(ws.close)
	close(ws.writeBufferChan)
	close(ws.closeMessageBufferChan)
	close(ws.pingMessageBufferChan)
	close(ws.pongMessageBufferChan)
	err := ws.c.Close()
	zelo.OnErr(err).Str("url", ws.WsUrl).Msg("[ws] close websocket error")
}
func (ws *WsConn) clearChannel(c chan struct{}) {
	for {
		if len(c) > 0 {
			<-c
		} else {
			break
		}
	}
}
