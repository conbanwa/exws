package binance

func (bn *Binance) header() map[string]string {
	return map[string]string{"X-MBX-APIKEY": bn.accessKey}
}
