package ftx

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/conbanwa/logs"
)

const URL = "https://ftx.com/api/"

func (client *Client) signRequest(method string, path string, body []byte) *http.Request {
	ts := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	signaturePayload := ts + method + "/api/" + path + string(body)
	signature := client.sign(signaturePayload)
	req, _ := http.NewRequest(method, URL+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("FTX-KEY", client.Api)
	req.Header.Set("FTX-SIGN", signature)
	req.Header.Set("FTX-TS", ts)
	if client.Subaccount != "" {
		req.Header.Set("FTX-SUBACCOUNT", client.Subaccount)
	}
	return req
}
func (client *Client) _get(path string, body []byte) (*http.Response, error) {
	preparedRequest := client.signRequest("GET", path, body)
	resp, err := client.Client.Do(preparedRequest)
	return resp, err
}
func (client *Client) _post(path string, body []byte) (*http.Response, error) {
	preparedRequest := client.signRequest("POST", path, body)
	resp, err := client.Client.Do(preparedRequest)
	return resp, err
}
func (client *Client) _delete(path string, body []byte) (*http.Response, error) {
	preparedRequest := client.signRequest("DELETE", path, body)
	resp, err := client.Client.Do(preparedRequest)
	return resp, err
}
func _processResponse(resp *http.Response, result any) error {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.E("Error processing response:", err)
		return err
	}
	err = json.Unmarshal(body, result)
	if err != nil {
		logs.E("Error processing response:", err)
		return err
	}
	return nil
}
