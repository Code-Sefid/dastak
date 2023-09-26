package helper

import (
	"math/rand"
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

func GenerateFactorCode() string {
	result := ""
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charactersLength := len(characters)

	for counter := 0; counter < 5; counter++ {
		randomIndex := rand.Intn(charactersLength)
		result += string(characters[randomIndex])
	}

	return result
}
