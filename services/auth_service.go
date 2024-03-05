package services

import (
	"authentication_service/internal/jwt"
	"authentication_service/models"
	"authentication_service/repositories"
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	GenerateRefreshToken(userId, pairToken string) (models.RefreshToken, error)
	GenerateAccessToken(userId, pairToken string) models.AccessToken
	GetRefreshTokenForAccessToken(token string) (models.RefreshTokenWithHash, error)
	DeleteRefreshToken(pairToken string)
}

type authService struct {
	refreshRepo repositories.RefreshTokenRepository
	secret      string
}

func NewAuthService(refreshRepo repositories.RefreshTokenRepository, secret string) AuthService {
	return authService{
		refreshRepo: refreshRepo,
		secret:      secret,
	}
}

func (s authService) GenerateAccessToken(userId, pairToken string) models.AccessToken {
	exp := time.Now().Add(jwt.TokenMaxLifetimeDuration).Unix()
	jwtToken := jwt.GenerateJWT(userId, exp, s.secret)
	accessToken := models.AccessToken{
		UserId:    userId,
		PairToken: pairToken,
		Value:     jwtToken,
		ExpiresIn: exp,
	}
	return accessToken
}

func (s authService) GenerateRefreshToken(userId, pairToken string) (models.RefreshToken, error) {
	exp := time.Now().Add(time.Hour * 72).Unix()
	jwtToken := jwt.GenerateJWT(userId, exp, s.secret)
	refreshToken := models.RefreshToken{
		UserId:    userId,
		PairToken: pairToken,
		Value:     jwtToken,
		ExpiresId: exp,
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(jwtToken), bcrypt.DefaultCost)
	if err != nil {
		return models.RefreshToken{}, nil
	}

	refreshTokenWithHash := models.RefreshTokenWithHash{
		UserId:    userId,
		PairToken: pairToken,
		Hash:      hash,
		ExpiresId: exp,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.refreshRepo.SaveToken(ctx, refreshTokenWithHash)

	return refreshToken, err
}

func (s authService) GetRefreshTokenForPair(pairToken string) (models.RefreshTokenWithHash, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.refreshRepo.GetToken(ctx, pairToken)
}

func (s authService) ParseToken(token string) (string, error) {
	claims, err := jwt.ParseToken(token, s.secret)
	return claims.PairId, err
}

func (s authService) GetRefreshTokenForAccessToken(token string) (models.RefreshTokenWithHash, error) {
	claims, err := jwt.ParseToken(token, s.secret)
	if err != nil {
		return models.RefreshTokenWithHash{}, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.refreshRepo.GetToken(ctx, claims.PairId)
}

func (s authService) DeleteRefreshToken(pairToken string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.refreshRepo.DeleteToken(ctx, pairToken)
}
