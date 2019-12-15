package cyberstation

import "strings"

// Availability は予約可能状態を示します。
type Availability string

const (
	// Available は予約可能であることを示します。
	Available = Availability("○")
	// LittleAvailable は残席が少ないことを示します。
	LittleAvailable = Availability("△")
	// NotAvailable は満席であることを示します。
	NotAvailable = Availability("×")
	// NotReservable は予約対象外であることを示します。
	NotReservable = Availability("-")
	// UseSmokingRoom は座席は禁煙ですが、喫煙ルームがあることを示します。
	UseSmokingRoom = Availability("＊")
)

// Train は1列車あたりの予約状況を示します。
type Train struct {
	TrainName         string
	DepartureTime     string
	ArriveTime        string
	StandardNoSmoking Availability
	StandardSmoking   Availability
	GreenNoSmoking    Availability
	GreenSmoking      Availability
	SleeperANoSmoking Availability
	SleeperASmoking   Availability
	SleeperBNoSmoking Availability
	SleeperBSmoking   Availability
}

// IsReservable はこの列車が予約可能かどうか調べます。
func (t Train) IsReservable() bool {
	all := t.StandardNoSmoking + t.StandardSmoking + t.GreenNoSmoking + t.GreenSmoking +
		t.SleeperANoSmoking + t.SleeperASmoking + t.SleeperBNoSmoking + t.SleeperBSmoking
	if strings.Contains(string(all), string(Available)) {
		return true
	}
	if strings.Contains(string(all), string(LittleAvailable)) {
		return true
	}
	return false
}

// Emoji は絵文字に変換します。
func (a Availability) Emoji() string {
	switch a {
	case NotReservable:
		return "⚪"
	case UseSmokingRoom:
		return "⚫"
	case NotAvailable:
		return "🔴"
	case LittleAvailable:
		return "🔺"
	case Available:
		return "🔵"
	default:
		return string(a)
	}
}
