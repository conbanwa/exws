package okx

import (
	"context"
	"github.com/amir-the-h/okex"
	"github.com/amir-the-h/okex/api"
	"github.com/amir-the-h/okex/events"
	"github.com/amir-the-h/okex/events/public"
	ws_public_requests "github.com/amir-the-h/okex/requests/ws/public"
	"os"
	"testing"
	"time"
)

func init() {
	log.Output(os.Stderr)
}

func TestWs(t *testing.T) {
	apiKey := "YOUR-API-KEY"
	secretKey := "YOUR-SECRET-KEY"
	passphrase := "YOUR-PASS-PHRASE"
	dest := okex.NormalServer // The main API server
	ctx := context.Background()
	client, err := api.NewClient(ctx, apiKey, secretKey, passphrase, dest)
	if err != nil {
		t.Error(err)
	}

	log.Println("Starting")
	errChan := make(chan *events.Error)
	subChan := make(chan *events.Subscribe)
	uSubChan := make(chan *events.Unsubscribe)
	logChan := make(chan *events.Login)
	sucChan := make(chan *events.Success)
	client.Ws.SetChannels(errChan, subChan, uSubChan, logChan, sucChan)

	obCh := make(chan *public.OrderBook)
	err = client.Ws.Public.OrderBook(ws_public_requests.OrderBook{
		InstID:  "BTC-USD-SWAP",
		Channel: "books",
	}, obCh)
	if err != nil {
		t.Error(err)
	}
	var listener = func() {
		for {
			select {
			case <-logChan:
				log.Print("[Authorized]")
			case success := <-sucChan:
				log.Printf("[SUCCESS]\t%+v", success)
			case sub := <-subChan:
				channel, _ := sub.Arg.Get("channel")
				log.Printf("[Subscribed]\t%s", channel)
			case uSub := <-uSubChan:
				channel, _ := uSub.Arg.Get("channel")
				log.Printf("[Unsubscribed]\t%s", channel)
			case err := <-client.Ws.ErrChan:
				log.Printf("[Error]\t%+v", err)
				for _, datum := range err.Data {
					log.Printf("[Error]\t\t%+v", datum)
				}
			case i := <-obCh:
				ch, _ := i.Arg.Get("channel")
				log.Printf("[Event]\t%s", ch)
				for _, p := range i.Books {
					for i := len(p.Asks) - 1; i >= 0; i-- {
						log.Printf("\t\tAsk\t%+v", p.Asks[i])
					}
					for _, bid := range p.Bids {
						log.Printf("\t\tBid\t%+v", bid)
					}
				}
			case b := <-client.Ws.DoneChan:
				log.Printf("[End]:\t%v", b)
				return
			}
		}
	}
	go listener()
	time.Sleep(time.Millisecond * 10)
}
