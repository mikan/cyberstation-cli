package cyberstation

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const vacancyURL = "http://www1.jr.cyberstation.ne.jp/csws/Vacancy.do"

// Vacancy は空席情報を照会します。
func Vacancy(target time.Time, departure, arrival string, group int) ([]Train, error) {
	form := fmt.Sprintf("script=1&month=%d&day=%d&hour=%d&minute=%d&train=%d&dep_stn=%s&arr_stn=%s",
		target.Month(), target.Day(), target.Hour(), target.Minute(), group, shiftJIS(departure), shiftJIS(arrival))
	req, err := http.NewRequest(http.MethodPost, vacancyURL, strings.NewReader(form))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", vacancyURL)
	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, fmt.Errorf("リクエストを送信できませんでした: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	defer SafeClose(resp.Body, "response")
	return parseTable(bufio.NewScanner(transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())))
}

func shiftJIS(name string) string {
	w := bytes.Buffer{}
	tw := transform.NewWriter(&w, japanese.ShiftJIS.NewEncoder())
	defer SafeClose(tw, "shift_jis encoder")
	if _, err := tw.Write([]byte(name)); err != nil {
		return name
	}
	return w.String()
}

func parseTable(scanner *bufio.Scanner) ([]Train, error) {
	var records []Train
	tableSeeking := true
	tdSeeking := false
	var current *Train
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "受け付け時間外") {
			return nil, errors.New("受付時間外です")
		}
		if strings.Contains(line, "ご希望の乗車日の空席状況は照会できません") {
			return nil, errors.New("照会できない日時です")
		}
		if strings.Contains(line, "入力された項目に誤りがあります") {
			return nil, errors.New("照会できない駅名です")
		}
		if strings.Contains(line, "ご希望の情報はお取り扱いできません") {
			return nil, nil
		}
		if tableSeeking {
			tableSeeking = !strings.Contains(line, "列車名")
			continue
		}
		if tdSeeking {
			if strings.Contains(line, "<TR>") {
				break
			}
			if !strings.Contains(line, "<td align=\"left\">") {
				continue
			}
			tdSeeking = false
			trainName := strings.Replace(strings.Replace(line, "<td align=\"left\">", "", 1), "</td>", "", 1)
			current = &Train{TrainName: trainName}
			continue
		}
		if strings.Contains(line, "</tr>") {
			tdSeeking = true
			if current != nil {
				records = append(records, *current)
				current = nil
			}
			continue
		}
		if !strings.Contains(line, "<td align=\"center\">") {
			continue
		}
		content := strings.Replace(strings.Replace(line, "<td align=\"center\">", "", 1), "</td>", "", 1)
		if current == nil {
			continue
		}
		if len(current.DepartureTime) == 0 {
			current.DepartureTime = content
		} else if len(current.ArriveTime) == 0 {
			current.ArriveTime = content
		} else if len(current.StandardNoSmoking) == 0 {
			current.StandardNoSmoking = Availability(content)
		} else if len(current.StandardSmoking) == 0 {
			current.StandardSmoking = Availability(content)
		} else if len(current.GreenNoSmoking) == 0 {
			current.GreenNoSmoking = Availability(content)
		} else if len(current.GreenSmoking) == 0 {
			current.GreenSmoking = Availability(content)
		} else if len(current.SleeperANoSmoking) == 0 {
			current.SleeperANoSmoking = Availability(content)
		} else if len(current.SleeperASmoking) == 0 {
			current.SleeperASmoking = Availability(content)
		} else if len(current.SleeperBNoSmoking) == 0 {
			current.SleeperBNoSmoking = Availability(content)
		} else if len(current.SleeperBSmoking) == 0 {
			current.SleeperBSmoking = Availability(content)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("照会結果の解析に失敗しました: %w", err)
	}
	return records, nil
}
