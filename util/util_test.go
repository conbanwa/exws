package util

import (
	"github.com/conbanwa/logs"
	"github.com/conbanwa/num"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
	"testing"
)

func BenchmarkFloatToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		num.FloatToString(12.232321, 0.001)
	}
}

func TestFloatToString(t *testing.T) {
	assert.Equal(t, "1", num.FloatToString(1.10231000, 1))
	assert.Equal(t, "0.102", num.FloatToString(0.10231000, 0.001))
	assert.Equal(t, "189.61", num.FloatToString(189.61020000, 0.01))
	assert.NotEqual(t, "1.10231000", num.FloatToString(1.10231000, 1e-8))
	assert.Equal(t, 0.129999, num.FloatToFixed(0.1299999, 1e-6))
	logs.Assert(num.FloatToFixed(1.10231000, 0.000000001) != 1.10231000, num.FloatToFixed(1.10231000, 0.000000001))
	//logs.Assert(0.102-FloatToFixed(0.10231000, 0.001) < 0.00000000000001, FloatToFixed(0.10231000, 0.001))
	assert.Equal(t, num.FloatToString(1.10231000, 1), "1")
	//assert.Equal(t, num.FloatToString(1.10231000, 0), "1.10231")
	assert.Equal(t, num.FloatToString(0.10231000, 0.001), "0.102")
	assert.Equal(t, num.FloatToString(189.61020000, 0.01), "189.61")
	assert.Equal(t, num.FloatToString(1.10231000, 1e-7), "1.1023100")
	assert.Equal(t, num.FloatToString(0.1299999, 0.0001), "0.1299")
	assert.Equal(t, num.FloatToString(6.7597, 0.01), "6.75")
	assert.Equal(t, num.FloatToFixed(1.10231000, 1), 1.0)
	assert.Equal(t, num.FloatToFixed(1.10231000, 0), 1.10231)
	assert.Equal(t, num.FloatToFixed(189.61020000, 0.01), 189.61)
	logs.ErrorIfNotSame(num.FloatToFixed(0.10231000, 0.001), 0.102)
	logs.ErrorIfNotSame(num.FloatToString(1.10231000, 1e-8), "1.10231")
	logs.ErrorIfNotSame(num.FloatToFixed(1.10231000, 0.000000001), 1.10231)
	logs.ErrorIfNotSame(num.FloatToFixed(0.1299999, 0.0001), 0.1299)
	logs.ErrorIfNotSame(num.FloatToFixed(6.7597, 0.01), 6.75)
	type args struct {
		v    float64
		step float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "0",
			args: args{v: 1111.55, step: 0.1},
			want: "1111.5",
		},
		{
			name: "1",
			args: args{v: 3341.055, step: 0.01},
			want: "3341.05",
		},
		{
			name: "2",
			args: args{v: 61.0555, step: 0.001},
			want: "61.055",
		},
		{
			name: "3",
			args: args{v: 5551.0555, step: 10},
			want: "5550",
		},
		{
			name: "4",
			args: args{v: 441.0555, step: 100},
			want: "400",
		},
		{
			name: "5",
			args: args{v: 2.9999999, step: 0.1},
			want: "2.9",
		},
		{
			name: "6",
			args: args{v: 1.9999, step: 1},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := num.FloatToString(tt.args.v, tt.args.step); got != tt.want {
				t.Errorf("num.FloatToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloatToFixed(t *testing.T) {
	type args struct {
		v    float64
		step float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "0",
			args: args{v: 1.0555, step: 0.1},
			want: 1,
		},
		{
			name: "1",
			args: args{v: 1.0555, step: 0.01},
			want: 1.05,
		},
		{
			name: "2",
			args: args{v: 1.0555, step: 0.001},
			want: 1.055,
		},
		{
			name: "3",
			args: args{v: 5551.0555, step: 10},
			want: 5550,
		},
		{
			name: "4",
			args: args{v: 441.0555, step: 100},
			want: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, num.FloatToFixed(tt.args.v, tt.args.step), "FloatToFixed(%v, %v)", tt.args.v, tt.args.step)
		})
	}
}

func TestParseInt(t *testing.T) {
	type args struct {
		v any
	}
	type testCase[T interface{ constraints.Integer }] struct {
		name string
		args args
		want T
	}
	type integer int32
	tests := []testCase[integer]{
		{name: "1", args: args{v: 23}, want: 23},
		{name: "2", args: args{v: -23}, want: -23},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, num.ToInt[integer](tt.args.v), "ParseInt(%v)", tt.args.v)
		})
	}
}

func TestGenerateOrderClientId(t *testing.T) {
	t.Log(len(GenerateOrderClientId(32)), GenerateOrderClientId(32))
}
