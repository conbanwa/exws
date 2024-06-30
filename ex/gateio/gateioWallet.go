package gateio

import (
	"context"
	"encoding/json"
	"github.com/conbanwa/logs"
	"github.com/conbanwa/num"
	"github.com/gateio/gateapi-go/v6"
	"sync"
)

func (g *Gate) Balances() (availables, frozens *sync.Map, err error) {
	availables, frozens = new(sync.Map), new(sync.Map)
	var method = "POST"
	var url = "https://api.gateio.life/api2/1/private/balances"
	var param = ""
	resp := g.httpDo(method, url, param)
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil || respMap["result"] != "true" {
		logs.I(err, respMap)
		return
	}
	for kind, wallet := range respMap {
		switch kind {
		case "available":
			af, ok := wallet.(map[string]any)
			if !ok {
				logs.E(wallet)
			}
			for k, v := range af {
				availables.Store(k, num.ToFloat64(v))
			}
		case "locked":
			af, ok := wallet.(map[string]any)
			if !ok {
				logs.E(wallet)
			}
			for k, v := range af {
				frozens.Store(k, num.ToFloat64(v))
			}
		}
	}
	return
}
func (g *Gate) Fee() float64 {
	client := gateapi.NewAPIClient(gateapi.NewConfiguration())
	// uncomment the next line if you are testing against testnet
	// client.ChangeBasePath("https://fx-api-testnet.gateio.ws/api/v4")
	ctx := context.WithValue(context.Background(),
		gateapi.ContextGateAPIV4, gateapi.GateAPIV4{
			Key:    g.accesskey,
			Secret: g.secretkey,
		})
	result, _, err := client.SpotApi.GetFee(ctx, nil)
	if err != nil {
		if e, ok := err.(gateapi.GateAPIError); ok {
			log.Debug().Msgf("gate api error: %s\n", e.Error())
		} else {
			log.Debug().Msgf("generic error: %s\n", err.Error())
		}
	} else {
		logs.D(result)
	}
	return 0.003
}
