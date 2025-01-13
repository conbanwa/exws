package build

import (
	"context"
	"errors"
	"fmt"
	"github.com/conbanwa/exws"
	"github.com/conbanwa/exws/config"
	"github.com/conbanwa/exws/cons"
	"github.com/conbanwa/exws/ex/bigone"
	"github.com/conbanwa/exws/ex/binance"
	"github.com/conbanwa/exws/ex/bitfinex"
	"github.com/conbanwa/exws/ex/bithumb"
	"github.com/conbanwa/exws/ex/bitmex"
	"github.com/conbanwa/exws/ex/bitstamp"
	"github.com/conbanwa/exws/ex/coinbase"
	"github.com/conbanwa/exws/ex/coinex"
	"github.com/conbanwa/exws/ex/gateio"
	"github.com/conbanwa/exws/ex/hitbtc"
	"github.com/conbanwa/exws/ex/huobi"
	"github.com/conbanwa/exws/ex/kraken"
	"github.com/conbanwa/exws/ex/kucoin"
	"github.com/conbanwa/exws/ex/okx"
	"github.com/conbanwa/exws/ex/okx/okex"
	"github.com/conbanwa/exws/ex/poloniex"
	"github.com/conbanwa/exws/stat/zelo"
	"net"
	"net/http"
	"net/url"
	"time"
)

var log = zelo.Writer

type APIBuilder struct {
	HttpClientConfig *config.HttpClientConfig
	client           *http.Client
	httpTimeout      time.Duration
	apiKey           string
	secretKey        string
	clientId         string
	apiPassphrase    string
	futuresEndPoint  string
	endPoint         string
	futuresLever     float64
}

var DefaultAPIBuilder = NewAPIBuilder()

func init() {
	config.SetProxy()
}

func NewAPIBuilder() (builder *APIBuilder) {
	return &APIBuilder{
		HttpClientConfig: config.DefaultHttpClientConfig,
		client: &http.Client{
			Timeout: config.DefaultHttpClientConfig.HttpTimeout,
			Transport: &http.Transport{
				Proxy: func(request *http.Request) (*url.URL, error) {
					return config.DefaultHttpClientConfig.Proxy, nil
				},
				MaxIdleConns:          config.DefaultHttpClientConfig.MaxIdleConns,
				IdleConnTimeout:       5 * config.DefaultHttpClientConfig.HttpTimeout,
				MaxConnsPerHost:       2,
				MaxIdleConnsPerHost:   2,
				TLSHandshakeTimeout:   config.DefaultHttpClientConfig.HttpTimeout,
				ResponseHeaderTimeout: config.DefaultHttpClientConfig.HttpTimeout,
				ExpectContinueTimeout: config.DefaultHttpClientConfig.HttpTimeout,
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.DialTimeout(network, addr, config.DefaultHttpClientConfig.HttpTimeout)
				}},
		}}
}

func NewCustomAPIBuilder(client *http.Client) (builder *APIBuilder) {
	return &APIBuilder{client: client}
}

func (builder *APIBuilder) BuildSpotWs(exName string) (exws.SpotWsApi, error) {
	switch exName {
	case cons.HUOBI_PRO, cons.HUOBI:
		return huobi.NewSpotWs(), nil
	case cons.BINANCE:
		return binance.NewSpotWs(), nil
	case cons.OKEX:
		return okx.NewSpotWs(), nil
	case cons.BITFINEX:
		return bitfinex.NewWs(), nil
	case cons.GATEIO:
		return gateio.NewWs(), nil
	}
	return nil, errors.New("not support the exchange " + exName)
}

