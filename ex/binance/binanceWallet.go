package binance

import (
	"encoding/json"
	"fmt"
	"github.com/conbanwa/num"
	"net/url"
	. "github.com/conbanwa/wstrader"
	. "github.com/conbanwa/wstrader/cons"
	. "github.com/conbanwa/wstrader/q"
	"github.com/conbanwa/wstrader/stat/zelo"
	. "github.com/conbanwa/wstrader/web"
	"strings"
	"sync"
)

func (bn *Binance) Balances() (available, frozen *sync.Map, err error) {
	available, frozen = new(sync.Map), new(sync.Map)
	bn.AssetDetail()
	params := url.Values{}
	bn.buildParamsSigned(&params)
	path := bn.apiV3 + AccountUri + params.Encode()
	resp, err := HttpGet2(bn.httpClient, path, bn.header())
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	if _, ok := resp["code"]; ok {
		err = fmt.Errorf(resp["msg"].(string))
		return
	}
	acc := Account{}
	acc.Exchange = bn.String()
	acc.SubAccounts = make(map[Currency]SubAccount)
	balances := resp["balances"].([]any)
	for _, v := range balances {
		vm := v.(map[string]any)
		if asset, ok := vm["asset"].(string); ok && !strings.HasPrefix(asset, "LD") {
			if free := num.ToFloat64(vm["free"]); free != 0 {
				available.Store(asset, free)
			}
			if locked := num.ToFloat64(vm["locked"]); locked != 0 {
				frozen.Store(asset, locked)
			}
		}
	}
	return
}
func (bn *Binance) TradeFee() (map[string]TradeFee, error) {
	postForm := url.Values{}
	bn.buildParamsSigned(&postForm)
	resp, err := HttpGet5(bn.httpClient, FeeUrl+"?"+postForm.Encode(), bn.header())
	if err != nil {
		return nil, err
	}
	var fees []any
	if err = json.Unmarshal(resp, &fees); err != nil {
		return nil, err
	}
	fm := make(map[string]TradeFee)
	for _, v := range fees {
		fee := v.(map[string]any)
		if symbol := fee["symbol"].(string); symbol != "" {
			TakerFee := num.ToFloat64(fee["takerCommission"])
			fm[symbol] = TradeFee{
				MakerFee: num.ToFloat64(fee["makerCommission"]),
				TakerFee: TakerFee}
			// log.Printf("response body: %+v", fee)
		}
	}
	return fm, nil
}

