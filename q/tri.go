package q

import (
	"github.com/conbanwa/logs"
	"github.com/conbanwa/slice"
	"strings"
)

type T [3]string

func (t T) sortTri() T {
	var FIAT = []string{"DAI", "EUR", "TRY", "BRL", "ARS", "BIDR", "IDRT", "NGN", "PLN", "RON", "RUB", "UAH", "ZAR"}
	var ALTS = []string{"BNB", "ETH", "BTC", "TRX", "GBP"}
	var priority = append([]string{"US"}, FIAT...)
	priority = append(priority, ALTS...)
	for _, s := range priority {
		if strings.Contains(t[0], s) {
			return t
		}
		if strings.Contains(t[1], s) {
			return T{t[1], t[0], t[2]}
		}
		if strings.Contains(t[2], s) {
			return T{t[2], t[0], t[1]}
		}
	}
	return t
}

func (t T) Edge(i int, s bool) D {
	return D{Base: t[i], Quote: t[(i+1)%3]}.Onward(s)
}

func (t T) Foreach(v [3]bool, f func(D)) {
	for i := 0; i < 3; i++ {
		d := t.Edge(i, v[i])
		f(d)
	}
}

func (t T) Reverse() T {
	return T{t[0], t[2], t[1]}
}

func DirectVectors(p map[D]P) (TriArr map[T][3]bool) {
	TriArr = make(map[T][3]bool)
	ds := slice.MapKeys(p)
	length := len(p)
	for i := 0; i < length; i++ {
		for j := i + 1; j < length; j++ {
			i1 := ds[i].Base
			i2 := ds[i].Quote
			j1 := ds[j].Base
			j2 := ds[j].Quote
			target := D{Base: i2, Quote: j2}
			t3 := i1
			if i1 == j1 {
			} else if i1 == j2 {
				target.Quote = j1
			} else if i2 == j1 {
				target.Base = i1
				t3 = i2
			} else if i2 == j2 {
				target = D{Base: i1, Quote: j1}
				t3 = i2
			} else {
				continue
			}
			for k := j + 1; k < length; k++ {
				if strings.Contains(target.Base, "multi-collateral-dai") || strings.Contains(target.Base, "usd") {
					logs.W(target, "has tether")
					target.Base = "tether"
				}
				if ds[k].same(target) {
					if target.Base != "tether" && (!target.validSlug() || target.Base == t3 || target.Quote == t3) {
						logs.W(!target.validSlug(), target, t3)
						logs.F(i, ds[i:i+4], ds[k], target, T{target.Base, target.Quote, t3})
					}
					var tri = T{target.Base, target.Quote, t3}.sortTri()
					TriArr[tri] = tri.findVector(p)
					TriArr[tri.Reverse()] = tri.Reverse().findVector(p)
				}
			}
		}
	}
	return TriArr
}
func (t T) findVector(p map[D]P) (vector [3]bool) {
	for i := range vector {
		if _, ok := p[t.blindEdge(i)]; ok {
			vector[i] = true
		} else if _, ok := p[t.blindEdge(i).Reverse()]; ok {
		} else {
			logs.F(t, "no BBO", t.blindEdge(i))
		}
	}
	return
}
func (t T) blindEdge(i int) D {
	return D{Base: t[i], Quote: t[(i+1)%3]}
}

func (t T) Contain(d D) bool {
	var count = 0
	for _, cur := range []string{d.Base, d.Quote} {
		for _, ti := range t {
			if ti == cur {
				count++
			}
		}
	}
	return count == 2
}
func (t T) Has(s string) bool {
	for _, ti := range t {
		if ti == s {
			return true
		}
	}
	return false
}

func (t T) String() string {
	return "[" + t[0] + " " + t[1] + " " + t[2] + "]"
}

func GetEdges(vs map[T][3]bool) (ds []D) {
	return slice.MapKeys(GetEdgesMap(vs))
}

func GetEdgesMap(vs map[T][3]bool) map[D]bool {
	ds := map[D]bool{}
	for t, v := range vs {
		t.Foreach(v, func(d D) {
			ds[d] = true
		})
	}
	return ds
}
