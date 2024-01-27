package zelo

import (
	"fmt"
	"os"
	"testing"
)

var Err = fmt.Errorf("test")

func init() {
	file, err := os.OpenFile("notes.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		Writer.Fatal().Err(err)
	}
	Writer = Json(file)
}

func TestOut(t *testing.T) {
	f := 1.0
	Writer = Colored(os.Stderr)
	Writer.Error().Discard().Float64("1", f).Msg("")
	NotGreater(1, 0).Int("1", 1).Msg("never")
	PanicLess(1, 1).Int("1", 1).Msg("NotEqualWithLevel")
	Assert(1, 1).Int("1", 1).Msg("never")
	NotEqual(1, 3).Int("1", 1).Msg("NotEqualWithLevel")
	NotEqual(f+2, f+3).Float64("1", f).Msg("")
	Writer = Plain(os.Stderr, false)
	Less(f+2, f+3).Float64("1", f).Msg("")
	NotGreater(f+2, f+3).Float64("1", f).Msg("")
	Writer = Colored(os.Stderr)
	Writer.Error().Err(Err).Float64("1", f).Msg("")
	Writer.Err(Err).Float64("1", f).Msg("")
	Writer.Err(nil).Float64("1", f).Msg("")
	//Writer.Debug().Bytes("22", mp.m.Range()).Send()
}

func BenchmarkNop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := float64(i)
		//if f+2 < f+3 {
		Writer.Error().Float64("1", f).Msg("")
		//}
	}
}
func BenchmarkErr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := float64(i)
		if f+2 < f+3 {
			if Err != nil {
				Writer.Err(Err).Float64("1", f).Msg("")
			}
		}
	}
}
func BenchmarkOnErr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := float64(i)
		if f+2 < f+3 {
			OnErr(Err).Float64("1", f).Msg("")
		}
	}
}
func BenchmarkIfNotEqual(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f := float64(i)
		NotEqual(f+2, f+3).Float64("1", f).Msg("")
	}
}
func BenchmarkLessThan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := float64(i)
		Less(f+2, f+3).Float64("1", f).Msg("")
	}
}

func BenchmarkLess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := float64(i)
		NotGreater(f+2, f+3).Float64("1", f).Msg("")
	}
}
