package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ekrresa/invoice-api/pkg/models"
	"github.com/ekrresa/invoice-api/pkg/repository"
	"github.com/ekrresa/invoice-api/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
)

type UserHandler struct {
	repo repository.Repository
}

func NewUserHandler(repo repository.Repository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

type user struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=8"`
}

func (c *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var responseBody user

	decodeErr := utils.DecodeJSONBody(w, r.Body, &responseBody)

	if decodeErr != nil {
		var requestError *utils.RequestError

		if errors.As(decodeErr, &requestError) {
			utils.ErrorResponse(w, requestError.Message, requestError.StatusCode)
		} else {
			utils.ErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	validate := validator.New()
	validateErr := validate.Struct(&responseBody)
	if validateErr != nil {
		utils.ErrorResponse(w, validateErr.Error(), http.StatusBadRequest)
		return
	}

	_, userByEmailError := c.repo.GetUserByEmail(responseBody.Email)
	if userByEmailError == nil {
		utils.ErrorResponse(w, "User already exists", http.StatusBadRequest)
		return
	}

	newUser := &models.User{}
	newUser.ID = strings.ToLower(ulid.Make().String())
	newUser.Name = responseBody.Name
	newUser.Email = responseBody.Email
	newUser.ApiKey = strings.ToLower(ulid.Make().String())

	userPassword, passwordError := utils.HashPassword(responseBody.Password)
	if passwordError != nil {
		utils.ErrorResponse(w, "Password failed validation", http.StatusInternalServerError)
		return
	}

	newUser.Password = userPassword

	createErr := c.repo.CreateUser(newUser)
	if createErr != nil {
		utils.ErrorResponse(w, createErr.Error(), http.StatusBadRequest)
		return
	}

	responsePayload := make(map[string]string)
	responsePayload["api_key"] = newUser.ApiKey

	utils.SuccessResponse(w, &responsePayload, "User created", http.StatusCreated)
}

type regeneratePayload struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=8"`
}

func (c *UserHandler) RegenerateApiKey(w http.ResponseWriter, r *http.Request) {
	var responseBody regeneratePayload

	decodeErr := utils.DecodeJSONBody(w, r.Body, &responseBody)
	if decodeErr != nil {
		var requestError *utils.RequestError

		if errors.As(decodeErr, &requestError) {
			utils.ErrorResponse(w, requestError.Message, requestError.StatusCode)
		} else {
			utils.ErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	validate := validator.New()
	validateErr := validate.Struct(&responseBody)
	if validateErr != nil {
		utils.ErrorResponse(w, validateErr.Error(), http.StatusBadRequest)
		return
	}

	user, userByEmailError := c.repo.GetUserByEmail(responseBody.Email)
	if userByEmailError != nil {
		utils.ErrorResponse(w, "Invalid email/password", http.StatusBadRequest)
		return
	}

	passwordIsCorrect := utils.CheckPasswordHash(responseBody.Password, user.Password)
	if !passwordIsCorrect {
		utils.ErrorResponse(w, "Invalid email/password", http.StatusBadRequest)
		return
	}

	newApiKey := strings.ToLower(ulid.Make().String())
	user.ApiKey = newApiKey

	updateErr := c.repo.UpdateUser(user)
	if updateErr != nil {
		utils.ErrorResponse(w, "Failed to regenerate api key", http.StatusBadRequest)
		return
	}

	responsePayload := make(map[string]string)
	responsePayload["api_key"] = newApiKey

	utils.SuccessResponse(w, &responsePayload, "Api key regenerated", http.StatusOK)
}
