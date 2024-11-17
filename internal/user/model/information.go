package model

// embedded
type Information struct {
	FirstName   string `gorm:"column:first_name"`
	LastName    string `gorm:"column:last_name"`
	City        string `gorm:"column:city"`
	PhoneNumber string `gorm:"column:phone_number"`
}
