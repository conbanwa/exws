package bigone

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"qa3/wstrader"
	"qa3/wstrader/cons"
	"qa3/wstrader/q"
	"qa3/wstrader/web"
	"time"

	"github.com/conbanwa/num"
	"github.com/google/uuid"
	"github.com/nubo/jwt"
)

const (
	V2          = "https://big.one/api/v2"
	V3          = "https://big.one/api/v3"
	TICKER_URI  = "%s/markets/%s/ticker"
	DEPTH_URI   = "%s/markets/%s/depth"
	ACCOUNT_URI = "%s/viewer/accounts"
	ORDERS_URI  = "%s/viewer/orders"
)

type Bigone struct {
	accessKey,
	secretKey string
	httpClient *http.Client
	uid        string
	baseUri    string
	timeOffset int64
}

func New(client *http.Client, api_key, secret_key string) *Bigone {
	return &Bigone{accessKey: api_key, secretKey: secret_key, httpClient: client, uid: uuid.New().String(), baseUri: V2}
}
func (bo *Bigone) String() string {
	return cons.BIGONE
}

type TickerResp struct {
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
		Ask struct {
			Amount string `json:"amount"`
			Price  string `json:"price"`
		} `json:"ask"`
		Bid struct {
			Amount string `json:"amount"`
			Price  string `json:"price"`
		} `json:"bid"`
		Close           string `json:"close"`
		DailyChange     string `json:"daily_change"`
		DailyChangePerc string `json:"daily_change_perc"`
		High            string `json:"high"`
		Low             string `json:"low"`
		MarketID        string `json:"market_id"`
		MarketUUID      string `json:"market_uuid"`
		Open            string `json:"open"`
		Volume          string `json:"volume"`
	} `json:"data"`
}

func (bo *Bigone) GetTicker(currency cons.CurrencyPair) (*wstrader.Ticker, error) {
	tickerURI := fmt.Sprintf(TICKER_URI, bo.baseUri, currency.ToSymbol("-"))
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

type PlaceOrderResp struct {
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
		Amount        string `json:"amount"`
		AvgDealPrice  string `json:"avg_deal_price"`
		FilledAmount  string `json:"filled_amount"`
		ID            string `json:"id"`
		OrderID       int64  `json:"id"`
		InsertedAt    string `json:"inserted_at"`
		CreatedAt     string `json:"created_at"`
		MarketID      string `json:"market_id"`
		AssetPairName string `json:"asset_pair_name"`
		MarketUUID    string `json:"market_uuid"`
		Price         string `json:"price"`
		Side          string `json:"side"`
		State         string `json:"state"`
		UpdatedAt     string `json:"updated_at"`
	} `json:"data"`
}

func (bo *Bigone) placeOrder(amount, price string, pair cons.CurrencyPair, orderType, orderSide string) (*q.Order, error) {
	path := fmt.Sprintf(ORDERS_URI, bo.baseUri)
	params := make(map[string]string)
	params["market_id"] = pair.ToSymbol("-")
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
		OrderID2:   resp.Data.ID,
		Price:      num.ToFloat64(resp.Data.Price),
		Amount:     num.ToFloat64(resp.Data.Amount),
		DealAmount: 0,
		AvgPrice:   0,
		Side:       side,
		Status:     cons.ORDER_UNFINISH,
		OrderTime:  int(time.Now().Unix())}, nil
}
func (bo *Bigone) LimitBuy(amount, price string, currency cons.CurrencyPair, opt ...cons.LimitOrderOptionalParameter) (*q.Order, error) {
	return bo.placeOrder(amount, price, currency, "LIMIT", "BID")
}
func (bo *Bigone) LimitSell(amount, price string, currency cons.CurrencyPair, opt ...cons.LimitOrderOptionalParameter) (*q.Order, error) {
	return bo.placeOrder(amount, price, currency, "LIMIT", "ASK")
}
func (bo *Bigone) MarketBuy(amount, price string, currency cons.CurrencyPair) (*q.Order, error) {
	panic("not implements")
}
func (bo *Bigone) MarketSell(amount, price string, currency cons.CurrencyPair) (*q.Order, error) {
	panic("not implements")
}
func (bo *Bigone) privateHeader() map[string]string {
	claims := jwt.ClaimSet{
		"type":  "OpenAPI",
		"sub":   bo.accessKey,
		"nonce": time.Now().UnixNano(),
	}
	token, err := claims.Sign(bo.secretKey)
	if nil != err {
		log.Printf("privateHeader - cliam.Sign failed : %v", err)
		return nil
	}
	return map[string]string{"Authorization": "Bearer " + token}
}

type OrderListResp struct {
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
		Edges []struct {
			Cursor string `json:"cursor"`
			Node   struct {
				Amount       string `json:"amount"`
				AvgDealPrice string `json:"avg_deal_price"`
				FilledAmount string `json:"filled_amount"`
				ID           string `json:"id"`
				InsertedAt   string `json:"inserted_at"`
				MarketID     string `json:"market_id"`
				MarketUUID   string `json:"market_uuid"`
				Price        string `json:"price"`
				Side         string `json:"side"`
				State        string `json:"state"`
				UpdatedAt    string `json:"updated_at"`
			} `json:"node"`
		} `json:"edges"`
		PageInfo struct {
			EndCursor       string `json:"end_cursor"`
			HasNextPage     bool   `json:"has_next_page"`
			HasPreviousPage bool   `json:"has_previous_page"`
			StartCursor     string `json:"start_cursor"`
		} `json:"page_info"`
	} `json:"data"`
}

