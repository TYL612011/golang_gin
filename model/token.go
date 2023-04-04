package model

type AccessToken struct {
	Token	string	`json:"token"`
	Exp		int64	`json:"exp"`
}

type RefreshToken struct {
	Token	string	`json:"token"`
}

