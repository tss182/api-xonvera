package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"regexp"
	"time"
)

var (
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func Ternary[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

func GetTimeDuration(durationType string, duration int) time.Duration {
	switch durationType {
	case "minute":
		return time.Minute * time.Duration(duration)
	case "hour":
		return time.Hour * time.Duration(duration)
	case "day":
		return time.Hour * 24 * time.Duration(duration)
	default:
		return 0
	}
}

func HashSha256(plaintext string) string {
	hash := sha256.New()
	hash.Write([]byte(plaintext))
	return hex.EncodeToString(hash.Sum(nil))
}

func VerifyHashSha256(plaintext string, hash string) bool {
	return HashSha256(plaintext) == hash
}

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func PhoneNumberFormat(phoneNumber string) string {
	reg := regexp.MustCompile("[^0-9]+")
	//just show number only
	phoneNumber = reg.ReplaceAllString(phoneNumber, "")
	//if number starts with 08 or 8, it will replace 628
	reg = regexp.MustCompile("^(08|8)")
	phoneNumber = reg.ReplaceAllString(phoneNumber, "628")
	return phoneNumber
}

func InArray[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func DecimalSeparator(i int) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", i)
}
