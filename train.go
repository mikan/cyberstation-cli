package cyberstation

import "strings"

// Availability は予約可能状態を示します。
type Availability string

// TrainGroup は列車種別を示します。
type TrainGroup int

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
	// NozomiGroup は "のぞみ" 系 (東海道・山陽・九州新幹線) の列車種別を示します。
	NozomiGroup = TrainGroup(1)
	// KodamaGroup は "こだま" 系 (東海道・山陽新幹線) の列車種別を示します。
	KodamaGroup = TrainGroup(2)
	// HayabusaGroup は "はやぶさ" 系 (東北・北海道・山形・秋田新幹線) の列車種別を示します。
	HayabusaGroup = TrainGroup(3)
	// TokiGriup は "とき" 系 (上越・北陸新幹線) の列車種別を示します。
	TokiGroup = TrainGroup(4)
	// ZairaiGriup は在来線の列車種別を示します。
	ZairaiGroup = TrainGroup(5)
)

// Train は1列車あたりの予約状況を示します。
type Train struct {
	TrainName          string
	DepartureTime      string
	ArriveTime         string
	StandardNoSmoking  Availability
	StandardSmoking    Availability
	GreenNoSmoking     Availability
	GreenSmoking       Availability
	SleeperANoSmoking  Availability
	SleeperASmoking    Availability
	SleeperBNoSmoking  Availability
	SleeperBSmoking    Availability
	GranClassNoSmoking Availability
}

// IsAvailable はこの列車の任意の席タイプで空席があるか調べます。
func (t Train) IsAvailable() bool {
	all := t.StandardNoSmoking + t.StandardSmoking + t.GreenNoSmoking + t.GreenSmoking +
		t.SleeperANoSmoking + t.SleeperASmoking + t.SleeperBNoSmoking + t.SleeperBSmoking + t.GranClassNoSmoking
	if strings.Contains(string(all), string(Available)) {
		return true
	}
	if strings.Contains(string(all), string(LittleAvailable)) {
		return true
	}
	return false
}

// IsSleeperReservable はこの列車の寝台席が予約対象かどうか調べます。
func (t Train) IsSleeperReservable() bool {
	sleepers := t.SleeperANoSmoking + t.SleeperASmoking + t.SleeperBNoSmoking + t.SleeperBSmoking
	if strings.Contains(string(sleepers), string(Available)) {
		return true
	}
	if strings.Contains(string(sleepers), string(LittleAvailable)) {
		return true
	}
	if strings.Contains(string(sleepers), string(NotAvailable)) {
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
