# wstrader

A Golang cryptocurrency websocket trading API with support for more than 10 bitcoin/altcoin exchanges

# usage:

## 本项目已完成自测代码100%覆盖 在各个目录下的_test文件中可找到使用方法

## 创建项目
```go
package main

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/build"
	"github.com/conbanwa/wstrader/config"
)
//创建api
func NewCryptoMarket() wstrader.API {
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