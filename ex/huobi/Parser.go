package huobi

import (
	"fmt"
	"github.com/conbanwa/exws"
	"github.com/conbanwa/exws/cons"
	"sort"
	"strings"
)

func parseDepthFromResponse(r DepthResponse) exws.Depth {
	var dep exws.Depth
	for _, bid := range r.Bids {
		dep.BidList = append(dep.BidList, exws.DepthRecord{Price: bid[0], Amount: bid[1]})
	}
	for _, ask := range r.Asks {
		dep.AskList = append(dep.AskList, exws.DepthRecord{Price: ask[0], Amount: ask[1]})
	}
	sort.Sort(sort.Reverse(dep.BidList))
	sort.Sort(sort.Reverse(dep.AskList))
	return dep
}
func parseCurrencyPairFromSpotWsCh(ch string) cons.CurrencyPair {
	if ch == "" {
		return cons.UNKNOWN_PAIR
	}
	meta := strings.Split(ch, ".")
	if len(meta) < 2 {
		log.Error().Msgf("parse error, ch=%s", ch)
		return cons.UNKNOWN_PAIR
	}
	currencyPairStr := meta[1]
	if strings.HasSuffix(currencyPairStr, "usdt") {
		currencyA := strings.TrimSuffix(currencyPairStr, "usdt")
		return cons.NewCurrencyPair2(fmt.Sprintf("%s_usdt", currencyA))
	}
	if strings.HasSuffix(currencyPairStr, "btc") {
		currencyA := strings.TrimSuffix(currencyPairStr, "btc")
		return cons.NewCurrencyPair2(fmt.Sprintf("%s_btc", currencyA))
	}
	if strings.HasSuffix(currencyPairStr, "eth") {
		currencyA := strings.TrimSuffix(currencyPairStr, "eth")
		return cons.NewCurrencyPair2(fmt.Sprintf("%s_eth", currencyA))
	}
	if strings.HasSuffix(currencyPairStr, "husd") {
		currencyA := strings.TrimSuffix(currencyPairStr, "husd")
		return cons.NewCurrencyPair2(fmt.Sprintf("%s_husd", currencyA))
	}
	if strings.HasSuffix(currencyPairStr, "ht") {
		currencyA := strings.TrimSuffix(currencyPairStr, "ht")
		return cons.NewCurrencyPair2(fmt.Sprintf("%s_ht", currencyA))
	}
	if strings.HasSuffix(currencyPairStr, "trx") {
		currencyA := strings.TrimSuffix(currencyPairStr, "trx")
		return cons.NewCurrencyPair2(fmt.Sprintf("%s_trx", currencyA))
	}
	return cons.UNKNOWN_PAIR
}
