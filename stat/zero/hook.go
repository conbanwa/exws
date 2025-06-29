package zero

import (
	"cmp"

	. "github.com/rs/zerolog"
)

func ErrHook(err error) HookFunc {
	return func(e *Event, level Level, message string) {
		if err == nil {
			e.Discard()
		}
		e.Err(err)
	}
}
func EqualHook[T comparable](a, b T) HookFunc {
	return func(e *Event, level Level, message string) {
		if a == b {
			e.Discard()
		}
		e.Any("equals", a).Any("expect", b)
	}
}
func GreaterHook[T cmp.Ordered](a, b T) HookFunc {
	return func(e *Event, level Level, message string) {
		if a > b {
			e.Discard()
		}
		e.Any("less", a).Any("greater", b)
	}
}
func LessHook(a, b float64) HookFunc {
	return func(e *Event, level Level, message string) {
		if a >= b {
			e.Discard()
		}
		e.Float64("less", a).Float64("greater", b)
	}
}

func OnErrWithLevel(err error, level Level) (e *Event) {
	hf := ErrHook(err)
	logger := Writer.Hook(hf)
	return lv(err == nil, level, logger)
}
func NotEqualWithLevel[T comparable](a, b T, level Level) (e *Event) {
	hf := EqualHook(a, b)
	logger := Writer.Hook(hf)
	return lv(a == b, level, logger)
}
func NotGreaterWithLevel[T cmp.Ordered](a, b T, level Level) (e *Event) {
	hf := GreaterHook(a, b)
	logger := Writer.Hook(hf)
	return lv(a > b, level, logger)
}
func LessWithLevel(a, b float64, level Level) (e *Event) {
	hf := LessHook(a, b)
	logger := Writer.Hook(hf)
	return lv(a >= b, level, logger)
}

func lv(discard bool, level Level, logger Logger) *Event {
	if !discard && level >= FatalLevel {
		//if level == FatalLevel {
		//	return logger.Fatal()
		//}
		return logger.Panic()
	}
	return logger.WithLevel(level)
}

func OnErr(err error) (e *Event) {
	return OnErrWithLevel(err, ErrorLevel)
}
func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func NotEqual[T comparable](a, b T) (e *Event) {
	return NotEqualWithLevel(a, b, ErrorLevel)
}
func Assert[T comparable](a, b T) (e *Event) {
	return NotEqualWithLevel(a, b, PanicLevel)
}

func NotGreater[T cmp.Ordered](a, b T) (e *Event) {
	return NotGreaterWithLevel(a, b, ErrorLevel)
}
func PanicNotGreater[T cmp.Ordered](a, b T) (e *Event) {
	return NotGreaterWithLevel(a, b, PanicLevel)
}
func Less(a, b float64) (e *Event) {
	return LessWithLevel(a, b, ErrorLevel)
}
func PanicLess(a, b float64) (e *Event) {
	return LessWithLevel(a, b, PanicLevel)
}
