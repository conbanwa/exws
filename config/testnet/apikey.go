package testnet

import "github.com/conbanwa/wstrader/cons"

var (
	huobi   = [4]string{}
	Binance = [4]string{}
	Gateio  = [4]string{}
)

func Keys(e string) (s [4]string) {
	switch e {
	case cons.HUOBI:
		return huobi
	case cons.BINANCE:
		return Binance
	case cons.GATEIO:
		return Gateio
	case "gatev2":
		return Gateio
	}
	return
}
