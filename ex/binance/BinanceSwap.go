package binance

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/conbanwa/num"
	"github.com/conbanwa/slice"
	"net/url"
	. "qa3/wstrader"
	. "qa3/wstrader/cons"
	. "qa3/wstrader/q"
	. "qa3/wstrader/util"
	. "qa3/wstrader/web"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	baseUrl = "https://fapi.binance.com"
)

type Swap struct {
	Binance
	f *Futures
}

func NewBinanceSwap(config *APIConfig) *Swap {
	if config.Endpoint == "" {
		config.Endpoint = baseUrl
	}
	bs := &Swap{
		Binance: Binance{
			baseUrl:    config.Endpoint,
			accessKey:  config.ApiKey,
			apiV1:      config.Endpoint + "/fapi/v1/",
			secretKey:  config.ApiSecretKey,
			httpClient: config.HttpClient,
		},
		f: NewBinanceFutures(&APIConfig{
			Endpoint:     strings.ReplaceAll(config.Endpoint, "fapi", "dapi"),
			HttpClient:   config.HttpClient,
			ApiKey:       config.ApiKey,
			ApiSecretKey: config.ApiSecretKey,
			Lever:        config.Lever,
		}),
	}
	bs.setTimeOffset()
	return bs
}
func (bs *Swap) SetBaseUri(uri string) {
	bs.baseUrl = uri
}
func (bs *Swap) String() string {
	return BINANCE_SWAP
}
func (bs *Swap) Ping() bool {
	if _, err := HttpGet(bs.httpClient, bs.apiV1+"ping"); err != nil {
		return false
	}
	return true
}
func (bs *Swap) setTimeOffset() error {
	respMap, err := HttpGet(bs.httpClient, bs.apiV1+ServerTimeUrl)
	if err != nil {
		return err
	}
	stime := int64(num.ToInt[int](respMap["serverTime"]))
	st := time.Unix(stime/1000, 1000000*(stime%1000))
	lt := time.Now()
	offset := st.Sub(lt).Nanoseconds()
	bs.timeOffset = offset
	return nil
}
func (bs *Swap) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	panic("not supported.")
}
func (bs *Swap) GetFutureTicker(currency CurrencyPair, contractType string) (*Ticker, error) {
	if contractType == SWAP_CONTRACT {
		return bs.f.GetFutureTicker(currency.AdaptUsdtToUsd(), SWAP_CONTRACT)
	}
	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}
	currency2 := bs.adaptCurrencyPair(currency)
	tickerPriceUri := bs.apiV1 + "ticker/price?symbol=" + currency2.ToSymbol("")
	tickerBookUri := bs.apiV1 + "ticker/bookTicker?symbol=" + currency2.ToSymbol("")
	tickerPriceMap := make(map[string]any)
	tickerBookMap := make(map[string]any)
	var err1 error
	var err2 error
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		tickerPriceMap, err1 = HttpGet(bs.httpClient, tickerPriceUri)
	}()
	go func() {
		defer wg.Done()
		tickerBookMap, err2 = HttpGet(bs.httpClient, tickerBookUri)
	}()
	wg.Wait()
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	var ticker Ticker
	ticker.Pair = currency
	ticker.Date = uint64(time.Now().UnixNano() / int64(time.Millisecond))
	ticker.Last = num.ToFloat64(tickerPriceMap["price"])
	ticker.Buy = num.ToFloat64(tickerBookMap["bidPrice"])
	ticker.Sell = num.ToFloat64(tickerBookMap["askPrice"])
	return &ticker, nil
}
func (bs *Swap) GetFutureDepth(currency CurrencyPair, contractType string, size int) (*Depth, error) {
	if contractType == SWAP_CONTRACT {
		return bs.f.GetFutureDepth(currency.AdaptUsdtToUsd(), SWAP_CONTRACT, size)
	}
	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}
	if size <= 5 {
		size = 5
	} else if size <= 10 {
		size = 10
	} else if size <= 20 {
		size = 20
	} else if size <= 50 {
		size = 50
	} else if size <= 100 {
		size = 100
	} else if size <= 500 {
		size = 500
	} else {
		size = 1000
	}
	currencyPair2 := bs.adaptCurrencyPair(currency)
	apiUrl := fmt.Sprintf(bs.apiV1+DepthUri, currencyPair2.ToSymbol(""), size)
	resp, err := HttpGet(bs.httpClient, apiUrl)
	if err != nil {
		return nil, err
	}
	if _, isok := resp["code"]; isok {
		return nil, errors.New(resp["msg"].(string))
	}
	bids := resp["bids"].([]any)
	asks := resp["asks"].([]any)
	depth := new(Depth)
	depth.Pair = currency
	depth.UTime = time.Now()
	n := 0
	for _, bid := range bids {
		_bid := bid.([]any)
		amount := num.ToFloat64(_bid[1])
		price := num.ToFloat64(_bid[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
		n++
		if n == size {
			break
		}
	}
	n = 0
	for _, ask := range asks {
		_ask := ask.([]any)
		amount := num.ToFloat64(_ask[1])
		price := num.ToFloat64(_ask[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
		n++
		if n == size {
			break
		}
	}
	return depth, nil
}
func (bs *Swap) GetFutureOrderHistory(pair CurrencyPair, contractType string, optional ...OptionalParameter) ([]FutureOrder, error) {
	panic("implement me")
}
func (bs *Swap) GetTrades(contractType string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	if contractType == SWAP_CONTRACT {
		return bs.f.GetTrades(SWAP_CONTRACT, currencyPair.AdaptUsdtToUsd(), since)
	}
	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}
	param := url.Values{}
	param.Set("symbol", bs.adaptCurrencyPair(currencyPair).ToSymbol(""))
	param.Set("limit", "500")
	if since > 0 {
		param.Set("fromId", strconv.Itoa(int(since)))
	}
	apiUrl := bs.apiV1 + "historicalTrades?" + param.Encode()
	resp, err := HttpGet3(bs.httpClient, apiUrl, map[string]string{
		"X-MBX-APIKEY": bs.accessKey})
	if err != nil {
		return nil, err
	}
	var trades []Trade
	for _, v := range resp {
		m := v.(map[string]any)
		ty := SELL
		if m["isBuyerMaker"].(bool) {
			ty = BUY
		}
		trades = append(trades, Trade{
			Tid:    num.ToInt[int64](m["id"]),
			Type:   ty,
			Amount: num.ToFloat64(m["qty"]),
			Price:  num.ToFloat64(m["price"]),
			Date:   num.ToInt[int64](m["time"]),
			Pair:   currencyPair,
		})
	}
	return trades, nil
}
func (bs *Swap) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	respMap, err := HttpGet(bs.httpClient, bs.apiV1+"premiumIndex?symbol="+bs.adaptCurrencyPair(currencyPair).ToSymbol(""))
	if err != nil {
		return 0.0, err
	}
	return num.ToFloat64(respMap["markPrice"]), nil
}
func (bs *Swap) GetFutureUserinfo(currencyPair ...CurrencyPair) (*FutureAccount, error) {
	acc, err := bs.f.GetFutureUserinfo(currencyPair...)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	bs.buildParamsSigned(&params)
	path := bs.apiV1 + AccountUri + params.Encode()
	respMap, err := HttpGet2(bs.httpClient, path, map[string]string{"X-MBX-APIKEY": bs.accessKey})
	if err != nil {
		return nil, err
	}
	if _, isok := respMap["code"]; isok == true {
		return nil, errors.New(respMap["msg"].(string))
	}
	balances := respMap["assets"].([]any)
	for _, v := range balances {
		vv := v.(map[string]any)
		currency := NewCurrency(vv["asset"].(string), "").AdaptBccToBch()
		acc.FutureSubAccounts[currency] = FutureSubAccount{
			Currency:      currency,
			AccountRights: num.ToFloat64(vv["marginBalance"]),
			KeepDeposit:   num.ToFloat64(vv["maintMargin"]),
			ProfitUnreal:  num.ToFloat64(vv["unrealizedProfit"]),
		}
	}
	return acc, nil
}

// @deprecated please call the Wallet api
func (bs *Swap) Transfer(currency Currency, transferType int, amount float64) (int64, error) {
	params := url.Values{}
	params.Set("currency", currency.String())
	params.Set("amount", fmt.Sprint(amount))
	params.Set("type", strconv.Itoa(transferType))
	uri := GlobalApiBaseUrl + "/sapi/v1/futures/transfer"
	bs.buildParamsSigned(&params)
	resp, err := HttpPostForm2(bs.httpClient, uri, params,
		map[string]string{"X-MBX-APIKEY": bs.accessKey})
	if err != nil {
		return 0, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return 0, err
	}
	return num.ToInt[int64](respMap["tranId"]), nil
}
func (bs *Swap) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice int, leverRate float64) (string, error) {
	fOrder, err := bs.PlaceFutureOrder2(currencyPair, contractType, price, amount, openType, matchPrice)
	return fOrder.OrderID2, err
}
func (bs *Swap) PlaceFutureOrder2(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice int, opt ...LimitOrderOptionalParameter) (*FutureOrder, error) {
	if contractType == SWAP_CONTRACT {
		orderId, err := bs.f.PlaceFutureOrder2(currencyPair.AdaptUsdtToUsd(), contractType, price, amount, openType, matchPrice, opt...)
		return &FutureOrder{
			OrderID2:     orderId,
			Price:        num.ToFloat64(price),
			Amount:       num.ToFloat64(amount),
			Status:       ORDER_UNFINISH,
			Currency:     currencyPair,
			OType:        openType,
			LeverRate:    0,
			ContractName: contractType,
		}, err
	}
	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}
	fOrder := &FutureOrder{
		Currency:     currencyPair,
		ClientOid:    GenerateOrderClientId(32),
		Price:        num.ToFloat64(price),
		Amount:       num.ToFloat64(amount),
		OrderType:    openType,
		LeverRate:    0,
		ContractName: contractType,
	}
	pair := bs.adaptCurrencyPair(currencyPair)
	path := bs.apiV1 + OrderUri
	params := url.Values{}
	params.Set("symbol", pair.ToSymbol(""))
	params.Set("quantity", amount)
	params.Set("newClientOrderId", fOrder.ClientOid)
	switch openType {
	case OPEN_BUY, CLOSE_SELL:
		params.Set("side", "BUY")
		if len(opt) > 0 && opt[0] == FuturesTwoWayPositionMode {
			params.Set("positionSide", "LONG")
		}
	case OPEN_SELL, CLOSE_BUY:
		params.Set("side", "SELL")
		if len(opt) > 0 && opt[0] == FuturesTwoWayPositionMode {
			params.Set("positionSide", "SHORT")
		}
	}
	if matchPrice == 0 {
		params.Set("type", "LIMIT")
		params.Set("price", price)
		params.Set("timeInForce", "GTC")
	} else {
		params.Set("type", "MARKET")
	}
	bs.buildParamsSigned(&params)
	resp, err := HttpPostForm2(bs.httpClient, path, params,
		map[string]string{"X-MBX-APIKEY": bs.accessKey})
	if err != nil {
		return fOrder, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return fOrder, err
	}
	orderId := num.ToInt[int](respMap["orderId"])
	if orderId <= 0 {
		return fOrder, errors.New(slice.Bytes2String(resp))
	}
	fOrder.OrderID2 = strconv.Itoa(orderId)
	return fOrder, nil
}
func (bs *Swap) LimitFuturesOrder(currencyPair CurrencyPair, contractType, price, amount string, openType int, opt ...LimitOrderOptionalParameter) (*FutureOrder, error) {
	return bs.PlaceFutureOrder2(currencyPair, contractType, price, amount, openType, 0, opt...)
}
func (bs *Swap) MarketFuturesOrder(currencyPair CurrencyPair, contractType, amount string, openType int) (*FutureOrder, error) {
	return bs.PlaceFutureOrder2(currencyPair, contractType, "0", amount, openType, 1)
}
func (bs *Swap) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	if contractType == SWAP_CONTRACT {
		return bs.f.FutureCancelOrder(currencyPair.AdaptUsdtToUsd(), contractType, orderId)
	}
	if contractType != SWAP_USDT_CONTRACT {
		return false, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}
	currencyPair = bs.adaptCurrencyPair(currencyPair)
	path := bs.apiV1 + OrderUri
	params := url.Values{}
	params.Set("symbol", bs.adaptCurrencyPair(currencyPair).ToSymbol(""))
	if strings.HasPrefix(orderId, "goex") { //goex default clientOrderId Features
		params.Set("origClientOrderId", orderId)
	} else {
		params.Set("orderId", orderId)
	}
	bs.buildParamsSigned(&params)
	resp, err := HttpDeleteForm(bs.httpClient, path, params, map[string]string{"X-MBX-APIKEY": bs.accessKey})
	if err != nil {
		return false, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return false, err
	}
	orderIdCanceled := num.ToInt[int](respMap["orderId"])
	if orderIdCanceled <= 0 {
		return false, errors.New(slice.Bytes2String(resp))
	}
	return true, nil
}
func (bs *Swap) FutureCancelAllOrders(currencyPair CurrencyPair, contractType string) (bool, error) {
	if contractType == SWAP_CONTRACT {
		return false, errors.New("not support")
	}
	if contractType == SWAP_CONTRACT {
		return false, errors.New("not support")
	}
	if contractType != SWAP_USDT_CONTRACT {
		return false, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}
	currencyPair = bs.adaptCurrencyPair(currencyPair)
	path := bs.apiV1 + "allOpenOrders"
	params := url.Values{}
	params.Set("symbol", bs.adaptCurrencyPair(currencyPair).ToSymbol(""))
	bs.buildParamsSigned(&params)
	resp, err := HttpDeleteForm(bs.httpClient, path, params, map[string]string{"X-MBX-APIKEY": bs.accessKey})
	if err != nil {
		return false, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return false, err
	}
	if num.ToInt[int](respMap["code"]) != 200 {
		return false, errors.New(respMap["msg"].(string))
	}
	return true, nil
}
func (bs *Swap) FutureCancelOrders(currencyPair CurrencyPair, contractType string, orderIdList []string) (bool, error) {
	if contractType != SWAP_USDT_CONTRACT {
		return false, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}
	currencyPair = bs.adaptCurrencyPair(currencyPair)
	path := bs.apiV1 + "batchOrders"
	if len(orderIdList) == 0 {
		return false, errors.New("list is empty, no order will be cancel")
	}
	list, _ := json.Marshal(orderIdList)
	params := url.Values{}
	params.Set("symbol", bs.adaptCurrencyPair(currencyPair).ToSymbol(""))
	params.Set("orderIdList", slice.Bytes2String(list))
	bs.buildParamsSigned(&params)
	resp, err := HttpDeleteForm(bs.httpClient, path, params, map[string]string{"X-MBX-APIKEY": bs.accessKey})
	if err != nil {
		return false, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return false, err
	}
	if num.ToInt[int](respMap["code"]) != 200 {
		return false, errors.New(respMap["msg"].(string))
	}
	return true, nil
}
func (bs *Swap) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	if contractType == SWAP_CONTRACT {
		return bs.f.GetFuturePosition(currencyPair.AdaptUsdtToUsd(), contractType)
	}
	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}
	currencyPair1 := bs.adaptCurrencyPair(currencyPair)
	params := url.Values{}
	bs.buildParamsSigned(&params)
	path := bs.apiV1 + "positionRisk?" + params.Encode()
	result, err := HttpGet3(bs.httpClient, path, map[string]string{"X-MBX-APIKEY": bs.accessKey})
	if err != nil {
		return nil, err
	}
	var positions []FuturePosition
	for _, info := range result {
		cont := info.(map[string]any)
		if cont["symbol"] != currencyPair1.ToSymbol("") {
			continue
		}
		p := FuturePosition{
			LeverRate:      num.ToFloat64(cont["leverage"]),
			Symbol:         currencyPair,
			ForceLiquPrice: num.ToFloat64(cont["liquidationPrice"]),
		}
		amount := num.ToFloat64(cont["positionAmt"])
		price := num.ToFloat64(cont["entryPrice"])
		if upnl := num.ToFloat64(cont["unRealizedProfit"]); amount > 0 {
			p.BuyAmount = amount
			p.BuyPriceAvg = price
			p.BuyProfitReal = upnl
		} else if amount < 0 {
			p.SellAmount = amount
			p.SellPriceAvg = price
			p.SellProfitReal = upnl
		}
		positions = append(positions, p)
	}
	return positions, nil
}
func (bs *Swap) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	if contractType == SWAP_CONTRACT {
		return nil, errors.New("not support")
	}
	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}
	if len(orderIds) == 0 {
		return nil, errors.New("orderIds is empty")
	}
	currencyPair1 := bs.adaptCurrencyPair(currencyPair)
	params := url.Values{}
	params.Set("symbol", currencyPair1.ToSymbol(""))
	bs.buildParamsSigned(&params)
	path := bs.apiV1 + "allOrders?" + params.Encode()
	result, err := HttpGet3(bs.httpClient, path, map[string]string{"X-MBX-APIKEY": bs.accessKey})
	if err != nil {
		return nil, err
	}
	orders := make([]FutureOrder, 0)
	for _, info := range result {
		_ord := info.(map[string]any)
		if _ord["symbol"].(string) != currencyPair1.ToSymbol("") {
			continue
		}
		orderId := num.ToInt[int](_ord["orderId"])
		ordId := strconv.Itoa(orderId)
		for _, id := range orderIds {
			if id == ordId {
				order := bs.parseOrder(_ord)
				order.Currency = currencyPair
				orders = append(orders, *order)
				break
			}
		}
	}
	return orders, nil
}
func (bs *Swap) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	if contractType == SWAP_CONTRACT {
		return bs.f.GetFutureOrder(orderId, currencyPair.AdaptUsdtToUsd(), contractType)
	}
	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}
	currencyPair1 := bs.adaptCurrencyPair(currencyPair)
	params := url.Values{}
	params.Set("symbol", currencyPair1.ToSymbol(""))
	params.Set("orderId", orderId)
	bs.buildParamsSigned(&params)
	path := bs.apiV1 + "allOrders?" + params.Encode()
	result, err := HttpGet3(bs.httpClient, path, map[string]string{"X-MBX-APIKEY": bs.accessKey})
	if err != nil {
		return nil, err
	}
	ordId, _ := strconv.Atoi(orderId)
	for _, info := range result {
		_ord := info.(map[string]any)
		if _ord["symbol"].(string) != currencyPair1.ToSymbol("") {
			continue
		}
		if num.ToInt[int](_ord["orderId"]) != ordId {
			continue
		}
		order := bs.parseOrder(_ord)
		order.Currency = currencyPair
		return order, nil
	}
	return nil, errors.New(fmt.Sprintf("not found order:%s", orderId))
}
func (bs *Swap) parseOrder(rsp map[string]any) *FutureOrder {
	order := &FutureOrder{}
	order.Price = num.ToFloat64(rsp["price"])
	order.Amount = num.ToFloat64(rsp["origQty"])
	order.DealAmount = num.ToFloat64(rsp["executedQty"])
	order.AvgPrice = num.ToFloat64(rsp["avgPrice"])
	order.OrderTime = num.ToInt[int64](rsp["time"])
	status := rsp["status"].(string)
	order.Status = bs.parseOrderStatus(status)
	order.OrderID = num.ToInt[int64](rsp["orderId"])
	order.OrderID2 = strconv.Itoa(int(order.OrderID))
	order.OType = OPEN_BUY
	if rsp["side"].(string) == "SELL" {
		order.OType = OPEN_SELL
	}
	//GTC - Good Till Cancel 成交为止
	//IOC - Immediate or Cancel 无法立即成交(吃单)的部分就撤销
	//FOK - Fill or Kill 无法全部立即成交就撤销
	//GTX - Good Till Crossing 无法成为挂单方就撤销
	ot := rsp["timeInForce"].(string)
	switch ot {
	case "GTC":
		order.OrderType = ORDER_FEATURE_LIMIT
	case "IOC":
		order.OrderType = ORDER_FEATURE_IOC
	case "FOK":
		order.OrderType = ORDER_FEATURE_FOK
	case "GTX":
		order.OrderType = ORDER_FEATURE_IOC
	}
	//LIMIT 限价单
	//MARKET 市价单
	//STOP 止损限价单
	//STOP_MARKET 止损市价单
	//TAKE_RPOFIT 止盈限价单
	//TAKE_RPOFIT_MARKET 止盈市价单
	return order
}
func (bs *Swap) parseOrderStatus(sts string) TradeStatus {
	orderStatus := ORDER_UNFINISH
	switch sts {
	case "PARTIALLY_FILLED", "partially_filled":
		orderStatus = ORDER_PART_FINISH
	case "FILLED", "filled":
		orderStatus = ORDER_FINISH
	case "CANCELED", "REJECTED", "EXPIRED":
		orderStatus = ORDER_CANCEL
	}
	return orderStatus
}
func (bs *Swap) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	if contractType == SWAP_CONTRACT {
		return bs.f.GetUnfinishFutureOrders(currencyPair.AdaptUsdtToUsd(), contractType)
	}
	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}
	currencyPair1 := bs.adaptCurrencyPair(currencyPair)
	params := url.Values{}
	params.Set("symbol", currencyPair1.ToSymbol(""))
	bs.buildParamsSigned(&params)
	path := bs.apiV1 + "openOrders?" + params.Encode()
	result, err := HttpGet3(bs.httpClient, path, map[string]string{"X-MBX-APIKEY": bs.accessKey})
	if err != nil {
		return nil, err
	}
	orders := make([]FutureOrder, 0)
	for _, info := range result {
		_ord := info.(map[string]any)
		if _ord["symbol"].(string) != currencyPair1.ToSymbol("") {
			continue
		}
		order := bs.parseOrder(_ord)
		order.Currency = currencyPair
		orders = append(orders, *order)
	}
	return orders, nil
}
func (bs *Swap) GetFee() (float64, error) {
	panic("not supported.")
}
func (bs *Swap) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	panic("not supported.")
}
func (bs *Swap) GetDeliveryTime() (int, int, int, int) {
	panic("not supported.")
}
func (bs *Swap) GetKlineRecords(contractType string, currency CurrencyPair, period KlinePeriod, size int, opt ...OptionalParameter) ([]FutureKline, error) {
	if contractType == SWAP_CONTRACT {
		return bs.f.GetKlineRecords(contractType, currency.AdaptUsdtToUsd(), period, size, opt...)
	}
	if contractType != SWAP_USDT_CONTRACT {
		return nil, errors.New("contract is error,please incoming SWAP_CONTRACT or SWAP_USDT_CONTRACT")
	}
	currency2 := bs.adaptCurrencyPair(currency)
	params := url.Values{}
	params.Set("symbol", currency2.ToSymbol(""))
	params.Set("interval", internalKlinePeriodConverter[period])
	//params.Set("endTime", strconv.Itoa(int(time.Now().UnixNano()/1000000)))
	params.Set("limit", strconv.Itoa(size))
	MergeOptionalParameter(&params, opt...)
	klineUrl := bs.apiV1 + KlineUri + "?" + params.Encode()
	klines, err := HttpGet3(bs.httpClient, klineUrl, nil)
	if err != nil {
		return nil, err
	}
	var klineRecords []FutureKline
	for _, _record := range klines {
		r := Kline{Pair: currency}
		record := _record.([]any)
		r.Timestamp = int64(record[0].(float64)) / 1000 //to unix timestramp
		r.Open = num.ToFloat64(record[1])
		r.High = num.ToFloat64(record[2])
		r.Low = num.ToFloat64(record[3])
		r.Close = num.ToFloat64(record[4])
		r.Vol = num.ToFloat64(record[5])
		klineRecords = append(klineRecords, FutureKline{Kline: &r})
	}
	return klineRecords, nil
}
func (bs *Swap) GetServerTime() (int64, error) {
	respMap, err := HttpGet(bs.httpClient, bs.apiV1+ServerTimeUrl)
	if err != nil {
		return 0, err
	}
	stime := int64(num.ToInt[int](respMap["serverTime"]))
	return stime, nil
}
func (bs *Swap) adaptCurrencyPair(pair CurrencyPair) CurrencyPair {
	return pair.AdaptUsdToUsdt()
}
