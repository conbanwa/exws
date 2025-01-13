package bitget

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/conbanwa/exws"
	"github.com/conbanwa/exws/cons"
	"github.com/conbanwa/exws/q"
	"github.com/conbanwa/exws/util"
	"github.com/conbanwa/exws/web"
	"github.com/conbanwa/num"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	baseUrl = "https://capi.bitget.com"
)

type BitgetSwap struct {
	accessKey  string
	secretKey  string
	passphrase string
	baseUrl    string
	httpClient *http.Client
	timeOffset int64
}

func NewSwap(config *exws.APIConfig) *BitgetSwap {
	if config.Endpoint == "" {
		config.Endpoint = baseUrl
	}
	bs := &BitgetSwap{
		baseUrl:    config.Endpoint,
		accessKey:  config.ApiKey,
		secretKey:  config.ApiSecretKey,
		passphrase: config.ApiPassphrase,
		httpClient: config.HttpClient,
	}
	bs.setTimeOffset()
	return bs
}
func (bs *BitgetSwap) SetBaseUri(uri string) {
	bs.baseUrl = uri
}
func (bs *BitgetSwap) String() string {
	return cons.BITGET_SWAP
}
func (bs *BitgetSwap) setTimeOffset() error {
	stime, err := bs.GetServerTime()
	if err != nil {
		return err
	}
	lt := time.Now().UnixNano() / 1000000
	bs.timeOffset = lt - stime
	return nil
}

/**
 *获取交割预估价
 */
func (bs *BitgetSwap) GetFutureEstimatedPrice(currencyPair cons.CurrencyPair) (float64, error) {
	panic("not supported.")
}

/**
 * 期货行情
 * @param currency_pair   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType  合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 */
func (bs *BitgetSwap) GetFutureTicker(currency cons.CurrencyPair, contractType string) (*exws.Ticker, error) {
	url := fmt.Sprintf("%s/api/swap/v1/instruments/%s/ticker", bs.baseUrl, bs.adaptSymbol(currency))
	tickerMap, err := web.HttpGet(bs.httpClient, url)
	if err != nil {
		return nil, err
	}
	status, isOk := tickerMap["status"]
	if !isOk || status != "ok" {
		return nil, errors.New(tickerMap["err_msg"].(string))
	}
	data := tickerMap["data"].(any)
	_data := data.(map[string]any)
	var ticker exws.Ticker
	ticker.Pair = currency
	ticker.Date = num.ToInt[uint64](_data["timestamp"])
	ticker.Last = num.ToFloat64(_data["last"])
	ticker.Buy = num.ToFloat64(_data["bidPrice"])
	ticker.Sell = num.ToFloat64(_data["best_ask"])
	ticker.High = num.ToFloat64(_data["high_24h"])
	ticker.Low = num.ToFloat64(_data["low_24h"])
	ticker.Vol = num.ToFloat64(_data["volume_24h"])
	return &ticker, nil
}

/**
 * 期货深度
 * @param currencyPair  btc_usd:比特币    ltc_usd :莱特币
 * @param contractType  合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @param size 获取深度档数
 * @return
 */
func (bs *BitgetSwap) GetFutureDepth(currency cons.CurrencyPair, contractType string, size int) (*exws.Depth, error) {
	panic("not implement")
}
func (bs *BitgetSwap) GetTrades(contractType string, currencyPair cons.CurrencyPair, since int64) ([]q.Trade, error) {
	panic("not implement")
}

/**
 * 期货指数
 * @param currencyPair   btc_usd:比特币    ltc_usd :莱特币
 */
func (bs *BitgetSwap) GetFutureIndex(currencyPair cons.CurrencyPair) (float64, error) {
	panic("not implement")
}
func (bs *BitgetSwap) doAuthRequest(method, uri string, param map[string]any) ([]byte, error) {
	timestamp := time.Now().Unix() * 1000
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["ACCESS-KEY"] = bs.accessKey
	headers["ACCESS-PASSPHRASE"] = bs.passphrase
	headers["ACCESS-TIMESTAMP"] = strconv.Itoa(int(timestamp))
	headers["locale"] = "zh-CN"
	postBody := ""
	if param != nil {
		postBodyArray, _ := json.Marshal(param)
		postBody = string(postBodyArray)
	}
	payload := fmt.Sprintf("%d%s%s%s", timestamp, method, uri, postBody)
	sign, _ := exws.GetParamHmacSHA256Base64Sign(bs.secretKey, payload)
	headers["ACCESS-SIGN"] = sign
	resp, err := web.NewRequest(bs.httpClient, method, bs.baseUrl+uri, postBody, headers)
	return resp, err
}

