# exws

A Golang cryptocurrency websocket trading API with support for more than 10 bitcoin/altcoin exchanges

# usage:

## 本项目已完成自测代码100%覆盖 在各个目录下的_test文件中可找到使用方法
## This project has achieved 100% self-test code coverage. You can find usage examples in the _test files within each directory.

## update apikey info in ```config/testnet/apikey.go```

## 创建项目
## create project by import this repo
```go
package main

import (
	"github.com/conbanwa/exws"
	"github.com/conbanwa/exws/build"
	"github.com/conbanwa/exws/config"
)

//创建api
func NewCryptoMarket() exws.API {
	config.UseProxy = "localhost:7890"
	api := build.DefaultAPIBuilder.APIKey(apiKey).APISecretkey(Secretkey).ApiPassphrase(phrase).Build("alias")
	return api
}

//使用创建的api
func main() {
	api := NewCryptoMarket()
	name := api.String()
	pairs, err := api.PairArray()
	allTicker, err := api.AllTicker()
	fees, err := api.TradeFee()
}
```
