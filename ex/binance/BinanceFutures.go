package binance

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/conbanwa/num"
	"github.com/conbanwa/slice"
	. "github.com/conbanwa/wstrader"
	. "github.com/conbanwa/wstrader/cons"
	. "github.com/conbanwa/wstrader/q"
	. "github.com/conbanwa/wstrader/util"
	. "github.com/conbanwa/wstrader/web"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
type AccountResponse struct {
	FeeTier  int  `json:"feeTier"`
	CanTrade bool `json:"canTrade"`
	Assets   []struct {
		Asset            string  `json:"asset"`
		WalletBalance    float64 `json:"walletBalance,string"`
		MarginBalance    float64 `json:"marginBalance,string"`
		UnrealizedProfit float64 `json:"unrealizedProfit,string"`
		MaintMargin      float64 `json:"maintMargin,string"`
	} `json:"assets"`
}
type OrderInfoResponse struct {
	BaseResponse
	Symbol        string  `json:"symbol"`
	Pair          string  `json:"pair"`
	ClientOrderId string  `json:"clientOrderId"`
	OrderId       int64   `json:"orderId"`
	AvgPrice      float64 `json:"avgPrice,string"`
	ExecutedQty   float64 `json:"executedQty,string"`
	OrigQty       float64 `json:"origQty,string"`
	Price         float64 `json:"price,string"`
	Side          string  `json:"side"`
	PositionSide  string  `json:"positionSide"`
	Status        string  `json:"status"`
	Type          string  `json:"type"`
	Time          int64   `json:"time"`
	UpdateTime    int64   `json:"updateTime"`
}
type PositionRiskResponse struct {
	Symbol           string  `json:"symbol"`
	PositionAmt      float64 `json:"positionAmt,string"`
	EntryPrice       float64 `json:"entryPrice,string"`
	UnRealizedProfit float64 `json:"unRealizedProfit,string"`
	LiquidationPrice float64 `json:"liquidationPrice,string"`
	Leverage         float64 `json:"leverage,string"`
	MarginType       string  `json:"marginType"`
	PositionSide     string  `json:"positionSide"`
}
type SymbolInfo struct {
	Symbol         string
	Pair           string
	ContractType   string `json:"contractType"`
	DeliveryDate   int64  `json:"deliveryDate"`
	ContractStatus string `json:"contractStatus"`
	ContractSize   int    `json:"contractSize"`
	PricePrecision int    `json:"pricePrecision"`
}
type Futures struct {
	base         *Binance
	apikey       string
	exchangeInfo *struct {
		Symbols []SymbolInfo `json:"symbols"`
	}
}

