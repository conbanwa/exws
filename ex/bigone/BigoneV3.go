package bigone

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"qa3/wstrader"
	"qa3/wstrader/cons"
	"qa3/wstrader/q"
	"qa3/wstrader/web"
	"time"

	"github.com/conbanwa/num"
	"github.com/google/uuid"
	"github.com/nubo/jwt"
)

type BigoneV3 struct {
	Bigone
}

// accessKey,
// secretKey string
// httpClient *http.Client
// uid        string
// baseUri    string
var _INERNAL_KLINE_PERIOD_CONVERTER = map[int]string{
	cons.KLINE_PERIOD_1MIN:   "min1",
	cons.KLINE_PERIOD_5MIN:   "min5",
	cons.KLINE_PERIOD_15MIN:  "min15",
	cons.KLINE_PERIOD_30MIN:  "min30",
	cons.KLINE_PERIOD_60MIN:  "hour1",
	cons.KLINE_PERIOD_4H:     "hour4",
	cons.KLINE_PERIOD_6H:     "hour6",
	cons.KLINE_PERIOD_12H:    "hour12",
	cons.KLINE_PERIOD_1DAY:   "day1",
	cons.KLINE_PERIOD_1WEEK:  "week1",
	cons.KLINE_PERIOD_1MONTH: "month1",
}

func NewV3(client *http.Client, api_key, secret_key string) *BigoneV3 {
	b1 := &BigoneV3{}
	b1.secretKey = secret_key
	b1.accessKey = api_key
	b1.httpClient = client
	b1.uid = uuid.New().String()
	b1.baseUri = V3
	b1.setTimeOffset()
	return b1
}
func (bo *BigoneV3) String() string {
	return cons.BIGONE
}

