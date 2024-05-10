package wstrader

import (
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
	"sync"
)

type API interface {
	CancelOrder(orderId string, currency cons.CurrencyPair) (bool, error)
	// GetOneOrder(orderId string, currency CurrencyPair) (*Order, error)
	GetUnfinishedOrders(currency cons.CurrencyPair) ([]q.Order, error)
	// GetOrderHistorys(currency CurrencyPair, opt ...OptionalParameter) ([]Order, error)
	// GetAccount() (*Account, error)
	// GetDepth(size int, currency CurrencyPair) (*Depth, error)
	// GetKlineRecords(currency CurrencyPair, period KlinePeriod, size int, optional ...OptionalParameter) ([]Kline, error)
	//none-personal but whole exchange center
	// GetTicker(currency CurrencyPair) (*Ticker, error)
	String() string
	OneTicker(d q.D) (q.Bbo, error)
	AllTicker(SymPair map[string]q.D) (*sync.Map, error)
	PairArray() (map[string]q.D, map[q.D]q.P, error)
	Fee() float64
	TradeFee() (map[string]q.TradeFee, error)
	WithdrawFee() ([]q.NetworkWithdraw, error)
	PlaceOrders([3]q.Order) ([3]q.Order, error)
	Balances() (*sync.Map, *sync.Map, error)
	GetAttr() q.Attr
	Test() bool
	// Order([3]string)
	// Correct(string)string
}
type FuturesWsApi interface {
	DepthCallback(func(depth *Depth))
	TickerCallback(func(ticker *FutureTicker))
	TradeCallback(func(trade *q.Trade, contract string))
	//OrderCallback(func(order *FutureOrder))
	//PositionCallback(func(position *FuturePosition))
	//AccountCallback(func(account *FutureAccount))
	SubscribeDepth(pair cons.CurrencyPair, contractType string) error
	SubscribeTicker(pair cons.CurrencyPair, contractType string) error
	SubscribeTrade(pair cons.CurrencyPair, contractType string) error
	//Login() error
	//SubscribeOrder(pair CurrencyPair, contractType string) error
	//SubscribePosition(pair CurrencyPair, contractType string) error
	//SubscribeAccount(pair CurrencyPair) error
}
type SpotWsApi interface {
	DepthCallback(func(depth *Depth))
	TickerCallback(func(ticker *Ticker))
	TradeCallback(func(trade *q.Trade))
	BBOCallback(func(bbo *q.Bbo))
	//OrderCallback(func(order *Order))
	//AccountCallback(func(account *Account))
	SubscribeDepth(pair cons.CurrencyPair) error
	SubscribeTicker(pair cons.CurrencyPair) error
	SubscribeTrade(pair cons.CurrencyPair) error
	SubscribeBBO(sm []string) error
	//Login() error
	//SubscribeOrder(pair CurrencyPair) error
	//SubscribeAccount(pair CurrencyPair) error
}
