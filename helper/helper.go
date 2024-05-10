package helper

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/nyaruka/phonenumbers"
)

func IsPhoneNumber(phoneNumber string) bool {
	if phoneNumber[0] != '+' {	
		return false
	}
    // Parse the phone number
    parsedPhoneNumber, err := phonenumbers.Parse(phoneNumber, "")
    if err != nil {
        return false
    }

    // Check if the phone number is valid according to the library
    return phonenumbers.IsValidNumber(parsedPhoneNumber)
}
func Includes(target string, array []string)bool{
	for _, value := range array {
        if value == target {
            return true
        }
    }
    return false
}
func IsURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}
func FormatToIso860(s string)string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		
		return ""
	}

	// Format the time object into ISO 8601 format
	return t.Format("2006-01-02T15:04:05Z07:00")
}
func IsUUID(s string) bool{
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
    return uuidRegex.MatchString(s)
}
func RemoveSpaces(s string)string {
	return strings.ReplaceAll(s, " ","")
}
func ContainSpaces(s string)bool{
	return strings.Contains(s, " ")
}