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
	d     = flag.String("date", time.Now().Format(dateFormat), "日付 (例: 2019/1/12)")
	t     = flag.String("time", time.Now().Format(timeFormat), "時刻 (例: 12:30)")
	from  = flag.String("from", "", "出発駅 (例: 東京)")
	to    = flag.String("to", "", "到着駅 (例: 大垣)")
	group = flag.Int("group", 5, "1: のぞみ・みずほ等, 2: こだま, 3: はやぶさ等, 4: とき・かがやき等, 5: 在来線")
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
	trains, err := cyberstation.Vacancy(parsed, *from, *to, *group)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}
	if len(trains) == 0 {
		fmt.Println("エラー: 結果がありません")
		os.Exit(1)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	_, err = fmt.Fprintf(w, "列車名\t発時刻\t着時刻\t指🚭\t指🚬\tG🚭\tG🚬\tA寝🚭\tA寝🚬\tB寝🚭\tB寝🚬\n")
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}
	reservable := false
	for _, train := range trains {
		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			train.TrainName,
			train.DepartureTime, train.ArriveTime,
			train.StandardNoSmoking.Emoji(), train.StandardSmoking.Emoji(),
			train.GreenNoSmoking.Emoji(), train.GreenSmoking.Emoji(),
			train.SleeperANoSmoking.Emoji(), train.SleeperASmoking.Emoji(),
			train.SleeperBNoSmoking.Emoji(), train.SleeperBSmoking.Emoji(),
		); err != nil {
			fmt.Printf("エラー: %v\n", err)
			os.Exit(1)
		}
		if !reservable && train.IsReservable() {
			reservable = true
		}
	}
	if err := w.Flush(); err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}
	if reservable {
		fmt.Printf("%s %s %s▶%s 空席があります😃\n", *d, *t, *from, *to)
	} else {
		fmt.Printf("%s %s %s▶%s 満席です😢\n", *d, *t, *from, *to)
	}
}
