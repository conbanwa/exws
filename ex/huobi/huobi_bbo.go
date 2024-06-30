package huobi

import (
	"errors"
	"fmt"
)

type BBOResponse struct {
	Symbol                                       string
	QuoteTime, Bid, BidSize, Ask, AskSize, SeqId float64
}

func (ws *SpotWs) SubscribeBBO(sm []string) (err error) {
	if ws.bboCallback == nil {
		return errors.New("please set bbo callback func")
	}
	for _, sym := range sm {
		err = ws.subscribe(map[string]any{
			"id":  "spot.bbo",
			"sub": fmt.Sprintf("market.%s.bbo", sym),
		})
		if err != nil {
			log.Panic().Err(err).Send()
			return
		}
	}
	return nil
}
