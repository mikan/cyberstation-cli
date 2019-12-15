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
	tg := cyberstation.TrainGroup(*group)
	trains, err := cyberstation.Vacancy(parsed, *from, *to, tg)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}
	if len(trains) == 0 {
		fmt.Println("エラー: 結果がありません")
		os.Exit(1)
	}
	printHumanReadable(trains, tg)
}

func printHumanReadable(trains []cyberstation.Train, group cyberstation.TrainGroup) {
	available := false
	sleeperReservable := false
	for _, train := range trains {
		if !available && train.IsAvailable() {
			available = true
		}
		if !sleeperReservable && train.IsSleeperReservable() {
			sleeperReservable = true
		}
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	if group == cyberstation.HayabusaGroup || group == cyberstation.TokiGroup {
		exitOnError(fmt.Fprintf(w, "指🚭\tG🚭\tGC🚭\t発時刻\t着時刻\t列車名\n"))
		for _, train := range trains {
			exitOnError(fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				train.StandardNoSmoking.Emoji(), train.GreenNoSmoking.Emoji(), train.GranClassNoSmoking.Emoji(),
				train.DepartureTime, train.ArriveTime, train.TrainName,
			))
		}
	} else if sleeperReservable {
		exitOnError(fmt.Fprintf(w, "指🚭\t指🚬\tG🚭\tG🚬\tA寝🚭\tA寝🚬\tB寝🚭\tB寝🚬\t発時刻\t着時刻\t列車名\n"))
		for _, train := range trains {
			exitOnError(fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				train.StandardNoSmoking.Emoji(), train.StandardSmoking.Emoji(),
				train.GreenNoSmoking.Emoji(), train.GreenSmoking.Emoji(),
				train.SleeperANoSmoking.Emoji(), train.SleeperASmoking.Emoji(),
				train.SleeperBNoSmoking.Emoji(), train.SleeperBSmoking.Emoji(),
				train.DepartureTime, train.ArriveTime, train.TrainName,
			))
		}
	} else {
		exitOnError(fmt.Fprintf(w, "指🚭\t指🚬\tG🚭\tG🚬\t発時刻\t着時刻\t列車名\n"))
		for _, train := range trains {
			exitOnError(fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				train.StandardNoSmoking.Emoji(), train.StandardSmoking.Emoji(),
				train.GreenNoSmoking.Emoji(), train.GreenSmoking.Emoji(),
				train.DepartureTime, train.ArriveTime, train.TrainName,
			))
		}
	}
	if err := w.Flush(); err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}
	if available {
		fmt.Printf("%s %s %s▶%s 空席があります😃\n", *d, *t, *from, *to)
	} else {
		fmt.Printf("%s %s %s▶%s 満席です😢\n", *d, *t, *from, *to)
	}
}

func exitOnError(_ int, err error) {
	if err != nil {
		fmt.Printf("エラー: %v", err)
		os.Exit(1)
	}
}
