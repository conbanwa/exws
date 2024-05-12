package binance

import (
	"fmt"
	"qa3/wstrader/cons"
	"strings"
)

func adaptStreamToCurrencyPair(stream string) cons.CurrencyPair {
	symbol := strings.Split(stream, "@")[0]
	if strings.HasSuffix(symbol, "usdt") {
		return cons.NewCurrencyPair2(fmt.Sprintf("%s_usdt", strings.TrimSuffix(symbol, "usdt")))
	}
	if strings.HasSuffix(symbol, "usd") {
		return cons.NewCurrencyPair2(fmt.Sprintf("%s_usd", strings.TrimSuffix(symbol, "usd")))
	}
	if strings.HasSuffix(symbol, "btc") {
		return cons.NewCurrencyPair2(fmt.Sprintf("%s_btc", strings.TrimSuffix(symbol, "btc")))
	}
	return cons.UNKNOWN_PAIR
}
func adaptSymbolToCurrencyPair(symbol string) cons.CurrencyPair {
	symbol = strings.ToUpper(symbol)
	if strings.HasSuffix(symbol, "USD") {
		return cons.NewCurrencyPair2(fmt.Sprintf("%s_USD", strings.TrimSuffix(symbol, "USD")))
	}
	if strings.HasSuffix(symbol, "USDT") {
		return cons.NewCurrencyPair2(fmt.Sprintf("%s_USDT", strings.TrimSuffix(symbol, "USDT")))
	}
	if strings.HasSuffix(symbol, "PAX") {
		return cons.NewCurrencyPair2(fmt.Sprintf("%s_PAX", strings.TrimSuffix(symbol, "PAX")))
	}
	if strings.HasSuffix(symbol, "BTC") {
		return cons.NewCurrencyPair2(fmt.Sprintf("%s_BTC", strings.TrimSuffix(symbol, "BTC")))
	}
	return cons.UNKNOWN_PAIR
}
func adaptOrderStatus(status string) cons.TradeStatus {
	var tradeStatus cons.TradeStatus
	switch status {
	case "NEW":
		tradeStatus = cons.ORDER_UNFINISH
	case "FILLED":
		tradeStatus = cons.ORDER_FINISH
	case "PARTIALLY_FILLED":
		tradeStatus = cons.ORDER_PART_FINISH
	case "CANCELED":
		tradeStatus = cons.ORDER_CANCEL
	case "PENDING_CANCEL":
		tradeStatus = cons.ORDER_CANCEL_ING
	case "REJECTED":
		tradeStatus = cons.ORDER_REJECT
	}
	return tradeStatus
}
