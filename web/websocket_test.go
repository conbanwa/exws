package web

import (
	"encoding/json"
	"github.com/conbanwa/slice"
	"qa3/wstrader/config"
	"testing"
	"time"
)

func Test_time(t *testing.T) {
	t.Log(time.Now().Unix())
}
func ProtoHandle(data []byte) error {
	println(slice.Bytes2String(data))
	return nil
}
func TestNewWsConn(t *testing.T) {
	clientId := "a"
	args := make([]any, 0)
	var heartbeatFunc = func() []byte {
		ts := time.Now().Unix()*1000 + 42029
		args = append(args, ts)
		//ping := fmt.Sprintf("{\"cmd\":\"ping\",\"args\":[%d],\"id\":\"%s\"}", ts, clientId)
		ping2 := map[string]any{
			"cmd":  "ping",
			"id":   clientId,
			"args": args}
		ping3, _ := json.Marshal(ping2)
		return ping3
	}
	//logs.D(ping)
	ws := NewWsBuilder().Dump().WsUrl("wss://api.fcoin.com/v2/ws").
		ProxyUrl("socks5://"+config.Proxy).
		AutoReconnect().
		Heartbeat(heartbeatFunc, 5*time.Second).ProtoHandleFunc(ProtoHandle).Build()
	t.Log(ws.Subscribe(map[string]string{
		"cmd": "sub", "args": "ticker.btcusdt", "id": clientId}))
	time.Sleep(time.Second * 20)
	ws.c.Close()
	time.Sleep(time.Second * 120)
}