// 全仓账户
func (bs *BitgetSwap) GetFutureUserinfo(currencyPair ...cons.CurrencyPair) (*exws.FutureAccount, error) {
	if len(currencyPair) > 1 {
		panic("not support")
	}
	uri := "/api/swap/v3/account/accounts"
	if len(currencyPair) == 1 {
		uri = "/api/swap/v3/account/account?symbol=" + bs.adaptSymbol(currencyPair[0])
	}
	resp, err := bs.doAuthRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	type Account struct {
		Equity              float64 `json:"equity,string"`
		FixedBalance        string  `json:"fixed_balance"`
		ForwardContractFlag bool    `json:"forwardContractFlag"`
		Margin              string  `json:"margin"`
		MarginFrozen        string  `json:"margin_frozen"`
		MarginMode          string  `json:"margin_mode"`
		RealizedPnl         float64 `json:"realized_pnl,string"`
		Symbol              string  `json:"symbol"`
		Timestamp           string  `json:"timestamp"`
		TotalAvailBalance   float64 `json:"total_avail_balance,string"`
		UnrealizedPnl       float64 `json:"unrealized_pnl,string"`
	}
	subAccount := make(map[cons.Currency]exws.FutureSubAccount)
	if len(currencyPair) == 0 {
		accs := make([]Account, 0)
		err = json.Unmarshal(resp, &accs)
		if err != nil {
			return nil, err
		}
		for _, acc := range accs {
			currency := cons.Currency{
				Symbol: acc.Symbol,
				Desc:   "",
			}
			subAccount[currency] = exws.FutureSubAccount{
				Currency:      currency,
				AccountRights: acc.Equity,
				KeepDeposit:   acc.TotalAvailBalance,
				ProfitReal:    acc.RealizedPnl,
				ProfitUnreal:  acc.UnrealizedPnl,
				RiskRate:      0,
			}
		}
	} else {
		acc := Account{}
		err = json.Unmarshal(resp, &acc)
		if err != nil {
			return nil, err
		}
		currency := cons.Currency{
			Symbol: acc.Symbol,
			Desc:   "",
		}
		subAccount[currency] = exws.FutureSubAccount{
			Currency:      currency,
			AccountRights: acc.Equity,
			KeepDeposit:   acc.TotalAvailBalance,
			ProfitReal:    acc.RealizedPnl,
			ProfitUnreal:  acc.UnrealizedPnl,
			RiskRate:      0,
		}
	}
	return &exws.FutureAccount{
		FutureSubAccounts: subAccount,
	}, nil
}
func (bs *BitgetSwap) PlaceFutureOrder(currencyPair cons.CurrencyPair, contractType, price, amount string, openType, matchPrice int, leverRate float64) (string, error) {
	fOrder, err := bs.PlaceFutureOrder2(currencyPair, contractType, price, amount, openType, matchPrice, leverRate)
	return fOrder.OrderID2, err
}

/**
 * @deprecated
 * 期货下单
 * @param currencyPair   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType   合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @param price  价格
 * @param amount  委托数量
 * @param openType   1:开多   2:开空   3:平多   4:平空
 * @param matchPrice  是否为对手价 0:不是    1:是   ,当取值为1时,price无效
 */
func (bs *BitgetSwap) PlaceFutureOrder2(currencyPair cons.CurrencyPair, contractType, price, amount string, openType, matchPrice int, leverRate float64) (*exws.FutureOrder, error) {
	fOrder := &exws.FutureOrder{
		Currency:     currencyPair,
		ClientOid:    util.GenerateOrderClientId(32),
		Price:        num.ToFloat64(price),
		Amount:       num.ToFloat64(amount),
		OrderType:    openType,
		LeverRate:    leverRate,
		ContractName: contractType,
	}
	symbol := bs.adaptSymbol(currencyPair)
	uri := "/api/swap/v3/order/placeOrder"
	params := make(map[string]any)
	params["symbol"] = symbol
	params["size"] = amount
	params["client_oid"] = fOrder.ClientOid
	params["type"] = strconv.Itoa(openType)
	params["match_price"] = strconv.Itoa(matchPrice)
	params["order_type"] = "0"
	if matchPrice == 0 {
		params["price"] = price
	}
	resp, err := bs.doAuthRequest(http.MethodPost, uri, params)
	if err != nil {
		return fOrder, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return fOrder, err
	}
	orderId := num.ToInt[int](respMap["order_id"])
	if orderId <= 0 {
		return fOrder, errors.New(string(resp))
	}
	fOrder.OrderID2 = respMap["order_id"].(string)
	return fOrder, nil
}
func (bs *BitgetSwap) LimitFuturesOrder(currencyPair cons.CurrencyPair, contractType, price, amount string, openType int, opt ...cons.LimitOrderOptionalParameter) (*exws.FutureOrder, error) {
	return bs.PlaceFutureOrder2(currencyPair, contractType, price, amount, openType, 0, 10)
}
func (bs *BitgetSwap) MarketFuturesOrder(currencyPair cons.CurrencyPair, contractType, amount string, openType int) (*exws.FutureOrder, error) {
	return bs.PlaceFutureOrder2(currencyPair, contractType, "0", amount, openType, 1, 10)
}

