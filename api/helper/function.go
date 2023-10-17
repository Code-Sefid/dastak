package helper

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
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


func Separate(Number int) string {
	numberStr := strconv.Itoa(Number)

	parts := strings.Split(numberStr, ".")

	wholePart := parts[0]
	wholePartWithCommas := addCommas(wholePart)

	if len(parts) > 1 {
		return wholePartWithCommas + "." + parts[1]
	}
	return wholePartWithCommas
}

func addCommas(input string) string {
	var result string
	n := len(input)
	for i, char := range input {
		result += string(char)
		if (n-i-1)%3 == 0 && i != n-1 {
			result += ","
		}
	}
	return result
}
