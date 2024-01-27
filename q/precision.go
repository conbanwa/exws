package q

import (
	"github.com/conbanwa/logs"
	"github.com/conbanwa/num"
)

type P struct {
	Base, BaseLimit, Quote, Price float64 //step
	// BaseN, BaseLimitN, QuoteN, PriceN int
	MinBase, MinBaseLimit, MinQuote float64
	TakerNett, MakerNett            float64
}

func (p P) GetMinSpend(bid bool) map[bool]float64 {
	if bid {
		return map[bool]float64{true: p.Base, false: p.Quote}
	}
	return map[bool]float64{true: p.Quote, false: p.Base}
}

func validateMinSize(fixedSize, min float64) error {
	if fixedSize < min {
		return logs.Err(fixedSize, "<", min)
	}
	return nil
}
func findFixedLost(size float64, fixedSize float64) float64 {
	return (size - fixedSize) / fixedSize
}
func (prec P) GetFixedSize(ratio, spend float64, sell bool) (sAmount string, fixedBaseSize, fixedQuoteSize, fixedLost float64, isForceMarket bool, err error) {
	baseSize, quoteSize := spend*ratio, spend
	if sell {
		baseSize, quoteSize = spend, spend*ratio
	}
	fixedBaseSize = num.FloatToFixed(baseSize, prec.Base)
	fixedQuoteSize = num.FloatToFixed(quoteSize, prec.Quote)
	if err = validateMinSize(fixedBaseSize, prec.MinBase); err != nil {
		return
	}
	if err = validateMinSize(fixedQuoteSize, prec.MinQuote); err != nil {
		return
	}
	if err := validateMinSize(fixedBaseSize, prec.MinBaseLimit); err != nil {
		logs.D("ForceMarket", sell, spend, fixedBaseSize, "<", prec.MinBaseLimit)
		isForceMarket = true
	}
	if sell {
		fixedLost = findFixedLost(quoteSize, fixedQuoteSize)
	} else {
		fixedLost = findFixedLost(baseSize, fixedBaseSize)
	}
	// tofixed TODO decrease lost for if inc = numerical.FloatToFixed((dec-incStep)*bbo.Bid*p.TakerNett, incStep)
	logs.I(prec.TakerNett, sell, spend, fixedBaseSize, fixedQuoteSize, ratio)
	sAmount = num.FloatToString(baseSize, prec.Base)
	return
}

func DefaultPrecision() P {
	return P{
		Base:         1e-06,
		BaseLimit:    1e-06,
		Quote:        1e-06,
		Price:        1e-06,
		MinBase:      1e-02,
		MinBaseLimit: 1e-02,
		MinQuote:     1e-02,
		TakerNett:    0.99925}
}

func NoFeePrecision() P {
	return P{
		Base:         1e-04,
		BaseLimit:    1e-04,
		Quote:        1e-04,
		Price:        1e-04,
		MinBase:      5,
		MinBaseLimit: 5,
		MinQuote:     5,
		TakerNett:    1.0}
}

//	type Balance struct {
//		Available float64
//		Frozen    float64
//      Loan      float64
// }