/**
 *
 * 取消订单
 * @param symbol   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType    合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @param orderId   订单ID
 */
func (bs *BitgetSwap) FutureCancelOrder(currencyPair cons.CurrencyPair, contractType, orderId string) (bool, error) {
	uri := "/api/swap/v3/order/cancel_order"
	params := make(map[string]any)
	params["symbol"] = bs.adaptSymbol(currencyPair)
	params["orderId"] = orderId
	resp, err := bs.doAuthRequest(http.MethodPost, uri, params)
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return false, err
	}
	result := respMap["result"].(bool)
	if !result {
		return false, errors.New(respMap["err_msg"].(string))
	}
	return true, nil
}

/**
 * 用户持仓查询
 * @param symbol   btc_usd:比特币    ltc_usd :莱特币
 * @param contractType   合约类型: this_week:当周   next_week:下周   month:当月   quarter:季度
 * @return
 */
func (bs *BitgetSwap) GetFuturePosition(currencyPair cons.CurrencyPair, contractType string) ([]exws.FuturePosition, error) {
	symbol := bs.adaptSymbol(currencyPair)
	uri := "/api/swap/v3/position/singlePosition?symbol=" + symbol
	resp, err := bs.doAuthRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	type PositionRsp struct {
		Holding []struct {
			AvailPosition    float64 `json:"avail_position,string"`
			AvgCost          float64 `json:"avg_cost,string"` //开仓平均价
			Leverage         float64 `json:"leverage,string"`
			LiquidationPrice float64 `json:"liquidation_price,string"`
			Margin           string  `json:"margin"`
			Position         float64 `json:"position,string"`
			RealizedPnl      float64 `json:"realized_pnl,string"`
			Side             string  `json:"side"`
			Symbol           string  `json:"symbol"`
			Timestamp        string  `json:"timestamp"`
		} `json:"holding"`
		MarginMode string `json:"margin_mode"`
	}
	pos := PositionRsp{}
	err = json.Unmarshal(resp, &pos)
	if err != nil {
		return nil, err
	}
	if len(pos.Holding) != 2 {
		return nil, fmt.Errorf("position is not correct:%s", string(resp))
	}
	var positions []exws.FuturePosition
	p := exws.FuturePosition{
		LeverRate:      pos.Holding[0].Leverage,
		Symbol:         currencyPair,
		ForceLiquPrice: pos.Holding[0].LiquidationPrice,
	}
	for _, info := range pos.Holding {
		if info.Symbol != symbol {
			continue
		}
		if info.Side == "long" {
			p.BuyAmount = info.Position
			p.BuyAvailable = info.AvailPosition
			p.BuyPriceAvg = info.AvgCost
			p.BuyProfitReal = info.RealizedPnl
		} else {
			p.SellAmount = info.Position
			p.SellAvailable = info.AvailPosition
			p.SellPriceAvg = info.AvgCost
			p.SellProfitReal = info.RealizedPnl
		}
	}
	positions = append(positions, p)
	return positions, nil
}

/**
 *获取订单信息
 */
func (bs *BitgetSwap) GetFutureOrders(orderIds []string, currencyPair cons.CurrencyPair, contractType string) ([]exws.FutureOrder, error) {
	panic("not implement")
}

/**
 *获取单个订单信息
 */
