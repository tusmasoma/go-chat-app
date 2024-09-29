package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-chat-app/usecase"
)

type UserHandler interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	uuc usecase.UserUseCase
}

func NewUserHandler(uuc usecase.UserUseCase) UserHandler {
	return &userHandler{
		uuc: uuc,
	}
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uh *userHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody SignUpRequest
	defer r.Body.Close()
	if !uh.isValidSignUpRequest(r.Body, &requestBody) {
		log.Warn("Invalid request body", log.Fstring("email", requestBody.Email))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := uh.uuc.SignUpAndGenerateToken(ctx, requestBody.Email, requestBody.Password)
	if err != nil {
		log.Error("Failed to create user and generate token", log.Fstring("email", requestBody.Email), log.Ferror(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Info("User sign up successfully", log.Fstring("email", requestBody.Email))
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}

func (uh *userHandler) isValidSignUpRequest(body io.ReadCloser, requestBody *SignUpRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Error("Failed to decode request body: %v", err)
		return false
	}
	if requestBody.Email == "" || requestBody.Password == "" {
		log.Warn("Invalid request body: %v", requestBody)
		return false
	}
	return true
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uh *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody LoginRequest
	if ok := isValidLoginRequest(r.Body, &requestBody); !ok {
		log.Info("Invalid user login request", log.Fstring("method", r.Method), log.Fstring("url", r.URL.String()))
		http.Error(w, "Invalid user login request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	jwt, err := uh.uuc.LoginAndGenerateToken(ctx, requestBody.Email, requestBody.Password)
	if err != nil {
		log.Error("Failed to login or generate token", log.Fstring("email", requestBody.Email), log.Ferror(err))
		http.Error(w, "Failed to Login or generate token", http.StatusInternalServerError)
		return
	}

	log.Info("User login successfully", log.Fstring("email", requestBody.Email))
	w.Header().Set("Authorization", "Bearer "+jwt)
	w.WriteHeader(http.StatusOK)
}

func isValidLoginRequest(body io.ReadCloser, requestBody *LoginRequest) bool {
	// リクエストボディのJSONを構造体にデコード
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Error("Invalid request body", log.Ferror(err))
		return false
	}
	if requestBody.Email == "" || requestBody.Password == "" {
		log.Info("Missing required fields", log.Fstring("email", requestBody.Email))
		return false
	}
	return true
}