func (builder *APIBuilder) GetHttpClientConfig() *config.HttpClientConfig {
	return builder.HttpClientConfig
}
func (builder *APIBuilder) GetHttpClient() *http.Client {
	return builder.client
}
func (builder *APIBuilder) HttpProxy(proxyUrl string) (_builder *APIBuilder) {
	if proxyUrl == "" {
		return builder
	}
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		return builder
	}
	builder.HttpClientConfig.Proxy = proxy
	transport := builder.client.Transport.(*http.Transport)
	transport.Proxy = http.ProxyURL(proxy)
	return builder
}
func (builder *APIBuilder) HttpTimeout(timeout time.Duration) (_builder *APIBuilder) {
	builder.HttpClientConfig.HttpTimeout = timeout
	builder.httpTimeout = timeout
	builder.client.Timeout = timeout
	if transport := builder.client.Transport.(*http.Transport); transport != nil {
		//transport.ResponseHeaderTimeout = timeout
		//transport.TLSHandshakeTimeout = timeout
		transport.IdleConnTimeout = timeout
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		}
	}
	return builder
}
func (builder *APIBuilder) APIKey(key string) (_builder *APIBuilder) {
	builder.apiKey = key
	return builder
}
func (builder *APIBuilder) APISecretkey(key string) (_builder *APIBuilder) {
	builder.secretKey = key
	return builder
}
func (builder *APIBuilder) ClientID(id string) (_builder *APIBuilder) {
	builder.clientId = id
	return builder
}
func (builder *APIBuilder) ApiPassphrase(apiPassphrase string) (_builder *APIBuilder) {
	builder.apiPassphrase = apiPassphrase
	return builder
}
func (builder *APIBuilder) FuturesEndpoint(endpoint string) (_builder *APIBuilder) {
	builder.futuresEndPoint = endpoint
	return builder
}
func (builder *APIBuilder) Endpoint(endpoint string) (_builder *APIBuilder) {
	builder.endPoint = endpoint
	return builder
}
func (builder *APIBuilder) FuturesLever(lever float64) (_builder *APIBuilder) {
	builder.futuresLever = lever
	return builder
}
func (builder *APIBuilder) Build(exName string) (api exws.API) {
	var _api exws.API
	switch exName {
	case cons.KUCOIN:
		_api = kucoin.New(builder.apiKey, builder.secretKey, builder.apiPassphrase)
	//case OKCOIN_CN:
	//	_api = okcoin.New(builder.client, builder.apiKey, builder.secretKey)
	case cons.POLONIEX:
		_api = poloniex.New(builder.client, builder.apiKey, builder.secretKey)
	//case OKCOIN_COM:
	//	_api = okcoin.NewCOM(builder.client, builder.apiKey, builder.secretKey)
	case cons.BITSTAMP:
		_api = bitstamp.NewBitstamp(builder.client, builder.apiKey, builder.secretKey, builder.clientId)
	case cons.HUOBI_PRO:
		//_api = huobi.NewHuoBiProSpot(builder.client, builder.apiKey, builder.secretKey)
		_api = huobi.NewHuobiWithConfig(&exws.APIConfig{
			HttpClient:   builder.client,
			Endpoint:     builder.endPoint,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretKey})
	case cons.OKEX:
		_api = okx.NewOKExV5Spot(&exws.APIConfig{
			HttpClient:    builder.client,
			ApiKey:        builder.apiKey,
			ApiSecretKey:  builder.secretKey,
			ApiPassphrase: builder.apiPassphrase,
			Endpoint:      builder.endPoint,
		})
	case cons.BITFINEX:
		_api = bitfinex.New(builder.client, builder.apiKey, builder.secretKey)
	case cons.KRAKEN:
		_api = kraken.New(builder.client, builder.apiKey, builder.secretKey)
	case cons.BINANCE:
		//_api = binance.New(builder.client, builder.apiKey, builder.secretKey)
		_api = binance.NewWithConfig(&exws.APIConfig{
			HttpClient:   builder.client,
			Endpoint:     builder.endPoint,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretKey})
	case cons.BITHUMB:
		_api = bithumb.New(builder.client, builder.apiKey, builder.secretKey)
	case cons.COINBASE:
		_api = coinbase.New(builder.client, builder.apiKey, builder.secretKey)
	case cons.COINEX:
		_api = coinex.New(builder.client, builder.apiKey, builder.secretKey)
	case cons.BIGONE:
		_api = bigone.New(builder.client, builder.apiKey, builder.secretKey)
	case cons.HITBTC:
		_api = hitbtc.New(builder.client, builder.apiKey, builder.secretKey)
	default:
		log.Warn().Str("ex", exName).Msg("exchange name error")
	}
	return _api
}
func (builder *APIBuilder) BuildFuture(exName string) (api exws.FutureRestAPI) {
	switch exName {
	case cons.BITMEX:
		return bitmex.New(&exws.APIConfig{
			//Endpoint:     "https://www.bitmex.com/",
			Endpoint:     builder.futuresEndPoint,
			HttpClient:   builder.client,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretKey})
	case cons.BITMEX_TEST:
		return bitmex.New(&exws.APIConfig{
			HttpClient:   builder.client,
			Endpoint:     "https://testnet.bitmex.com",
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretKey,
		})
	case cons.OKEX_FUTURE:
		//return okcoin.NewOKEx(builder.client, builder.apiKey, builder.secretKey)
		return okex.NewOKEx(&exws.APIConfig{
			HttpClient: builder.client,
			//	Endpoint:      "https://www.okx.com",
			Endpoint:      builder.futuresEndPoint,
			ApiKey:        builder.apiKey,
			ApiSecretKey:  builder.secretKey,
			ApiPassphrase: builder.apiPassphrase,
			Lever:         builder.futuresLever}).OKExFuture
	case cons.HBDM:
		return huobi.NewHbdm(&exws.APIConfig{
			HttpClient:   builder.client,
			Endpoint:     builder.futuresEndPoint,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretKey,
			Lever:        builder.futuresLever})
	case cons.HBDM_SWAP:
		return huobi.NewHbdmSwap(&exws.APIConfig{
			HttpClient:   builder.client,
			Endpoint:     builder.endPoint,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretKey,
			Lever:        builder.futuresLever,
		})
	case cons.OKEX_SWAP:
		return okex.NewOKEx(&exws.APIConfig{
			HttpClient:    builder.client,
			Endpoint:      builder.futuresEndPoint,
			ApiKey:        builder.apiKey,
			ApiSecretKey:  builder.secretKey,
			ApiPassphrase: builder.apiPassphrase,
			Lever:         builder.futuresLever}).OKExSwap
	case cons.BINANCE_SWAP:
		return binance.NewBinanceSwap(&exws.APIConfig{
			HttpClient:   builder.client,
			Endpoint:     builder.futuresEndPoint,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretKey,
			Lever:        builder.futuresLever,
		})
	case cons.BINANCE, cons.BINANCE_FUTURES:
		return binance.NewBinanceFutures(&exws.APIConfig{
			HttpClient:   builder.client,
			Endpoint:     builder.futuresEndPoint,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretKey,
			Lever:        builder.futuresLever,
		})
	default:
		println(fmt.Sprintf("%s not support future", exName))
		return nil
	}
}
func (builder *APIBuilder) BuildFuturesWs(exName string) (exws.FuturesWsApi, error) {
	switch exName {
	case cons.OKEX, cons.OKEX_FUTURE:
		return okex.NewOKExV3FuturesWs(okex.NewOKEx(&exws.APIConfig{
			HttpClient: builder.client,
			Endpoint:   builder.futuresEndPoint,
		})), nil
	case cons.HBDM:
		return huobi.NewHbdmWs(), nil
	case cons.HBDM_SWAP:
		return huobi.NewHbdmSwapWs(), nil
	case cons.BINANCE, cons.BINANCE_FUTURES, cons.BINANCE_SWAP:
		return binance.NewFuturesWs(), nil
	case cons.BITMEX:
		return bitmex.NewSwapWs(), nil
	}
	return nil, errors.New("not support the exchange " + exName)
}
func (builder *APIBuilder) BuildWallet(exName string) (exws.WalletApi, error) {
	switch exName {
	case cons.OKEX:
		return okex.NewOKEx(&exws.APIConfig{
			HttpClient:    builder.client,
			ApiKey:        builder.apiKey,
			ApiSecretKey:  builder.secretKey,
			ApiPassphrase: builder.apiPassphrase,
		}).OKExWallet, nil
	case cons.HUOBI_PRO:
		return huobi.NewWallet(&exws.APIConfig{
			HttpClient:   builder.client,
			Endpoint:     builder.endPoint,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretKey,
		}), nil
	case cons.BINANCE:
		return binance.NewWallet(&exws.APIConfig{
			HttpClient:   builder.client,
			Endpoint:     builder.endPoint,
			ApiKey:       builder.apiKey,
			ApiSecretKey: builder.secretKey,
		}), nil
	}
	return nil, errors.New("not support the wallet api for  " + exName)
}
