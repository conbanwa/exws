package q

import (
	"github.com/conbanwa/logs"
	"github.com/conbanwa/wstrader/stat/zelo"
	"strings"
	"sync"
)

type D struct {
	Base, Quote string
}
type Attr struct {
	// OrderLimit uint32
	MultiOrder, Round bool
}

func (d D) same(d2 D) bool {
	return d.Base == d2.Quote && d2.Base == d.Quote || (d.Quote == d2.Quote && d.Base == d2.Base)
}

func (d D) Onward(side bool) D {
	if side {
		return d
	}
	return d.Reverse()
}
func (d D) Check(v Bbo, BBOs *sync.Map) bool {
	if !d.validSlug() {
		logs.E("Duo not valid", d, v)
		return false
	}
	if !v.Valid() {
		logs.E("BBO not valid", d, v)
		return false
	}
	x, xs, f := ras(bothBbo(BBOs, d))
	y, ys, b := ras(bothBbo(BBOs, d.Reverse()))
	if f == b {
		if !d.same(D{"BTC", "USDC"}) {
			logs.I("same side in both trade ", d, v, v.Bid, v.Ask, f, b)
		}
	} else if x*y > 1.002 {
		logs.F("EARN in single trade ", d, v, v.Bid, v.Ask, x*y)
	}
	if x*y*xs*ys == 0 {
		logs.F("0 in single trade ", d, v, v.Bid, v.Ask, x, xs, y, ys)
	}
	return true
}
func ras(bid, ask Bbo) (rate, size float64, front bool) {
	if bid.Bid != 0 {
		rate = bid.Bid
		size = bid.BidSize
		front = true
	} else if ask.Ask != 0 {
		rate = 1 / ask.Ask
		size = ask.AskSize * ask.Ask
	}
	return
}
func bothBbo(bbo *sync.Map, d D) (Bbo, Bbo) {
	bid, bk := bbo.Load(d)
	ask, ak := bbo.Load(d.Reverse())
	if bk {
		return bid.(Bbo), Bbo{}
	}
	if ak {
		return Bbo{}, ask.(Bbo)
	}
	logs.F(d, "both direction have no BBO")
	return Bbo{}, Bbo{}
}
func (d D) validSlug() bool {
	isValid := validCurrency(d.Base, 49) && validCurrency(d.Quote, 39) && d.Base != d.Quote
	if len(d.Quote) < 2 || len(d.Base) < 2 && d.Quote != "T" && d.Base == "T" {
		zelo.Writer.Info().Str("one letter pair", d).Send()
	}
	// if d.Quote == "tether" {
	// 	return true
	// }
	// if IsUpper(d.Quote) && IsUpper(d.Base) {
	// 	logs.W(d, strings.Contains(d.Quote, "USD"))
	// }
	zelo.NotEqual(isValid, true).Msg("not validSlug")
	return isValid
}
func (d D) Valid() bool {
	isValid := validCurrency(d.Base, 19) && validCurrency(d.Quote, 19) && d.Base != d.Quote
	zelo.NotEqual(isValid, true).Msg("not valid")
	return isValid
}
func validCurrency(currency string, limit int) bool {
	return len(currency) >= 1 && len(currency) < limit && currency != "UNKNOWN"
}

func (d D) Reverse() D {
	return D{Base: d.Quote, Quote: d.Base}
}
func (d D) Clip(s string) string {
	return d.Base + s + d.Quote
}
func (d D) String() string {
	return d.Clip("-")
}
func (d D) Contains(currency string) bool {
	return strings.Contains(d.Quote, currency) || strings.Contains(d.Base, currency)
}
func (d D) QuoteContains(currency string) bool {
	return strings.Contains(d.Quote, currency)
}
func (d D) GetRelatedVs(cl any, ok bool) (relatedVs map[T][3]bool, err error) {
	if cl == nil || !ok {
		err = logs.Err(d, "not in cluster")
		return
	}
	relatedVs, ok = cl.(map[T][3]bool)
	if relatedVs == nil || !ok {
		err = logs.Err(cl, "cluster is not Tri")
		return
	}
	zelo.PanicNotGreater(len(relatedVs), 0).Msg("not in cluster")
	return
}

func (d D) ToInstrument(duoSym map[D]string) string {
	if instruments, ok := duoSym[d]; ok {
		return instruments
	}
	panic(d.Base + " " + d.Quote)
}
func DsToInstruments(ds []D, duoSym map[D]string) (instruments []string) {
	for _, d := range ds {
		instruments = append(instruments, d.ToInstrument(duoSym))
	}
	return
}
