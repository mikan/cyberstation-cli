package cyberstation

import "strings"

const (
	// Available 予約可能であることを示します。
	Available = "○"
	// LittleAvailable 残席が少ないことを示します。
	LittleAvailable = "△"
	// NotAvailable 満席であることを示します。
	NotAvailable = "×"
	// NotReservable 予約対象外であることを示します。
	NotReservable = "-"
	// UseSmokingRoom 座席は禁煙ですが、喫煙ルームがあることを示します。
	UseSmokingRoom = "＊"
)

// Train 1列車あたりの予約状況を示します。
type Train struct {
	TrainName         string
	DepartureTime     string
	ArriveTime        string
	StandardNoSmoking string
	StandardSmoking   string
	GreenNoSmoking    string
	GreenSmoking      string
	SleeperANoSmoking string
	SleeperASmoking   string
	SleeperBNoSmoking string
	SleeperBSmoking   string
}

// IsReservable この列車が予約可能かどうか調べます。
func (t Train) IsReservable() bool {
	all := t.StandardNoSmoking + t.StandardSmoking + t.GreenNoSmoking + t.GreenSmoking +
		t.SleeperANoSmoking + t.SleeperASmoking + t.SleeperBNoSmoking + t.SleeperBSmoking
	if strings.Contains(all, Available) {
		return true
	}
	if strings.Contains(all, LittleAvailable) {
		return true
	}
	return false
}
