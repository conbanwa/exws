package wstrader

import (
	"errors"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
	"github.com/conbanwa/wstrader/stat/zelo"
	"reflect"
	"time"
)

var log = zelo.Writer

/*
  - 本函数只适合，返回两个参数的API重试调用，其中一个参数必须是error
    @retry  重试次数
    @delay  每次重试延迟时间间隔
    @method 调用的函数，比如: api.GetTicker ,注意：不是api.GetTicker(...)
    @params 参数,顺序一定要按照实际调用函数入参顺序一样
    @return 返回
*/
func RE(retry int, delay time.Duration, method any, params ...any) any {
	invokeM := reflect.ValueOf(method)
	if invokeM.Kind() != reflect.Func {
		return errors.New("method not a function")
	}
	var value = make([]reflect.Value, len(params))
	var i = 0
	for ; i < len(params); i++ {
		value[i] = reflect.ValueOf(params[i])
	}
	var retV any
	var retryC = 0
_CALL:
	if retryC > 0 {
		log.Info().Msgf("sleep %v after re call", delay)
		time.Sleep(delay)
	}
	retValues := invokeM.Call(value)
	for _, v := range retValues {
		if v.Type().String() != "error" {
			retV = v.Interface()
			continue
		}
		if v.IsNil() {
			continue
		}
		log.Error().Any("invokeM", v).Msg("[api error]")
		retryC++
		if retryC <= retry-1 {
			log.Info().Msgf("Invoke Method[%s] Error , Begin Retry Call [%d] ...", invokeM.String(), retryC)
			goto _CALL
		}
		log.Error().Msg("Invoke Method Fail ???" + invokeM.String())
		return v.Interface()
	}
	return retV
}

/**
 * call all unfinished orders
 */
func CancelAllUnfinishedOrders(api API, currencyPair cons.CurrencyPair) int {
	if api == nil {
		log.Error().Msg("api instance is nil ??? , please new a api instance")
		return -1
	}
	c := 0
	for {
		ret := RE(2, 200*time.Millisecond, api.GetUnfinishedOrders, currencyPair)
		if err, isok := ret.(error); isok {
			log.Error().Err(err).Msg("[api error]")
			break
		}
		if ret == nil {
			break
		}
		orders, isok := ret.([]q.Order)
		if !isok || len(orders) == 0 {
			break
		}
		for _, ord := range orders {
			_, err := api.CancelOrder(ord.OrderID2, currencyPair)
			if err != nil {
				log.Error().Err(err).Send()
			} else {
				c++
			}
			time.Sleep(120 * time.Millisecond) //控制频率
		}
	}
	return c
}

/**
 * call all unfinished future orders
 * @return c 成功撤单数量
 */
func CancelAllUnfinishedFutureOrders(api FutureRestAPI, contractType string, currencyPair cons.CurrencyPair) int {
	if api == nil {
		log.Error().Msg("api instance is nil ??? , please new a api instance")
		return 0
	}
	c := 0
	for {
		ret := RE(10, 200*time.Millisecond, api.GetUnfinishFutureOrders, currencyPair, contractType)
		if err, isOk := ret.(error); isOk {
			log.Error().Err(err).Msg("[api error]")
			break
		}
		if ret == nil {
			break
		}
		orders, isOk := ret.([]FutureOrder)
		if !isOk || len(orders) == 0 {
			break
		}
		for _, ord := range orders {
			_, err := api.FutureCancelOrder(currencyPair, contractType, ord.OrderID2)
			if err != nil {
				log.Error().Err(err).Send()
			} else {
				c++
			}
			time.Sleep(120 * time.Millisecond) //控制频率
		}
	}
	return c
}
