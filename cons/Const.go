package cons

import (
	"fmt"
)

type (
	TradeSide    int
	TradeStatus  int
	KlinePeriod  int
	OrderFeature int
	OrderType    int
)

const (
	BUY TradeSide = 1 + iota
	SELL
	BUY_MARKET
	SELL_MARKET
)

func (ts TradeSide) String() string {
	switch ts {
	case 1:
		return "BUY"
	case 2:
		return "SELL"
	case 3:
		return "BUY_MARKET"
	case 4:
		return "SELL_MARKET"
	default:
		return "UNKNOWN"
	}
}

func (ts TradeStatus) String() string {
	return tradeStatusSymbol[ts]
}

var tradeStatusSymbol = [...]string{"UNFINISH", "PART_FINISH", "FINISH", "CANCEL", "REJECT", "CANCEL_ING", "FAIL"}

const (
	ORDER_UNFINISH TradeStatus = iota
	ORDER_PART_FINISH
	ORDER_FINISH
	ORDER_CANCEL
	ORDER_REJECT
	ORDER_CANCEL_ING
	ORDER_FAIL
)
const (
	OPEN_BUY   = 1 + iota //开多
	OPEN_SELL             //开空
	CLOSE_BUY             //平多
	CLOSE_SELL            //平空
)

// k线周期
const (
	KLINE_PERIOD_1MIN = 1 + iota
	KLINE_PERIOD_3MIN
	KLINE_PERIOD_5MIN
	KLINE_PERIOD_15MIN
	KLINE_PERIOD_30MIN
	KLINE_PERIOD_60MIN
	KLINE_PERIOD_1H
	KLINE_PERIOD_2H
	KLINE_PERIOD_3H
	KLINE_PERIOD_4H
	KLINE_PERIOD_6H
	KLINE_PERIOD_8H
	KLINE_PERIOD_12H
	KLINE_PERIOD_1DAY
	KLINE_PERIOD_3DAY
	KLINE_PERIOD_1WEEK
	KLINE_PERIOD_1MONTH
	KLINE_PERIOD_1YEAR
)

const (
	ORDER_FEATURE_ORDINARY = iota
	ORDER_FEATURE_POST_ONLY
	ORDER_FEATURE_FOK
	ORDER_FEATURE_IOC
	ORDER_FEATURE_FAK
	ORDER_FEATURE_LIMIT
)

func (of OrderFeature) String() string {
	if of > 0 && int(of) < len(orderFeatureSymbol) {
		return orderFeatureSymbol[of]
	}
	return fmt.Sprintf("UNKNOWN_ORDER_TYPE(%d)", of)
}

var orderFeatureSymbol = [...]string{"ORDINARY", "POST_ONLY", "FOK", "IOC", "FAK", "LIMIT"}

func (ot OrderType) String() string {
	if ot > 0 && int(ot) <= len(orderTypeSymbol) {
		return orderTypeSymbol[ot-1]
	}
	return fmt.Sprintf("UNKNOWN_ORDER_TYPE(%d)", ot)
}

var orderTypeSymbol = [...]string{"LIMIT", "MARKET"}

const (
	ORDER_TYPE_LIMIT = 1 + iota
	ORDER_TYPE_MARKET
)

const (
	SUB_ACCOUNT         = iota         //子账户
	SPOT                               // 币币交易
	FUTURE              = 1 + iota     //交割合约
	C2C                                //法币
	SPOT_MARGIN                        //币币杠杆交易
	WALLET                             // 资金账户
	TIPS                = 2 + iota     //余币宝
	SWAP                               //永续合约
	SWAP_USDT                          //usdt本位永续合约
	THIS_WEEK_CONTRACT  = "this_week"  //周合约
	NEXT_WEEK_CONTRACT  = "next_week"  //次周合约
	QUARTER_CONTRACT    = "quarter"    //季度合约
	BI_QUARTER_CONTRACT = "bi_quarter" // NEXT QUARTER
	SWAP_CONTRACT       = "swap"       //永续合约
	SWAP_USDT_CONTRACT  = "swap-usdt"
	// exchanges const
	KUCOIN          = "kucoin.com"
	OKCOIN_COM      = "okcoin.com"
	OKEX            = "okx.com"
	OKEX_FUTURE     = "okx.com_future"
	OKEX_SWAP       = "okx.com_swap"
	HUOBI           = "huobi.com"
	HUOBI_PRO       = "huobi.pro"
	BITSTAMP        = "bitstamp.net"
	KRAKEN          = "kraken.com"
	BITFINEX        = "bitfinex.com"
	BINANCE         = "binance.com"
	BINANCE_SWAP    = "binance.com_swap"
	BINANCE_FUTURES = "binance.com_futures"
	POLONIEX        = "poloniex.com"
	COINEX          = "coinex.com"
	BITHUMB         = "bithumb.com"
	GATEIO          = "gate.io"
	COINBASE            = "coinbase.com"
	BIGONE          = "big.one"
	FMEX            = "fmex.com"
	HITBTC          = "hitbtc.com"
	BITMEX          = "bitmex.com"
	BITMEX_TEST     = "testnet.bitmex.com"
	CRYPTOPIA       = "cryptopia.co.nz"
	HBDM            = "hbdm.com"
	HBDM_SWAP       = "hbdm.com_swap"
	BITGET_SWAP     = "bitget_swap"
)

type LimitOrderOptionalParameter int

const (
	PostOnly LimitOrderOptionalParameter = iota + 1
	Ioc
	Fok
	FuturesTwoWayPositionMode //币安双向持仓模式
)

func (opt LimitOrderOptionalParameter) String() string {
	switch opt {
	case PostOnly:
		return "post_only"
	case Fok:
		return "fok"
	case Ioc:
		return "ioc"
	default:
		return "error-order-optional-parameter"
	}
}