func (bn *Binance) AssetDetail() (sf map[string]AssetDetail, err error) {
	sf = make(map[string]AssetDetail)
	postForm := url.Values{}
	bn.buildParamsSigned(&postForm)
	resp, err := HttpGet5(bn.httpClient, NetworkWithdrawUrl+"?"+postForm.Encode(), bn.header())
	if err != nil {
		return nil, err
	}

	var fees any
	if err = json.Unmarshal(resp, &fees); err != nil {
		return nil, err
	}
	for _, v := range fees.([]any) {
		if fee := v.(map[string]any); fee["coin"] != nil {
			sf[fee["coin"].(string)] = AssetDetail{
				Coin:             v.(map[string]any)["coin"].(string),
				DepositAllEnable: v.(map[string]any)["depositAllEnable"].(bool),
				Free:             v.(map[string]any)["free"].(string),
				Freeze:           v.(map[string]any)["freeze"].(string),
				Ipoable:          v.(map[string]any)["ipoable"].(string),
				Ipoing:           v.(map[string]any)["ipoing"].(string),
				IsLegalMoney:     v.(map[string]any)["isLegalMoney"].(bool),
				Locked:           v.(map[string]any)["locked"].(string),
				Name:             v.(map[string]any)["name"].(string),
			}
		}
	}
	return sf, nil
}
func (bn *Binance) WithdrawFee() (sf []NetworkWithdraw, err error) {
	postForm := url.Values{}
	bn.buildParamsSigned(&postForm)
	resp, err := HttpGet5(bn.httpClient, NetworkWithdrawUrl+"?"+postForm.Encode(), bn.header())
	if err != nil {
		return nil, err
	}
	var fees []any
	if err = json.Unmarshal(resp, &fees); err != nil {
		return nil, err
	}
	exFee, err := bn.ExWithdrawFee()
	zelo.PanicOnErr(err).Send()
	for _, v := range fees {
		fee := v.(map[string]any)
		if fee["coin"] == nil {
			continue
		}
		networks := fee["networkList"].([]any)
		for _, v := range networks {
			network := v.(map[string]any)
			coi, ti := network["coin"], network["specialTips"]
			if coi == nil {
				continue
			}
			tip := ""
			if ti != nil {
				tip = ti.(string)
			}
			coin := coi.(string)
			withdrawFee := num.ToFloat64(network["withdrawFee"].(string))
			withdrawEnable := network["withdrawEnable"].(bool)
			if withdrawFee+exFee[coin].WithdrawFee == 0 && withdrawEnable {
				log.Debug().Msg(coin + "(network " + network["network"].(string) + ") has 0 withdrawFee. tip: " + tip)
			}
			sf = append(sf, NetworkWithdraw{
				SpecialTips:             tip,
				AddressRegex:            network["addressRegex"].(string),
				Coin:                    coin,
				DepositDesc:             network["depositDesc"].(string),
				DepositEnable:           network["depositEnable"].(bool),
				IsDefault:               network["isDefault"].(bool),
				MemoRegex:               network["memoRegex"].(string),
				MinConfirm:              network["minConfirm"].(float64),
				Name:                    network["name"].(string),
				Network:                 network["network"].(string),
				ResetAddressStatus:      network["resetAddressStatus"].(bool),
				UnLockConfirm:           network["unLockConfirm"].(float64),
				WithdrawDesc:            network["withdrawDesc"].(string),
				WithdrawEnable:          withdrawEnable,
				WithdrawIntegerMultiple: num.ToFloat64(network["withdrawIntegerMultiple"].(string)),
				SameAddress:             network["sameAddress"].(bool),
				EstimatedArrivalTime:    network["estimatedArrivalTime"].(float64),
				Busy:                    network["busy"].(bool),
				ExFee:                   exFee[coin].WithdrawFee,
				ExMin:                   exFee[coin].MinWithdrawAmount,
				Fee:                     withdrawFee,
				Min:                     num.ToFloat64(network["withdrawMin"].(string)),
				Max:                     num.ToFloat64(network["withdrawMax"].(string)),
			})
		}
	}
	//for _, v := range sf {
	//	if v.Coin == "XRP" {
	//		fmt.Printf("%+v\n", v)
	//	}
	//}
	return sf, nil
}
func (bn *Binance) ExWithdrawFee() (sf map[string]ExWithdraw, err error) {
	sf = make(map[string]ExWithdraw)
	postForm := url.Values{}
	bn.buildParamsSigned(&postForm)
	resp, err := HttpGet5(bn.httpClient, ExWithdrawUrl+"?"+postForm.Encode(), bn.header())
	if err != nil {
		return nil, err
	}
	var fees any
	if err = json.Unmarshal(resp, &fees); err != nil {
		return nil, err
	}
	for k, v := range fees.(map[string]any) {
		fee := v.(map[string]any)
		var withdraw = ExWithdraw{
			//Coin:             fee["coin"].(string),
			MinWithdrawAmount: num.ToFloat64(fee["minWithdrawAmount"].(string)),
			DepositStatus:     fee["depositStatus"].(bool),
			WithdrawFee:       num.ToFloat64(fee["withdrawFee"].(string)),
			WithdrawStatus:    fee["withdrawStatus"].(bool),
		}
		if fee["depositTip"] != nil {
			withdraw.DepositTip = fee["depositTip"].(string)
		}
		sf[k] = withdraw
	}
	return sf, nil
}

type AssetDetail struct {
	Coin             string            `json:"coin"`
	DepositAllEnable bool              `json:"depositAllEnable"`
	Free             string            `json:"free"`
	Freeze           string            `json:"freeze"`
	Ipoable          string            `json:"ipoable"`
	Ipoing           string            `json:"ipoing"`
	IsLegalMoney     bool              `json:"isLegalMoney"`
	Locked           string            `json:"locked"`
	Name             string            `json:"name"`
	NetworkList      []NetworkWithdraw `json:"networkList"`
}

type ExWithdraw struct {
	MinWithdrawAmount float64 `json:"minWithdrawAmount"`
	DepositStatus     bool    `json:"depositStatus"`
	WithdrawFee       float64 `json:"withdrawFee"`
	WithdrawStatus    bool    `json:"withdrawStatus"`
	DepositTip        string  `json:"depositTip"`
}
