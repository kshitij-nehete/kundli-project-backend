package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/kshitij-nehete/astro-report/internal/auth"
	"github.com/kshitij-nehete/astro-report/internal/handler/dto"
	"github.com/kshitij-nehete/astro-report/internal/response"
	"github.com/kshitij-nehete/astro-report/internal/usecase"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
	jwtService  *auth.JWTService
	validator   *validator.Validate
}

func NewAuthHandler(
	authUsecase *usecase.AuthUsecase,
	jwtService *auth.JWTService,
) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		jwtService:  jwtService,
		validator:   validator.New(),
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req dto.RegisterRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		response.WriteJSONError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.authUsecase.Register(
		r.Context(),
		req.Name,
		req.Email,
		req.Password,
	)

	if err != nil {
		response.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"user registered successfully"}`))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req dto.LoginRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		response.WriteJSONError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authUsecase.Login(
		r.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		response.WriteJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := h.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		response.WriteJSONError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	response := map[string]string{
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
