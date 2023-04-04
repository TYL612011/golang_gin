package model


type User struct {
	Username		string `gorm:"primary_key" json:"username"`
	Password		string `json:"password"`
	Email			string `json:"email"`
}