func NewBinanceFutures(config *APIConfig) *Futures {
	if config.Endpoint == "" {
		config.Endpoint = "https://dapi.binance.com"
	}
	if config.HttpClient == nil {
		config.HttpClient = http.DefaultClient
	}
	bs := &Futures{
		apikey: config.ApiKey,
		base:   NewWithConfig(config),
	}
	bs.base.apiV1 = config.Endpoint + "/dapi/v1/"
	go bs.GetExchangeInfo()
	return bs
}
func (bs *Futures) SetBaseUri(uri string) {
	bs.base.baseUrl = uri
}
func (bs *Futures) String() string {
	return BINANCE_FUTURES
}
func (bs *Futures) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	symbol, err := bs.adaptToSymbol(currencyPair, contractType)
	if err != nil {
		return nil, err
	}
	ticker24hrUri := bs.base.apiV1 + "ticker/24hr?symbol=" + symbol
	tickerBookUri := bs.base.apiV1 + "ticker/bookTicker?symbol=" + symbol
	var (
		ticker24HrResp []any
		tickerBookResp []any
		err1           error
		err2           error
		wg             = sync.WaitGroup{}
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		ticker24HrResp, err1 = HttpGet3(bs.base.httpClient, ticker24hrUri, map[string]string{})
	}()
	go func() {
		defer wg.Done()
		tickerBookResp, err2 = HttpGet3(bs.base.httpClient, tickerBookUri, map[string]string{})
	}()
	wg.Wait()
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	if len(ticker24HrResp) == 0 {
		return nil, errors.New("response is empty")
	}
	if len(tickerBookResp) == 0 {
		return nil, errors.New("response is empty")
	}
	ticker24HrMap := ticker24HrResp[0].(map[string]any)
	tickerBookMap := tickerBookResp[0].(map[string]any)
	var ticker Ticker
	ticker.Pair = currencyPair
	ticker.Date = num.ToInt[uint64](tickerBookMap["time"])
	ticker.Last = num.ToFloat64(ticker24HrMap["lastPrice"])
	ticker.Buy = num.ToFloat64(tickerBookMap["bidPrice"])
	ticker.Sell = num.ToFloat64(tickerBookMap["askPrice"])
	ticker.High = num.ToFloat64(ticker24HrMap["highPrice"])
	ticker.Low = num.ToFloat64(ticker24HrMap["lowPrice"])
	ticker.Vol = num.ToFloat64(ticker24HrMap["volume"])
	return &ticker, nil
}
func (bs *Futures) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	symbol, err := bs.adaptToSymbol(currencyPair, contractType)
	if err != nil {
		return nil, err
	}
	limit := 5
	if size <= 5 {
		limit = 5
	} else if size <= 10 {
		limit = 10
	} else if size <= 20 {
		limit = 20
	} else if size <= 50 {
		limit = 50
	} else if size <= 100 {
		limit = 100
	} else if size <= 500 {
		limit = 500
	} else {
		limit = 1000
	}
	depthUri := bs.base.apiV1 + "depth?symbol=%s&limit=%d"
	ret, err := HttpGet(bs.base.httpClient, fmt.Sprintf(depthUri, symbol, limit))
	if err != nil {
		return nil, err
	}
	// logs.D(ret)
	var dep Depth
	dep.ContractType = contractType
	dep.Pair = currencyPair
	eT := int64(ret["E"].(float64))
	dep.UTime = time.Unix(0, eT*int64(time.Millisecond))
	for _, item := range ret["asks"].([]any) {
		ask := item.([]any)
		dep.AskList = append(dep.AskList, DepthRecord{
			Price:  num.ToFloat64(ask[0]),
			Amount: num.ToFloat64(ask[1]),
		})
	}
	for _, item := range ret["bids"].([]any) {
		bid := item.([]any)
		dep.BidList = append(dep.BidList, DepthRecord{
			Price:  num.ToFloat64(bid[0]),
			Amount: num.ToFloat64(bid[1]),
		})
	}
	sort.Sort(sort.Reverse(dep.AskList))
	return &dep, nil
}
func (bs *Futures) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	panic("not supported.")
}
func (bs *Futures) GetFutureUserinfo(currencyPair ...CurrencyPair) (*FutureAccount, error) {
	accountUri := bs.base.apiV1 + "account"
	param := url.Values{}
	bs.base.buildParamsSigned(&param)
	respData, err := HttpGet5(bs.base.httpClient, accountUri+"?"+param.Encode(), map[string]string{
		"X-MBX-APIKEY": bs.apikey})
	if err != nil {
		return nil, err
	}
	// logs.D(slice.Bytes2String(respData))
	var (
		accountResp    AccountResponse
		futureAccounts FutureAccount
	)
	err = json.Unmarshal(respData, &accountResp)
	if err != nil {
		return nil, fmt.Errorf("response body: %s , %w", slice.Bytes2String(respData), err)
	}
	futureAccounts.FutureSubAccounts = make(map[Currency]FutureSubAccount, 4)
	for _, asset := range accountResp.Assets {
		currency := NewCurrency(asset.Asset, "")
		futureAccounts.FutureSubAccounts[currency] = FutureSubAccount{
			Currency:      NewCurrency(asset.Asset, ""),
			AccountRights: asset.MarginBalance,
			KeepDeposit:   asset.MaintMargin,
			ProfitReal:    0,
			ProfitUnreal:  asset.UnrealizedProfit,
			RiskRate:      0,
		}
	}
	return &futureAccounts, nil
}
func (bs *Futures) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice int, leverRate float64) (string, error) {
	return bs.PlaceFutureOrder2(currencyPair, contractType, price, amount, openType, matchPrice)
}
func (bs *Futures) PlaceFutureOrder2(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice int, opt ...LimitOrderOptionalParameter) (string, error) {
	apiPath := "order"
	symbol, err := bs.adaptToSymbol(currencyPair, contractType)
	if err != nil {
		return "", err
	}
	param := url.Values{}
	param.Set("symbol", symbol)
	param.Set("newClientOrderId", GenerateOrderClientId(32))
	param.Set("quantity", amount)
	param.Set("newOrderRespType", "ACK")
	if matchPrice == 0 {
		param.Set("type", "LIMIT")
		param.Set("timeInForce", "GTC")
		param.Set("price", price)
	} else {
		param.Set("type", "MARKET")
	}
	switch openType {
	case OPEN_BUY, CLOSE_SELL:
		param.Set("side", "BUY")
		if len(opt) > 0 && opt[0] == FuturesTwoWayPositionMode {
			param.Set("positionSide", "LONG")
		}
	case OPEN_SELL, CLOSE_BUY:
		param.Set("side", "SELL")
		if len(opt) > 0 && opt[0] == FuturesTwoWayPositionMode {
			param.Set("positionSide", "SHORT")
		}
	}
	bs.base.buildParamsSigned(&param)
	resp, err := HttpPostForm2(bs.base.httpClient, fmt.Sprintf("%s%s", bs.base.apiV1, apiPath), param,
		map[string]string{"X-MBX-APIKEY": bs.apikey})
	if err != nil {
		return "", err
	}
	// logs.D(slice.Bytes2String(resp))
	var response struct {
		BaseResponse
		OrderId int64 `json:"orderId"`
	}
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return "", err
	}
	if response.Code == 0 {
		return fmt.Sprint(response.OrderId), nil
	}
	return "", errors.New(response.Msg)
}
func (bs *Futures) LimitFuturesOrder(currencyPair CurrencyPair, contractType, price, amount string, openType int, opt ...LimitOrderOptionalParameter) (*FutureOrder, error) {
	orderId, err := bs.PlaceFutureOrder2(currencyPair, contractType, price, amount, openType, 0, opt...)
	return &FutureOrder{
		OrderID2:     orderId,
		Currency:     currencyPair,
		ContractName: contractType,
		Amount:       num.ToFloat64(amount),
		Price:        num.ToFloat64(price),
		OType:        openType,
	}, err
}
func (bs *Futures) MarketFuturesOrder(currencyPair CurrencyPair, contractType, amount string, openType int) (*FutureOrder, error) {
	orderId, err := bs.PlaceFutureOrder2(currencyPair, contractType, "", amount, openType, 1)
	return &FutureOrder{
		OrderID2:     orderId,
		Currency:     currencyPair,
		ContractName: contractType,
		Amount:       num.ToFloat64(amount),
		OType:        openType,
	}, err
}
func (bs *Futures) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	apiPath := "order"
	symbol, err := bs.adaptToSymbol(currencyPair, contractType)
	if err != nil {
		return false, err
	}
	param := url.Values{}
	param.Set("symbol", symbol)
	if strings.HasPrefix(orderId, "goex") {
		param.Set("origClientOrderId", orderId)
	} else {
		param.Set("orderId", orderId)
	}
	bs.base.buildParamsSigned(&param)
	reqUrl := fmt.Sprintf("%s%s?%s", bs.base.apiV1, apiPath, param.Encode())
	resp, err := HttpDeleteForm(bs.base.httpClient, reqUrl, url.Values{}, map[string]string{"X-MBX-APIKEY": bs.apikey})
	if err != nil {
		log.Error().Bytes("resp", resp).Str("request url", reqUrl).Send()
		return false, err
	}
	// logs.D(slice.Bytes2String(resp))
	return true, nil
}
func (bs *Futures) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	symbol, err := bs.adaptToSymbol(currencyPair, contractType)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	bs.base.buildParamsSigned(&params)
	path := bs.base.apiV1 + "positionRisk?" + params.Encode()
	respBody, err := HttpGet5(bs.base.httpClient, path, map[string]string{"X-MBX-APIKEY": bs.apikey})
	if err != nil {
		return nil, err
	}
	// logs.D(slice.Bytes2String(respBody))
	var (
		positionRiskResponse []PositionRiskResponse
		positions            []FuturePosition
	)
	err = json.Unmarshal(respBody, &positionRiskResponse)
	if err != nil {
		log.Error().Err(err).Bytes("response data", respBody).Send()
		return nil, err
	}
	for _, info := range positionRiskResponse {
		if info.Symbol != symbol {
			continue
		}
		p := FuturePosition{
			LeverRate:      info.Leverage,
			Symbol:         currencyPair,
			ForceLiquPrice: info.LiquidationPrice,
		}
		if info.PositionAmt > 0 {
			p.BuyAmount = info.PositionAmt
			p.BuyAvailable = info.PositionAmt
			p.BuyPriceAvg = info.EntryPrice
			p.BuyPriceCost = info.EntryPrice
			p.BuyProfit = info.UnRealizedProfit
			p.BuyProfitReal = info.UnRealizedProfit
		} else if info.PositionAmt < 0 {
			p.SellAmount = math.Abs(info.PositionAmt)
			p.SellAvailable = math.Abs(info.PositionAmt)
			p.SellPriceAvg = info.EntryPrice
			p.SellPriceCost = info.EntryPrice
			p.SellProfit = info.UnRealizedProfit
			p.SellProfitReal = info.UnRealizedProfit
		}
		positions = append(positions, p)
	}
	return positions, nil
}
func (bs *Futures) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("not supported.")
}
func (bs *Futures) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	apiPath := "order"
	symbol, err := bs.adaptToSymbol(currencyPair, contractType)
	if err != nil {
		return nil, err
	}
	param := url.Values{}
	param.Set("symbol", symbol)
	param.Set("orderId", orderId)
	bs.base.buildParamsSigned(&param)
	reqUrl := fmt.Sprintf("%s%s?%s", bs.base.apiV1, apiPath, param.Encode())
	resp, err := HttpGet5(bs.base.httpClient, reqUrl, map[string]string{"X-MBX-APIKEY": bs.apikey})
	if err != nil {
		log.Error().Msgf("request url: %s", reqUrl)
		return nil, err
	}
	// logs.D(slice.Bytes2String(resp))
	var getOrderInfoResponse OrderInfoResponse
	err = json.Unmarshal(resp, &getOrderInfoResponse)
	if err != nil {
		log.Error().Msgf("response body: %s", slice.Bytes2String(resp))
		return nil, err
	}
	return &FutureOrder{
		Currency:     currencyPair,
		ClientOid:    getOrderInfoResponse.ClientOrderId,
		OrderID2:     fmt.Sprint(getOrderInfoResponse.OrderId),
		Price:        getOrderInfoResponse.Price,
		Amount:       getOrderInfoResponse.OrigQty,
		AvgPrice:     getOrderInfoResponse.AvgPrice,
		DealAmount:   getOrderInfoResponse.ExecutedQty,
		OrderTime:    getOrderInfoResponse.Time / 1000,
		Status:       bs.adaptStatus(getOrderInfoResponse.Status),
		OType:        bs.adaptOType(getOrderInfoResponse.Side, getOrderInfoResponse.PositionSide),
		ContractName: contractType,
		FinishedTime: getOrderInfoResponse.UpdateTime / 1000,
	}, nil
}
func (bs *Futures) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	apiPath := "openOrders"
	param := url.Values{}
	symbol, err := bs.adaptToSymbol(currencyPair, contractType)
	if err != nil {
		return nil, err
	}
	param.Set("symbol", symbol)
	bs.base.buildParamsSigned(&param)
	respbody, err := HttpGet5(bs.base.httpClient, fmt.Sprintf("%s%s?%s", bs.base.apiV1, apiPath, param.Encode()),
		map[string]string{
			"X-MBX-APIKEY": bs.apikey,
		})
	if err != nil {
		return nil, err
	}
	// logs.D(slice.Bytes2String(respbody))
	var (
		openOrderResponse []OrderInfoResponse
		orders            []FutureOrder
	)
	err = json.Unmarshal(respbody, &openOrderResponse)
	if err != nil {
		return nil, err
	}
	for _, ord := range openOrderResponse {
		orders = append(orders, FutureOrder{
			Currency:     currencyPair,
			ClientOid:    ord.ClientOrderId,
			OrderID:      ord.OrderId,
			OrderID2:     fmt.Sprint(ord.OrderId),
			Price:        ord.Price,
			Amount:       ord.OrigQty,
			AvgPrice:     ord.AvgPrice,
			DealAmount:   ord.ExecutedQty,
			Status:       bs.adaptStatus(ord.Status),
			OType:        bs.adaptOType(ord.Side, ord.PositionSide),
			ContractName: contractType,
			FinishedTime: ord.UpdateTime / 1000,
			OrderTime:    ord.Time / 1000,
		})
	}
	return orders, nil
}
func (bs *Futures) GetFutureOrderHistory(pair CurrencyPair, contractType string, optional ...OptionalParameter) ([]FutureOrder, error) {
	panic("implement me")
}
func (bs *Futures) GetFee() (float64, error) {
	panic("not supported.")
}
func (bs *Futures) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	switch currencyPair {
	case BTC_USD:
		return 100, nil
	default:
		return 10, nil
	}
}
func (bs *Futures) GetDeliveryTime() (int, int, int, int) {
	panic("not supported.")
}
func (bs *Futures) GetKlineRecords(contractType string, currency CurrencyPair, period KlinePeriod, size int, opt ...OptionalParameter) ([]FutureKline, error) {
	panic("not supported.")
}
func (bs *Futures) GetTrades(contractType string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not supported.")
}
func (bs *Futures) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	panic("not supported.")
}
func (bs *Futures) GetExchangeInfo() {
	exchangeInfoUri := bs.base.apiV1 + "exchangeInfo"
	ret, err := HttpGet5(bs.base.httpClient, exchangeInfoUri, map[string]string{})
	if err != nil {
		log.Error().Err(err).Msg("[exchangeInfo] Http Error")
		return
	}
	err = json.Unmarshal(ret, &bs.exchangeInfo)
	if err != nil {
		log.Error().Err(err).Bytes("response data", ret).Msg("json unmarshal response content error")
		return
	}
}
func (bs *Futures) adaptToSymbol(pair CurrencyPair, contractType string) (string, error) {
	if contractType == THIS_WEEK_CONTRACT || contractType == NEXT_WEEK_CONTRACT {
		return "", errors.New("binance only support contract quarter or bi_quarter")
	}
	if contractType == SWAP_CONTRACT {
		return fmt.Sprintf("%s_PERP", pair.AdaptUsdtToUsd().ToSymbol("")), nil
	}
	if bs.exchangeInfo == nil || len(bs.exchangeInfo.Symbols) == 0 {
		bs.GetExchangeInfo()
	}
	for _, info := range bs.exchangeInfo.Symbols {
		if info.ContractType != "PERPETUAL" &&
			info.ContractStatus == "TRADING" &&
			info.DeliveryDate <= time.Now().Unix()*1000 {
			log.Printf("pair=%s , contractType=%s, delivery date = %d ,  now= %d", info.Pair, info.ContractType, info.DeliveryDate, time.Now().Unix()*1000)
			bs.GetExchangeInfo()
		}
		if info.Pair == pair.ToSymbol("") {
			if info.ContractStatus != "TRADING" {
				return "", errors.New("contract status " + info.ContractStatus)
			}
			if info.ContractType == "CURRENT_QUARTER" && contractType == QUARTER_CONTRACT {
				return info.Symbol, nil
			}
			if info.ContractType == "NEXT_QUARTER" && contractType == BI_QUARTER_CONTRACT {
				return info.Symbol, nil
			}
			if info.Symbol == contractType {
				return info.Symbol, nil
			}
		}
	}
	return "", errors.New("binance not support " + pair.ToSymbol("") + " " + contractType)
}
func (bs *Futures) adaptStatus(status string) TradeStatus {
	switch status {
	case "NEW":
		return ORDER_UNFINISH
	case "CANCELED":
		return ORDER_CANCEL
	case "FILLED":
		return ORDER_FINISH
	case "PARTIALLY_FILLED":
		return ORDER_PART_FINISH
	case "PENDING_CANCEL":
		return ORDER_CANCEL_ING
	case "REJECTED":
		return ORDER_REJECT
	default:
		return ORDER_UNFINISH
	}
}
func (bs *Futures) adaptOType(side string, positionSide string) int {
	if positionSide == "BOTH" && side == "SELL" {
		return OPEN_SELL
	}
	if positionSide == "BOTH" && side == "BUY" {
		return OPEN_BUY
	}
	if positionSide == "LONG" {
		switch side {
		case "BUY":
			return OPEN_BUY
		default:
			return CLOSE_BUY
		}
	}
	if positionSide == "SHORT" {
		switch side {
		case "SELL":
			return OPEN_SELL
		default:
			return CLOSE_SELL
		}
	}
	return 0
}