func (bs *BitgetSwap) GetFutureOrder(orderId string, currencyPair cons.CurrencyPair, contractType string) (*exws.FutureOrder, error) {
	symbol := bs.adaptSymbol(currencyPair)
	uri := fmt.Sprintf("/api/swap/v3/order/detail?symbol=%s&orderId=%s", symbol, orderId)
	resp, err := bs.doAuthRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	result := make(map[string]any)
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}
	order := &exws.FutureOrder{}
	order.Currency = currencyPair
	order.Price = num.ToFloat64(result["price"])
	order.Amount = num.ToFloat64(result["size"])
	order.AvgPrice = num.ToFloat64(result["price_avg"])
	order.OrderID2 = orderId
	order.DealAmount = num.ToFloat64(result["filled_qty"])
	order.Fee = num.ToFloat64(result["fee"])
	order.OType = num.ToInt[int](result["type"])
	order.ClientOid, _ = result["clientOid"].(string)
	status := num.ToInt[int](result["status"])
	switch status {
	case -1:
		order.Status = cons.ORDER_CANCEL
	case 0:
		order.Status = cons.ORDER_UNFINISH
	case 1:
		order.Status = cons.ORDER_PART_FINISH
	case 2:
		order.Status = cons.ORDER_FINISH
	default:
		order.Status = cons.ORDER_UNFINISH
	}
	return order, nil
}

/**
 *获取未完成订单信息
 */
func (bs *BitgetSwap) GetUnfinishFutureOrders(currencyPair cons.CurrencyPair, contractType string) ([]exws.FutureOrder, error) {
	symbol := bs.adaptSymbol(currencyPair)
	uri := fmt.Sprintf("/api/swap/v3/order/orders?symbol=%s&from=1&to=1&limit=100&status=0", symbol)
	resp, err := bs.doAuthRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	result := make([]any, 0)
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}
	orders := make([]exws.FutureOrder, 0)
	for _, v := range result {
		vv := v.(map[string]any)
		order := exws.FutureOrder{}
		order.Currency = currencyPair
		order.Price = num.ToFloat64(vv["price"])
		order.Amount = num.ToFloat64(vv["size"])
		order.AvgPrice = num.ToFloat64(vv["price_avg"])
		order.OrderID2 = vv["order_id"].(string)
		order.DealAmount = num.ToFloat64(vv["filled_qty"])
		order.Fee = num.ToFloat64(vv["fee"])
		order.OType = num.ToInt[int](vv["type"])
		order.ClientOid = vv["client_oid"].(string)
		status := num.ToInt[int](vv["status"])
		switch status {
		case -1:
			order.Status = cons.ORDER_CANCEL
		case 0:
			order.Status = cons.ORDER_UNFINISH
		case 1:
			order.Status = cons.ORDER_PART_FINISH
		case 2:
			order.Status = cons.ORDER_FINISH
		default:
			order.Status = cons.ORDER_UNFINISH
		}
		orders = append(orders, order)
	}
	return orders, nil
}

/**
 *获取交易费
 */
func (bs *BitgetSwap) GetFee() (float64, error) {
	panic("not supported.")
}

/**
 *获取每张合约价值
 */
func (bs *BitgetSwap) GetContractValue(currencyPair cons.CurrencyPair) (float64, error) {
	panic("not supported.")
}

/**
 *获取交割时间 星期(0,1,2,3,4,5,6)，小时，分，秒
 */
func (bs *BitgetSwap) GetDeliveryTime() (int, int, int, int) {
	panic("not supported.")
}

/**
 * 获取K线数据
 */
func (bs *BitgetSwap) GetKlineRecords(contractType string, currency cons.CurrencyPair, period, size, since int) ([]exws.FutureKline, error) {
	panic("not supported.")
}
func (bs *BitgetSwap) GetServerTime() (int64, error) {
	respMap, err := web.HttpGet(bs.httpClient, fmt.Sprintf("%s/api/swap/v3/market/time", bs.baseUrl))
	if err != nil {
		return 0, err
	}
	stime := int64(num.ToInt[int](respMap["timestamp"]))
	return stime, nil
}
func (bs *BitgetSwap) adaptSymbol(pair cons.CurrencyPair) string {
	symbol := strings.ToLower(pair.ToSymbol(""))
	if pair.CurrencyB == cons.USDT {
		symbol = "cmt_" + symbol
	}
	return symbol
}

type MarginLeverage struct {
	LongLeverage        float64 `json:"long_leverage,string"`
	MarginMode          string  `json:"margin_mode"`
	ShortLeverage       float64 `json:"short_leverage,string"`
	ForwardContractFlag bool    `json:"forwardContractFlag"`
	Symbol              string  `json:"symbol"`
}

