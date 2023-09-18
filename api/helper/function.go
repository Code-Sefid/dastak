package helper

import (
	"regexp"
	"time"
)

func IsEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, err := regexp.MatchString(pattern, email)
	if err != nil {
		return false
	}
	return match
}

func DateToWeekday(date time.Time) int {
	weekday := int(date.Weekday())
	return (weekday + 6) % 7
}
