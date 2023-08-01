package utils

import "time"

func FormatTime(t string) (string, error) {
	trackTime, err := time.Parse("2006-01-02T15:04:05", t)
	if err != nil {
		return "", err
	}

	now := time.Now()
	if now.Year() == trackTime.Year() && now.YearDay() == trackTime.YearDay() {
		return "Сегодня в " + trackTime.Format("15:04"), nil
	}

	if now.Year() == trackTime.Year() && now.YearDay()-1 == trackTime.YearDay() {
		return "Вчера в " + trackTime.Format("15:04"), nil
	}

	return trackTime.Format("02.01.2006 в 15:04"), nil
}
