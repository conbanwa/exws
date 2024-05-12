package util

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/url"
	"qa3/wstrader"
	"qa3/wstrader/cons"
	"strings"
	"unicode"
)

func MergeOptionalParameter(values *url.Values, opts ...wstrader.OptionalParameter) url.Values {
	for _, opt := range opts {
		for k, v := range opt {
			values.Set(k, fmt.Sprint(v))
		}
	}
	return *values
}

func GzipDecompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

func FlateDecompress(data []byte) ([]byte, error) {
	return ioutil.ReadAll(flate.NewReader(bytes.NewReader(data)))
}

func GenerateOrderClientId(size int) string {
	uuidStr := strings.Replace(uuid.New().String(), "-", "", 32)
	return "q3a" + uuidStr[0:size-5]
}

func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func AdaptTradeSide(side string) cons.TradeSide {
	side2 := strings.ToUpper(side)
	switch side2 {
	case "SELL":
		return cons.SELL
	case "BUY":
		return cons.BUY
	case "BUY_MARKET":
		return cons.BUY_MARKET
	case "SELL_MARKET":
		return cons.SELL_MARKET
	default:
		return -1
	}
}
