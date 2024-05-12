package ftx

import (
	. "github.com/conbanwa/wstrader"
	. "github.com/conbanwa/wstrader/cons"
	"strings"
)

const (
	baseUrl = "https://capi.ftx.com"
	LILY    = "lily"
)

func (client *Client) String() string {
	return FTX
}

type Instrument struct {
	Coin                string `json:"coin"`
	ContractVal         string `json:"contract_val"`
	Delivery            []any  `json:"delivery"`
	ForwardContractFlag bool   `json:"forwardContractFlag"`
	Listing             any    `json:"listing"`
	PriceEndStep        int    `json:"priceEndStep"`
	QuoteCurrency       string `json:"quote_currency"`
	SizeIncrement       int    `json:"size_increment"`
	Symbol              string `json:"symbol"`
	TickSize            int    `json:"tick_size"`
	UnderlyingIndex     string `json:"underlying_index"`
}

//	func (ftx *FtxClient) GetInstruments() ([]Instrument, error) {
//		url := fmt.Sprintf("%s/api/swap/v3/market/contracts", baseUrl)
//		resp, err := HttpGet3(ftx.httpClient, url, nil)
//		if err != nil {
//			return nil, err
//		}
//		ins := make([]Instrument, 0)
//		for _, v := range resp {
//			contract := v.(map[string]any)
//			ins = append(ins, Instrument{
//				Coin:                contract["coin"].(string),
//				ContractVal:         contract["contract_val"].(string),
//				Delivery:            contract["delivery"].([]any),
//				ForwardContractFlag: contract["forwardContractFlag"].(bool),
//				PriceEndStep:        num.ToInt[int](contract["priceEndStep"]),
//				QuoteCurrency:       contract["quote_currency"].(string),
//				SizeIncrement:       num.ToInt[int](contract["size_increment"]),
//				Symbol:              contract["contract_val"].(string),
//				TickSize:            num.ToInt[int](contract["tick_size"]),
//				UnderlyingIndex:     contract["underlying_index"].(string),
//			})
//		}
//		return ins, nil
//	}
func (client *Client) GetAccount() (a *Account, err error) {
	return
}
func symbolToCurrencyPair(symbol string) CurrencyPair {
	currencyA := strings.ToUpper(symbol[0:3])
	currencyB := strings.ToUpper(symbol[3:])
	return NewCurrencyPair(NewCurrency(currencyA, ""), NewCurrency(currencyB, ""))
}
