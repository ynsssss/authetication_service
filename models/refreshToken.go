package models

type RefreshToken struct {
	Value     string
	PairToken string
	UserId    string
	ExpiresId int64 //timestamp
}

type RefreshTokenWithHash struct {
	Hash      []byte
	PairToken string
	UserId    string
	ExpiresId int64 //timestamp
}
