package wallet

import (
	"github.com/conbanwa/wstrader/q"
)

type Msg struct {
	q.T
}

var LowBalances = map[Msg]int{}

var Ch = make(chan Msg, 3)

func init() {
	go ReportLowBalance()
}

func ReportLowBalance() {
	for {
		select {
		case tri := <-Ch:
			if tri.Has("BTC") {
				break
			}
			LowBalances[tri]++
		default:
		}
	}
}
