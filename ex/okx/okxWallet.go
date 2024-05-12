package okx

import (
	"fmt"
	"net/http"
	"github.com/conbanwa/wstrader/q"
	"github.com/conbanwa/wstrader/stat/zelo"
	"sync"
)

func (ok *OKX) Balances() (availables, frozens *sync.Map, err error) {
	availables, frozens = new(sync.Map), new(sync.Map)
	return
}
func (ok *OKX) WithdrawFee() (wf []q.NetworkWithdraw, err error) {
	type currenciesResponse struct {
		Code int    `json:"code,string"`
		Msg  string `json:"msg"`
		Data []struct {
			CanDep               bool    `json:"canDep"`
			CanInternal          bool    `json:"canInternal"`
			CanWd                bool    `json:"canWd"`
			Ccy                  string  `json:"ccy"`
			Chain                string  `json:"chain"`
			DepQuotaFixed        string  `json:"depQuotaFixed"`
			DepQuoteDailyLayer2  string  `json:"depQuoteDailyLayer2"`
			LogoLink             string  `json:"logoLink"`
			MainNet              bool    `json:"mainNet"`
			MaxFee               float64 `json:"maxFee,string"`
			MaxWd                float64 `json:"maxWd,string"`
			MinDep               string  `json:"minDep"`
			MinDepArrivalConfirm string  `json:"minDepArrivalConfirm"`
			MinFee               float64 `json:"minFee,string"`
			MinWd                float64 `json:"minWd,string"`
			MinWdUnlockConfirm   float64 `json:"minWdUnlockConfirm,string"`
			Name                 string  `json:"name"`
			NeedTag              bool    `json:"needTag"`
			UsedDepQuotaFixed    string  `json:"usedDepQuotaFixed"`
			UsedWdQuota          string  `json:"usedWdQuota"`
			WdQuota              string  `json:"wdQuota"`
			WdTickSz             string  `json:"wdTickSz"`
		} `json:"data"`
	}
	var response currenciesResponse
	urlPath := "/api/v5/asset/currencies"
	err = ok.DoAuthorRequest(http.MethodGet, urlPath, "", &response)
	if err != nil {
		zelo.OnErr(err).Send()
		err = nil
		return
	}
	if response.Code != 0 {
		return nil, fmt.Errorf("GetTickerV5 error:%s", response.Msg)

	}
	for _, v := range response.Data {
		c := q.NetworkWithdraw{
			Coin:           v.Name,
			DepositEnable:  v.CanDep,
			MinConfirm:     v.MinWdUnlockConfirm, // min number for balance confirmation
			Name:           v.Name,
			Network:        v.Chain,
			WithdrawEnable: v.CanWd,
			Fee:            v.MinFee, // `json:"withdrawFee"`
			Max:            v.MaxWd,  // `json:"withdrawMax"`
			Min:            v.MinWd,  // `json:"withdrawMin"`
			// WithdrawIntegerMultiple : `json:"withdrawIntegerMultiple"`
			// ExFee                   :
			// ExMin                   :
		}
		wf = append(wf, c)
	}
	return
}
func (ok *OKX) Fee() (f float64) {
	f = 0.001
	// var response any
	// var err error
	// err = ok.DoRequest("GET", "/api/spot/v3/trade_fee?instrument_id=BTC-USDT", "", &response)
	// if err != nil {
	// 	return
	// }
	// logs.I(response)
	// taker, o := response.(map[string]any)["taker"]
	// if !o {
	// 	logs.W("okx ", err, response.(map[string]any))
	// 	err = errors.New("no taker string")
	// 	return
	// }
	// f, err = strconv.ParseFloat(taker.(string), 64)
	// if err != nil {
	// 	logs.F(err, response.(map[string]any))
	// 	return
	// }
	return f
}
