package server

import (
	"authentication_service/internal/token"
	"authentication_service/services"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthUserHandler interface {
	refreshHandler(w http.ResponseWriter, r *http.Request)
	registerHandler(w http.ResponseWriter, r *http.Request)
}

type authUserHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) AuthUserHandler {
	return authUserHandler{
		authService: authService,
	}
}

func (h authUserHandler) refreshHandler(w http.ResponseWriter, r *http.Request) {
	oldRefreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	oldAccessTokenHeader := r.Header.Get("Authentication")
	if oldAccessTokenHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	oldAccessToken := strings.TrimPrefix(oldAccessTokenHeader, "Bearer ")

	oldRefreshToken, err := h.authService.GetRefreshTokenForAccessToken(oldAccessToken)
	if err != nil || len(oldRefreshToken.Hash) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword(oldRefreshToken.Hash, []byte(oldRefreshTokenCookie.Value)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userId := oldRefreshToken.UserId

	h.authService.DeleteRefreshToken(oldRefreshToken.PairToken)

	pairToken := token.GenerateToken()

	accessToken := h.authService.GenerateAccessToken(userId, pairToken)
	refreshToken, err := h.authService.GenerateRefreshToken(userId, pairToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.Value,
		HttpOnly: true,
	}
	http.SetCookie(w, &refreshCookie)

	json.NewEncoder(w).Encode(accessToken)
}

func (h authUserHandler) registerHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pairToken := token.GenerateToken()

	accessToken := h.authService.GenerateAccessToken(userId, pairToken)
	refreshToken, err := h.authService.GenerateRefreshToken(userId, pairToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.Value,
		HttpOnly: true,
	}
	http.SetCookie(w, &refreshCookie)

	json.NewEncoder(w).Encode(accessToken)

}
