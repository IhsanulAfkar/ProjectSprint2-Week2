package forms

type StaffRegister struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	Password    string `json:"password"`
}

type StaffLogin struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}