package q

// Bbo ticker
type Bbo struct {
	Bid, BidSize float64
	Pair         string
	Ask, AskSize float64
	Updated      int64
	// D
}

func (bb Bbo) Valid() bool {
	return bb.AskSize > 0 && bb.Bid > 0 && bb.BidSize > 0 && bb.Ask >= bb.Bid
}
func (bb Bbo) SideRatioAmount(sell bool) (ratio float64, amount float64) {
	if sell {
		return bb.Bid, bb.BidSize
	}
	return 1 / bb.Ask, bb.AskSize * bb.Ask
}
func (bb Bbo) SideRatio(sell bool) float64 {
	if sell {
		return bb.Bid
	}
	return 1 / bb.Ask
}

// Deprecated: by chan
func ExpireTime(bl any, old Bbo) int64 {
	return bl.(Bbo).Updated - old.Updated
}
