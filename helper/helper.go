package helper

import (
	"Week2/forms"
	"regexp"
	"strings"
	"time"
)

func IsPhoneNumber(phoneNumber string) bool {
	if len(phoneNumber) == 0{
		return false
	}
	if phoneNumber[0] != '+' {	
		return false
	}

	return true
    // Parse the phone number
    // parsedPhoneNumber, err := phonenumbers.Parse(phoneNumber, "")
	// fmt.Println(parsedPhoneNumber)
    // if err != nil {
	// 	fmt.Println(err.Error())
    //     return false
    // }

    // // Check if the phone number is valid according to the library
    // return phonenumbers.IsValidNumber(parsedPhoneNumber)
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
	// fmt.Println(s)
	// u, err := url.ParseRequestURI(s)
	// fmt.Println(u.Scheme)
	// return err == nil && u.Scheme != "" && u.Host != ""
	regex := `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,4}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`
	pattern := regexp.MustCompile(regex)
	return pattern.MatchString(s)
	// return govalidator.IsURL(s)
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

// func SanitizePhoneNumber(s string) string{
// 	s = strings
// }
func IsProductsUnique(items []forms.ProductDetail) bool {
	seen := make(map[string]bool) // Map to store IDs already seen

	for _, item := range items {
		if seen[item.ProductId] {
			// ID is not unique, return false
			return false
		}
		// Mark the ID as seen
		seen[item.ProductId] = true
	}

	// All IDs are unique
	return true
}