type ServerTimestampResp struct {
	Errors []struct {
		Code      int `json:"code"`
		Locations []struct {
			Column int `json:"column"`
			Line   int `json:"line"`
		} `json:"locations"`
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data struct {
		Timetamp int64 `json:"timestamp"`
	} `json:"data"`
}

func (bo *BigoneV3) setTimeOffset() {
	pingUri := fmt.Sprintf("%s/ping", bo.baseUri)
	var resp ServerTimestampResp
	//log.Printf("GetPing -> %s", pingUri)
	err := web.HttpGet4(bo.httpClient, pingUri, nil, &resp)
	if err != nil {
		log.Printf("GetPing - HttpGet4 failed : %v", err)
		return
	}
	bo.timeOffset = time.Now().UnixNano() - resp.Data.Timetamp
	//log.Println(resp)
	return
}
func (bo *BigoneV3) GetTicker(currency cons.CurrencyPair) (*wstrader.Ticker, error) {
	params := url.Values{}
	params.Set("asset_pair_name", currency.ToSymbol("-"))
	tickerURI := fmt.Sprintf("%s/asset_pairs/%s/ticker?%s", bo.baseUri, currency.ToSymbol("-"), params.Encode())
	var resp TickerResp
	//log.Printf("GetTicker -> %s", tickerURI)
	err := web.HttpGet4(bo.httpClient, tickerURI, nil, &resp)
	if err != nil {
		log.Printf("GetTicker - HttpGet4 failed : %v", err)
		return nil, err
	}
	var ticker wstrader.Ticker
	ticker.Pair = currency
	ticker.Date = uint64(time.Now().Unix())
	ticker.Last = num.ToFloat64(resp.Data.Close)
	ticker.Buy = num.ToFloat64(resp.Data.Bid.Price)
	ticker.Sell = num.ToFloat64(resp.Data.Ask.Price)
	ticker.Low = num.ToFloat64(resp.Data.Low)
	ticker.High = num.ToFloat64(resp.Data.High)
	ticker.Vol = num.ToFloat64(resp.Data.Volume)
	return &ticker, nil
}
func (bo *BigoneV3) placeOrder(amount, price string, pair cons.CurrencyPair, orderType, orderSide string) (*q.Order, error) {
	path := fmt.Sprintf(ORDERS_URI, bo.baseUri)
	params := make(map[string]string)
	params["asset_pair_name"] = pair.ToSymbol("-")
	params["side"] = orderSide
	params["amount"] = amount
	params["price"] = price
	var resp PlaceOrderResp
	buf, err := web.HttpPostForm4(bo.httpClient, path, params, bo.privateHeader())
	if err != nil {
		log.Printf("placeOrder - HttpPostForm4 failed : %v", err)
		return nil, err
	}
	if err = json.Unmarshal(buf, &resp); nil != err {
		log.Printf("buf : %s", string(buf))
		log.Printf("placeOrder - json.Unmarshal failed : %v", err)
		return nil, err
	}
	if len(resp.Errors) > 0 {
		log.Printf("placeOrder - failed : %v", resp.Errors)
		return nil, fmt.Errorf(resp.Errors[0].Message)
	}
	side := cons.BUY
	if orderSide == "ASK" {
		side = cons.SELL
	}
	return &q.Order{
		Currency:   pair,
		OrderID:    int(resp.Data.OrderID),
		OrderID2:   fmt.Sprint(resp.Data.OrderID),
		Price:      num.ToFloat64(resp.Data.Price),
		Amount:     num.ToFloat64(resp.Data.Amount),
		DealAmount: 0,
		AvgPrice:   num.ToFloat64(resp.Data.AvgDealPrice),
		Side:       side,
		Status:     cons.ORDER_UNFINISH,
		OrderTime:  int(time.Now().Unix())}, nil
}
func (bo *BigoneV3) LimitBuy(amount, price string, currency cons.CurrencyPair) (*q.Order, error) {
	return bo.placeOrder(amount, price, currency, "LIMIT", "BID")
}
func (bo *BigoneV3) LimitSell(amount, price string, currency cons.CurrencyPair) (*q.Order, error) {
	return bo.placeOrder(amount, price, currency, "LIMIT", "ASK")
}
func (bo *BigoneV3) MarketBuy(amount, price string, currency cons.CurrencyPair) (*q.Order, error) {
	panic("not implements")
}
func (bo *BigoneV3) MarketSell(amount, price string, currency cons.CurrencyPair) (*q.Order, error) {
	panic("not implements")
}
func (bo *BigoneV3) privateHeader() map[string]string {
	claims := jwt.ClaimSet{
		"type":  "OpenAPI",
		"sub":   bo.accessKey,
		"nonce": time.Now().UnixNano() - bo.timeOffset,
	}
	token, err := claims.Sign(bo.secretKey)
	if nil != err {
		log.Printf("privateHeader - cliam.Sign failed : %v", err)
		return nil
	}
	return map[string]string{"Authorization": "Bearer " + token}
}

type OrderListV3Resp struct {
	Errors []struct {
		Code      int `json:"code"`
		Locations []struct {
			Column int `json:"column"`
			Line   int `json:"line"`
		} `json:"locations"`
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data []struct {
		ID            int64  `json:"id"`
		AssetPairName string `json:"asset_pair_name"`
		Price         string `json:"price"`
		Amount        string `json:"amount"`
		FilledAmount  string `json:"filled_amount"`
		AvgDealPrice  string `json:"avg_deal_price"`
		Side          string `json:"side"`
		State         string `json:"state"`
		CreatedAt     string `json:"created_at"`
		UpdatedAt     string `json:"updated_at"`
	} `json:"data"`
	PageToken string `json:"page_token"`
}

func (bo *BigoneV3) getOrdersList(currencyPair cons.CurrencyPair, size int, sts cons.TradeStatus) ([]q.Order, error) {
	apiURL := fmt.Sprintf(ORDERS_URI+"?asset_pair_name=%s&limit=%d",
		bo.baseUri, currencyPair.ToSymbol("-"), size)
	if sts == cons.ORDER_FINISH {
		apiURL += "&state=FILLED"
	} else {
		apiURL += "&state=PENDING"
	}
	//log.Printf("getOrdersList -> %s", apiURL)
	var resp OrderListV3Resp
	err := web.HttpGet4(bo.httpClient, apiURL, bo.privateHeader(), &resp)
	if err != nil {
		log.Printf("getOrdersList - HttpGet4 failed : %v", err)
		return nil, err
	}
	orders := make([]q.Order, 0)
	for _, order := range resp.Data {
		ord := q.Order{}
		switch order.State {
		case "PENDING":
			ord.Status = cons.ORDER_UNFINISH
		case "FILLED":
			ord.Status = cons.ORDER_FINISH
		case "CANCELLED":
			ord.Status = cons.ORDER_CANCEL
		}
		if ord.Status != sts {
			continue // discard
		}
		ord.Currency = currencyPair
		ord.OrderID2 = fmt.Sprint(order.ID)
		ord.OrderID = int(order.ID)
		if order.Side == "ASK" {
			ord.Side = cons.SELL
		} else {
			ord.Side = cons.BUY
		}
		ord.Amount = num.ToFloat64(order.Amount)
		ord.Price = num.ToFloat64(order.Price)
		ord.DealAmount = num.ToFloat64(order.FilledAmount)
		ord.AvgPrice = num.ToFloat64(order.Price)
		orders = append(orders, ord)
	}
	return orders, nil
}
func (bo *BigoneV3) CancelOrder(orderId string, currency cons.CurrencyPair) (bool, error) {
	path := fmt.Sprintf(ORDERS_URI+"/%s/cancel", bo.baseUri, orderId)
	params := make(map[string]string)
	params["id"] = orderId
	buf, err := web.HttpPostForm4(bo.httpClient, path, params, bo.privateHeader())
	if err != nil {
		log.Printf("CancelOrder - faield : %v", err)
		return false, err
	}
	var resp CancelOrderResp
	if err = json.Unmarshal(buf, &resp); nil != err {
		log.Printf("CancelOrder - json.Unmarshal failed : %v", err)
		return false, err
	}
	if len(resp.Errors) > 0 {
		log.Printf("getOrdersList - response error : %v", resp.Errors)
		return false, fmt.Errorf("%s", resp.Errors[0].Message)
	}
	return true, nil
}

type GetOneOrderResp struct {
	Errors []struct {
		Code      int `json:"code"`
		Locations []struct {
			Column int `json:"column"`
			Line   int `json:"line"`
		} `json:"locations"`
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data struct {
		OrderID       int64  `json:"id"`
		AssetPairName string `json:"asset_pair_name"`
		Price         string `json:"price"`
		Amount        string `json:"amount"`
		FilledAmount  string `json:"filled_amount"`
		AvgDealPrice  string `json:"avg_deal_price"`
		Side          string `json:"side"`
		State         string `json:"state"`
		CreatedAt     string `json:"created_at"`
		UpdatedAt     string `json:"updated_at"`
	}
}

func (bo *BigoneV3) GetOneOrder(orderId string, currencyPair cons.CurrencyPair) (*q.Order, error) {
	path := fmt.Sprintf(ORDERS_URI+"/%s?id=%s", bo.baseUri, orderId, orderId)
	//log.Printf("GetOneOrder -> %s", path)
	var resp GetOneOrderResp
	err := web.HttpGet4(bo.httpClient, path, bo.privateHeader(), &resp)
	if err != nil {
		log.Printf("GetOneOrder - faield : %v", err)
		return nil, err
	}
	state := cons.ORDER_UNFINISH
	switch resp.Data.State {
	case "PENDING":
		state = cons.ORDER_UNFINISH
	case "FILLED":
		state = cons.ORDER_FINISH
	case "CANCELLED":
		state = cons.ORDER_CANCEL
	}
	side := cons.BUY
	if resp.Data.Side == "ASK" {
		side = cons.SELL
	}
	return &q.Order{
		Price:      num.ToFloat64(resp.Data.Price),
		Amount:     num.ToFloat64(resp.Data.Amount),
		AvgPrice:   num.ToFloat64(resp.Data.AvgDealPrice),
		DealAmount: num.ToFloat64(resp.Data.FilledAmount),
		OrderID2:   fmt.Sprint(resp.Data.OrderID),
		OrderID:    int(resp.Data.OrderID),
		Status:     state,
		Currency:   currencyPair,
		Side:       side,
	}, nil
}
func (bo *BigoneV3) GetUnfinishedOrders(currencyPair cons.CurrencyPair) ([]q.Order, error) {
	return bo.getOrdersList(currencyPair, 200, cons.ORDER_UNFINISH)
}
func (bo *BigoneV3) GetOrderHistorys(currencyPair cons.CurrencyPair, opt wstrader.OptionalParameter) ([]q.Order, error) {
	return bo.getOrdersList(currencyPair, 200, cons.ORDER_FINISH)
}
func (bo *BigoneV3) GetAccount() (*wstrader.Account, error) {
	var resp AccountResp
	apiUrl := fmt.Sprintf(ACCOUNT_URI, bo.baseUri)
	err := web.HttpGet4(bo.httpClient, apiUrl, bo.privateHeader(), &resp)
	if err != nil {
		log.Println("GetAccount error:", err)
		return nil, err
	}
	//logs.D(resp)
	acc := wstrader.Account{}
	acc.Exchange = bo.String()
	acc.SubAccounts = make(map[cons.Currency]wstrader.SubAccount)
	for _, v := range resp.Data {
		//log.Println(v)
		currency := cons.NewCurrency(v.AssetSymbol, "")
		acc.SubAccounts[currency] = wstrader.SubAccount{
			Currency:     currency,
			Amount:       num.ToFloat64(v.Balance),
			ForzenAmount: num.ToFloat64(v.LockedBalance),
		}
	}
	return &acc, nil
}
func (bo *BigoneV3) GetDepth(size int, currencyPair cons.CurrencyPair) (*wstrader.Depth, error) {
	var resp DepthResp
	params := url.Values{}
	params.Set("asset_pair_name", currencyPair.ToSymbol("-"))
	params.Set("limit", fmt.Sprint(size))
	apiURL := fmt.Sprintf("%s/asset_pairs/%s/depth?%s", bo.baseUri, currencyPair.ToSymbol("-"), params.Encode())
	//log.Printf("GetDepth -> %s", apiURL)
	err := web.HttpGet4(bo.httpClient, apiURL, nil, &resp)
	if err != nil {
		log.Println("GetDepth error:", err)
		return nil, err
	}
	depth := new(wstrader.Depth)
	for _, bid := range resp.Data.Bids {
		var amount float64
		if bid.Amount != "" {
			amount = num.ToFloat64(bid.Amount)
		} else {
			amount = num.ToFloat64(bid.Quantity)
		}
		price := num.ToFloat64(bid.Price)
		dr := wstrader.DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
	}
	for _, ask := range resp.Data.Asks {
		var amount float64
		if ask.Amount != "" {
			amount = num.ToFloat64(ask.Amount)
		} else {
			amount = num.ToFloat64(ask.Quantity)
		}
		price := num.ToFloat64(ask.Price)
		dr := wstrader.DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
	}
	depth.Pair = currencyPair
	depth.UTime = time.Now()
	return depth, nil
}

type CandleResp struct {
	Errors []struct {
		Code      int `json:"code"`
		Locations []struct {
			Column int `json:"column"`
			Line   int `json:"line"`
		} `json:"locations"`
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data []struct {
		Close  string `json:"close"`
		High   string `json:"high"`
		Low    string `json:"low"`
		Open   string `json:"open"`
		Time   string `json:"time"`
		Volume string `json:"volume"`
	} `json:"data"`
}

func (bo *BigoneV3) GetKlineRecords(currency cons.CurrencyPair, period, size, since int) ([]wstrader.Kline, error) {
	apiUrl := fmt.Sprintf("%s/asset_pairs/%s/candles", bo.baseUri, currency.ToSymbol("-"))
	params := url.Values{}
	params.Set("asset_pair_name", currency.ToSymbol("-"))
	params.Set("period", _INERNAL_KLINE_PERIOD_CONVERTER[period])
	params.Set("limit", fmt.Sprint(size))
	//params["period"] = _INERNAL_KLINE_PERIOD_CONVERTER[period]
	//params["time"] =
	//params["limit"] = fmt.Sprint(size)
	var resp CandleResp
	err := web.HttpGet4(bo.httpClient, apiUrl+"?"+params.Encode(), bo.privateHeader(), &resp)
	if err != nil {
		log.Printf("GetKlineRecords - HttpGet4 failed : %v", err)
		return nil, err
	}
	klines := make([]wstrader.Kline, 0)
	for _, v := range resp.Data {
		ts, _ := time.Parse("2006-01-02T15:04:05Z", v.Time)
		klines = append(klines, wstrader.Kline{
			Pair:      currency,
			Open:      num.ToFloat64(v.Open),
			Close:     num.ToFloat64(v.Close),
			High:      num.ToFloat64(v.High),
			Low:       num.ToFloat64(v.Low),
			Vol:       num.ToFloat64(v.Volume),
			Timestamp: ts.Unix(),
		})
	}
	return klines, nil
}

// 非个人，整个交易所的交易记录
func (bo *BigoneV3) GetTrades(currencyPair cons.CurrencyPair, since int64) ([]q.Trade, error) {
	panic("not implements")
}
