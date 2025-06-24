package web

import (
	"encoding/json"
	"github.com/conbanwa/exws/stat/zero"
	"github.com/conbanwa/slice"
	"net/http"
	"net/url"
	"regexp"
)

var log = zero.Writer

func HttpGet(client *http.Client, reqUrl string) (map[string]any, error) {
	respData, err := NewRequest(client, "GET", reqUrl, "", nil)
	if err != nil {
		return nil, err
	}
	var bodyDataMap map[string]any
	err = json.Unmarshal(respData, &bodyDataMap)
	if err != nil {
		log.Error().Bytes("response data", respData).Send()
		return nil, err
	}
	return bodyDataMap, nil
}
func HttpGet2(client *http.Client, reqUrl string, headers map[string]string) (map[string]any, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	respData, err := NewRequest(client, "GET", reqUrl, "", headers)
	if err != nil {
		return nil, err
	}
	var bodyDataMap map[string]any
	err = json.Unmarshal(respData, &bodyDataMap)
	if err != nil {
		log.Error().Bytes("response data", respData).Send()
		return nil, err
	}
	return bodyDataMap, nil
}
func HttpGet3(client *http.Client, reqUrl string, headers map[string]string) ([]any, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	respData, err := NewRequest(client, "GET", reqUrl, "", headers)
	if err != nil {
		return nil, err
	}
	var bodyDataMap []any
	err = json.Unmarshal(respData, &bodyDataMap)
	if err != nil {
		log.Error().Int("len", len(respData)).Bytes("response data", respData).Str("reqUrl", reqUrl).Send()
		return nil, err
	}
	return bodyDataMap, nil
}
func HttpGet4(client *http.Client, reqUrl string, headers map[string]string, result any, reg ...string) error {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	respData, err := NewRequest(client, "GET", reqUrl, "", headers)
	if err != nil {
		return err
	}
	for _, r := range reg {
		reg := regexp.MustCompile(r)
		respData = reg.ReplaceAll(respData, []byte(`$1"0"`))
	}
	err = json.Unmarshal(respData, result)
	if err != nil {
		log.Error().Err(err).Int("len", len(respData)).Bytes("response data", respData).Str("reqUrl", reqUrl).Msg("HttpGet4 - json.Unmarshal failed")
		return err
	}
	return nil
}
func HttpGet5(client *http.Client, reqUrl string, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	respData, err := NewRequest(client, "GET", reqUrl, "", headers)
	if err != nil {
		return nil, err
	}
	return respData, nil
}
func HttpZB(client *http.Client, reqUrl string) (bodyDataMap map[string]map[string]string, err error) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	respData, err := NewRequest(client, "GET", reqUrl, "", headers)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respData, &bodyDataMap)
	if err != nil {
		log.Error().Bytes("response data", respData).Send()
		return nil, err
	}
	return bodyDataMap, nil
}
func HttpZBP(client *http.Client, reqUrl string) (bodyDataMap map[string]map[string]float64, err error) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	respData, err := NewRequest(client, "GET", reqUrl, "", headers)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respData, &bodyDataMap)
	if err != nil {
		log.Error().Bytes("response data", respData).Send()
		return nil, err
	}
	return bodyDataMap, nil
}
func HttpPostForm(client *http.Client, reqUrl string, postData url.Values) ([]byte, error) {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded"}
	return NewRequest(client, "POST", reqUrl, postData.Encode(), headers)
}
func HttpPostForm2(client *http.Client, reqUrl string, postData url.Values, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return NewRequest(client, "POST", reqUrl, postData.Encode(), headers)
}
func HttpPostForm3(client *http.Client, reqUrl string, postData string, headers map[string]string) ([]byte, error) {
	return NewRequest(client, "POST", reqUrl, postData, headers)
}
func HttpPostForm4(client *http.Client, reqUrl string, postData map[string]string, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/json"
	data, _ := json.Marshal(postData)
	return NewRequest(client, "POST", reqUrl, slice.Bytes2String(data), headers)
}
func HttpDeleteForm(client *http.Client, reqUrl string, postData url.Values, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return NewRequest(client, "DELETE", reqUrl, postData.Encode(), headers)
}
func HttpPut(client *http.Client, reqUrl string, postData url.Values, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return NewRequest(client, "PUT", reqUrl, postData.Encode(), headers)
}
