package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/mikan/cyberstation-cli"
)

const (
	dateFormat = "2006/1/2"
	timeFormat = "15:4"
)

var (
	d    = flag.String("date", time.Now().Format(dateFormat), "date")
	t    = flag.String("time", time.Now().Format(timeFormat), "time")
	from = flag.String("from", "", "departure station name")
	to   = flag.String("to", "", "arrival station name")
)

func main() {
	flag.Parse()
	if len(*from) == 0 || len(*to) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	parsed, err := time.Parse(dateFormat+" "+timeFormat, *d+" "+*t)
	if err != nil {
		flag.Usage()
		os.Exit(2)
	}
	trains, err := cyberstation.Query(parsed, *from, *to)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	}
	if len(trains) == 0 {
		fmt.Println("エラー: 結果がありません")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	_, err = fmt.Fprintf(w, "列車名\t発時刻\t着時刻\t指🚭\t指🚬\tG🚭\tG🚬\tA寝🚭\tA寝🚬\tB寝🚭\tB寝🚬\n")
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	}
	reservable := false
	for _, train := range trains {
		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			train.TrainName,
			train.DepartureTime, train.ArriveTime,
			emoji(train.StandardNoSmoking), emoji(train.StandardSmoking),
			emoji(train.GreenNoSmoking), emoji(train.GreenSmoking),
			emoji(train.SleeperANoSmoking), emoji(train.SleeperASmoking),
			emoji(train.SleeperBNoSmoking), emoji(train.SleeperBSmoking),
		); err != nil {
			fmt.Printf("エラー: %v\n", err)
		}
		if !reservable && train.IsReservable() {
			reservable = true
		}
	}
	if err := w.Flush(); err != nil {
		fmt.Printf("エラー: %v\n", err)
	}
	if reservable {
		fmt.Printf("%s %s %s▶%s 空席があります😃\n", *d, *t, *from, *to)
	} else {
		fmt.Printf("%s %s %s▶%s 満席です😢\n", *d, *t, *from, *to)
	}
}

func emoji(symbol string) string {
	switch symbol {
	case cyberstation.NotReservable:
		return "⚪"
	case cyberstation.UseSmokingRoom:
		return "⚫"
	case cyberstation.NotAvailable:
		return "🔴"
	case cyberstation.LittleAvailable:
		return "🔺"
	case cyberstation.Available:
		return "🔵"
	default:
		return symbol
	}
}