func (bo *Bigone) getOrdersList(currencyPair cons.CurrencyPair, size int, sts cons.TradeStatus) ([]q.Order, error) {
	apiURL := ""
	apiURL = fmt.Sprintf(ORDERS_URI+"?market_id=%s",
		bo.baseUri, currencyPair.ToSymbol("-"))
	if sts == cons.ORDER_FINISH {
		apiURL += "&state=FILLED"
	} else {
		apiURL += "&state=PENDING"
	}
	var resp OrderListResp
	err := web.HttpGet4(bo.httpClient, apiURL, bo.privateHeader(), &resp)
	if err != nil {
		log.Printf("getOrdersList - HttpGet4 failed : %v", err)
		return nil, err
	}
	orders := make([]q.Order, 0)
	for _, edge := range resp.Data.Edges {
		order := edge.Node
		status := order.State
		side := order.Side
		ord := q.Order{}
		switch status {
		case "PENDING":
			ord.Status = cons.ORDER_UNFINISH
		case "FILLED":
			ord.Status = cons.ORDER_FINISH
		case "CANCELED":
			ord.Status = cons.ORDER_CANCEL
		}
		if ord.Status != sts {
			continue // discard
		}
		ord.Currency = currencyPair
		ord.OrderID2 = order.ID
		if side == "ASK" {
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

type CancelOrderResp struct {
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
		ID            string `json:"id"`
		OrderID       string `json:"id"`
		MarketUUID    string `json:"market_uuid"`
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

func (bo *Bigone) CancelOrder(orderId string, currency cons.CurrencyPair) (bool, error) {
	path := fmt.Sprintf(ORDERS_URI+"/%s/cancel", bo.baseUri, orderId)
	params := make(map[string]string)
	params["order_id"] = orderId
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
func (bo *Bigone) GetOneOrder(orderId string, currencyPair cons.CurrencyPair) (*q.Order, error) {
	return nil, fmt.Errorf("GetOneOrder - not support yet")
}
func (bo *Bigone) GetUnfinishedOrders(currencyPair cons.CurrencyPair) ([]q.Order, error) {
	return bo.getOrdersList(currencyPair, -1, cons.ORDER_UNFINISH)
}
func (bo *Bigone) GetOrderHistorys(currencyPair cons.CurrencyPair, opt ...wstrader.OptionalParameter) ([]q.Order, error) {
	return bo.getOrdersList(currencyPair, -1, cons.ORDER_FINISH)
}

type AccountResp struct {
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
		AssetID       string `json:"asset_id"`
		AssetSymbol   string `json:"asset_symbol"`
		AssetUUID     string `json:"asset_uuid,omitempty"`
		Balance       string `json:"balance"`
		LockedBalance string `json:"locked_balance"`
	} `json:"data"`
}

func (bo *Bigone) GetAccount() (*wstrader.Account, error) {
	var resp AccountResp
	apiUrl := fmt.Sprintf(ACCOUNT_URI, bo.baseUri)
	err := web.HttpGet4(bo.httpClient, apiUrl, bo.privateHeader(), &resp)
	if err != nil {
		log.Println("GetAccount error:", err)
		return nil, err
	}
	acc := wstrader.Account{}
	acc.Exchange = bo.String()
	acc.SubAccounts = make(map[cons.Currency]wstrader.SubAccount)
	for _, v := range resp.Data {
		//log.Println(v)
		var currency cons.Currency
		if v.AssetID != "" {
			currency = cons.NewCurrency(v.AssetID, "")
		} else {
			currency = cons.NewCurrency(v.AssetSymbol, "")
		}
		acc.SubAccounts[currency] = wstrader.SubAccount{
			Currency:     currency,
			Amount:       num.ToFloat64(v.Balance),
			ForzenAmount: num.ToFloat64(v.LockedBalance),
		}
	}
	return &acc, nil
}

type DepthResp struct {
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
		MarketID      string `json:"market_id"`
		AssetPairName string `json:"asset_pair_name"`
		Bids          []struct {
			Price      string `json:"price"`
			OrderCount int    `json:"order_count"`
			Amount     string `json:"amount,omitempty"`
			Quantity   string `json:"quantity,omitempty"`
		} `json:"bids"`
		Asks []struct {
			Price      string `json:"price"`
			OrderCount int    `json:"order_count"`
			Amount     string `json:"amount,omitempty"`
			Quantity   string `json:"quantity,omitempty"`
		} `json:"asks"`
	}
}

func (bo *Bigone) GetDepth(size int, currencyPair cons.CurrencyPair) (*wstrader.Depth, error) {
	var resp DepthResp
	apiURL := fmt.Sprintf(DEPTH_URI, bo.baseUri, currencyPair.ToSymbol("-"))
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
	return depth, nil
}
func (bo *Bigone) GetKlineRecords(currency cons.CurrencyPair, period cons.KlinePeriod, size int, opt ...wstrader.OptionalParameter) ([]wstrader.Kline, error) {
	panic("not implements")
}

// 非个人，整个交易所的交易记录
func (bo *Bigone) GetTrades(currencyPair cons.CurrencyPair, since int64) ([]q.Trade, error) {
	panic("not implements")
}
