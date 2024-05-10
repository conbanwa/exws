package lap

import (
	"sync/atomic"
	"time"

	"github.com/conbanwa/logs"
)

type Elapse []int64
type Earn struct {
	Elapse
	Estimates, Profits, Duration, Count atomic.Int64
}

const MicroMultiplier = 1e6

func (earn *Earn) Mark(e, p float64, elapse Elapse) {
	elapse = elapse.PushNow()
	if len(earn.Elapse) < len(elapse) {
		earn.Elapse = make(Elapse, len(elapse))
	}
	earn.Count.Add(1)
	earn.Estimates.Add(int64(e * MicroMultiplier))
	earn.Profits.Add(int64(p * MicroMultiplier))
	earn.Duration.Add(elapse.Sum())
	earn.Elapse = elapse
	logs.I("estimate grow", e, "profit", p, " cycle time: ", elapse.Gap())
}
func (earn *Earn) Speak() {
	logs.I(" -  -  -  matches", earn.Count,
		"Estimates", earn.Average(earn.Estimates),
		"Profits", earn.Average(earn.Profits),
		"MeanTime", earn.Average(earn.Duration), "s", earn.Elapse)
	logs.I("Gap", earn.Elapse.Gap())
}

func (earn *Earn) Average(i atomic.Int64) float64 {
	return float64(i.Load()/earn.Count.Load()) / MicroMultiplier
}
func NewElapse() Elapse {
	return Elapse{time.Now().UnixMicro()}
}
func (elapse Elapse) PushNow() Elapse {
	return append(elapse, time.Now().UnixMicro()-elapse[0])
}
func (elapse Elapse) Sum() int64 {
	return elapse[len(elapse)-1]
}
func (elapse Elapse) RawGap() Elapse {
	gap := make(Elapse, len(elapse))
	for i := 2; i < len(elapse); i++ {
		gap[i] = elapse[i] - elapse[i-1]
	}
	gap[0] = elapse.Sum()
	return gap
}
func (elapse Elapse) Gap() []float64 {
	dived := make([]float64, len(elapse))
	for i, v := range elapse.RawGap() {
		dived[i] = float64(v) / MicroMultiplier
	}
	return dived
}
