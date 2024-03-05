package models

type AccessToken struct {
	Value     string
	UserId    string
	PairToken string
	ExpiresIn int64
}
