package lap

import (
	"testing"
	"time"
)

var earn Earn

func TestNewElapse(t *testing.T) {
	for i := 12; i < 15; i++ {
		mission(i)
	}
}

func mission(v int) {
	elapse := NewElapse()
	for i := 0; i < v-10; i++ {
		time.Sleep(10 * time.Millisecond * time.Duration(i))
		elapse = elapse.PushNow()
	}
	time.Sleep(100 * time.Millisecond)
	elapse = elapse.PushNow()
	time.Sleep(100 * time.Millisecond)
	earn.Mark((150-float64(v))/100, float64(v)/100, elapse)
	earn.Speak()
}
