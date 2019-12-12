package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mikan/cyberstation-cli"
)

type currentState int

const (
	dateFormat         = "2006/1/2"
	timeFormat         = "15:4"
	offlineState       = currentState(0)
	reservableState    = currentState(1)
	notReservableState = currentState(2)
)

var (
	webhook  = flag.String("webhook", "", "webhook URL")
	d        = flag.String("date", time.Now().Format(dateFormat), "date")
	t        = flag.String("time", time.Now().Format(timeFormat), "time")
	from     = flag.String("from", "", "departure station name")
	to       = flag.String("to", "", "arrival station name")
	interval = flag.Int("interval", 1, "checking interval in minutes")
)

func main() {
	flag.Parse()
	if len(*from) == 0 || len(*to) == 0 || len(*webhook) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	parsed, err := time.Parse(dateFormat+" "+timeFormat, *d+" "+*t)
	if err != nil {
		flag.Usage()
		os.Exit(2)
	}
	post(fmt.Sprintf("%s %s %s▶%s の監視を開始します...", *d, *t, *from, *to))
	currentState := offlineState
	for {
		trains, err := cyberstation.Vacancy(parsed, *from, *to)
		if err != nil || len(trains) == 0 {
			currentState = offlineState
		} else {
			reservable := false
			trainName := ""
			for _, train := range trains {
				if train.IsReservable() {
					reservable = true
					trainName = train.TrainName
					break
				}
			}
			if reservable && currentState != reservableState {
				currentState = reservableState
				post(trainName + "が予約可能になりました!")
			} else if !reservable && (currentState != notReservableState && currentState != offlineState) {
				currentState = notReservableState
				post("満席になりました...")
			}
		}
		time.Sleep(time.Duration(*interval) * time.Minute)
	}
}

func post(msg string) {
	log.Printf("sending webhook message: %s", msg)
	body, err := json.Marshal(struct {
		Text string `json:"text"`
	}{Text: msg})
	if err != nil {
		log.Fatalf("failed to marshal message: %v", err)
	}
	resp, err := http.Post(*webhook, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("failed to post payload: %v", err)
		return
	}
	if resp.StatusCode == http.StatusOK {
		return
	}
	defer cyberstation.SafeClose(resp.Body, "webhook response body")
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response: %v", err)
		return
	}
	log.Printf("failed to post webhook: %s", string(content))
}