// side
// 1:多仓
// 2:空仓
func (bs *BitgetSwap) SetMarginLevel(currencyPair cons.CurrencyPair, level, side int) (*MarginLeverage, error) {
	uri := "/api/swap/v3/account/leverage"
	reqBody := make(map[string]any)
	reqBody["leverage"] = strconv.Itoa(level)
	reqBody["side"] = strconv.Itoa(side)
	reqBody["symbol"] = bs.adaptSymbol(currencyPair)
	resp, err := bs.doAuthRequest(http.MethodPost, uri, reqBody)
	if err != nil {
		return nil, err
	}
	margin := MarginLeverage{}
	err = json.Unmarshal(resp, &margin)
	if err != nil {
		return nil, err
	}
	return &margin, nil
}
func (bs *BitgetSwap) GetMarginLevel(currencyPair cons.CurrencyPair) (*MarginLeverage, error) {
	uri := "/api/swap/v3/account/settings?symbol=" + bs.adaptSymbol(currencyPair)
	resp, err := bs.doAuthRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	margin := MarginLeverage{}
	err = json.Unmarshal(resp, &margin)
	if err != nil {
		return nil, err
	}
	return &margin, nil
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

func (bs *BitgetSwap) GetContractInfo(pair cons.CurrencyPair) (*Instrument, error) {
	url := fmt.Sprintf("%s/api/swap/v3/market/contracts", bs.baseUrl)
	resp, err := web.HttpGet3(bs.httpClient, url, nil)
	if err != nil {
		return nil, err
	}
	for _, v := range resp {
		contract := v.(map[string]any)
		if contract["quote_currency"].(string) == pair.CurrencyB.String() && contract["underlying_index"].(string) == pair.CurrencyA.String() {
			return &Instrument{
				Coin:                contract["coin"].(string),
				ContractVal:         contract["contract_val"].(string),
				Delivery:            contract["delivery"].([]any),
				ForwardContractFlag: contract["forwardContractFlag"].(bool),
				PriceEndStep:        num.ToInt[int](contract["priceEndStep"]),
				QuoteCurrency:       contract["quote_currency"].(string),
				SizeIncrement:       num.ToInt[int](contract["size_increment"]),
				Symbol:              contract["contract_val"].(string),
				TickSize:            num.ToInt[int](contract["tick_size"]),
				UnderlyingIndex:     contract["underlying_index"].(string),
			}, nil
		}
	}
	return nil, errors.New("not found")
}
func (bs *BitgetSwap) GetInstruments() ([]Instrument, error) {
	url := fmt.Sprintf("%s/api/swap/v3/market/contracts", bs.baseUrl)
	resp, err := web.HttpGet3(bs.httpClient, url, nil)
	if err != nil {
		return nil, err
	}
	ins := make([]Instrument, 0)
	for _, v := range resp {
		contract := v.(map[string]any)
		ins = append(ins, Instrument{
			Coin:                contract["coin"].(string),
			ContractVal:         contract["contract_val"].(string),
			Delivery:            contract["delivery"].([]any),
			ForwardContractFlag: contract["forwardContractFlag"].(bool),
			PriceEndStep:        num.ToInt[int](contract["priceEndStep"]),
			QuoteCurrency:       contract["quote_currency"].(string),
			SizeIncrement:       num.ToInt[int](contract["size_increment"]),
			Symbol:              contract["contract_val"].(string),
			TickSize:            num.ToInt[int](contract["tick_size"]),
			UnderlyingIndex:     contract["underlying_index"].(string),
		})
	}
	return ins, nil
}

// side
// 1:多仓
// 2:空仓
// autoAppend追加保证金类型
// 0 不自动追加 1 自动追加
func (bs *BitgetSwap) ModifyAutoAppendMargin(currencyPair cons.CurrencyPair, side int, autoAppend int) (bool, error) {
	uri := "/api/swap/v3/account/modifyAutoAppendMargin"
	reqBody := make(map[string]any)
	reqBody["append_type"] = autoAppend
	reqBody["side"] = side
	reqBody["symbol"] = bs.adaptSymbol(currencyPair)
	resp, err := bs.doAuthRequest(http.MethodPost, uri, reqBody)
	if err != nil {
		return false, err
	}
	autoMargin := make(map[string]any)
	err = json.Unmarshal(resp, &autoMargin)
	if err != nil {
		return false, err
	}
	if !autoMargin["result"].(bool) {
		return false, nil
	}
	return true, nil
}
