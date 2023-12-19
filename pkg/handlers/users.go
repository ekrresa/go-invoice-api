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

func (c *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var requestBody models.CreateUserInput

	var decodeErr = utils.DecodeJSONBody(w, r.Body, &requestBody)

	if decodeErr != nil {
		var requestError *utils.RequestError

		if errors.As(decodeErr, &requestError) {
			utils.ErrorResponse(w, requestError.Message, requestError.StatusCode)
		} else {
			utils.ErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	var validate = validator.New()
	var validateErr = validate.Struct(&requestBody)
	if validateErr != nil {
		utils.ErrorResponse(w, validateErr.Error(), http.StatusBadRequest)
		return
	}

	if _, err := c.repo.GetUserByEmail(requestBody.Email); err == nil {
		utils.ErrorResponse(w, "User already exists", http.StatusBadRequest)
		return
	}

	var hashedPassword, passwordError = utils.HashPassword(requestBody.Password)
	if passwordError != nil {
		utils.ErrorResponse(w, "Password failed validation", http.StatusInternalServerError)
		return
	}

	requestBody.Password = hashedPassword

	var apiKey, createErr = c.repo.CreateUser(requestBody)
	if createErr != nil {
		utils.ErrorResponse(w, createErr.Error(), http.StatusBadRequest)
		return
	}

	var responsePayload = map[string]string{
		"api_key": apiKey,
	}

	utils.SuccessResponse(w, &responsePayload, "User created", http.StatusCreated)
}

func (c *UserHandler) RegenerateApiKey(w http.ResponseWriter, r *http.Request) {
	var requestBody models.RegenerateAPIKeyInput

	var decodeErr = utils.DecodeJSONBody(w, r.Body, &requestBody)
	if decodeErr != nil {
		var requestError *utils.RequestError

		if errors.As(decodeErr, &requestError) {
			utils.ErrorResponse(w, requestError.Message, requestError.StatusCode)
		} else {
			utils.ErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	var validate = validator.New()
	var validateErr = validate.Struct(&requestBody)
	if validateErr != nil {
		utils.ErrorResponse(w, validateErr.Error(), http.StatusBadRequest)
		return
	}

	var user, userByEmailError = c.repo.GetUserByEmail(requestBody.Email)
	if userByEmailError != nil {
		utils.ErrorResponse(w, "Invalid email/password", http.StatusBadRequest)
		return
	}

	var passwordIsCorrect = utils.CheckPasswordHash(requestBody.Password, user.Password)
	if !passwordIsCorrect {
		utils.ErrorResponse(w, "Invalid email/password", http.StatusBadRequest)
		return
	}

	var newApiKey = strings.ToLower(ulid.Make().String())
	var hashedApiKey = utils.HashApiKey(newApiKey)

	var isUserUpdated, updateErr = c.repo.UpdateUserAPIKey(user.ID, hashedApiKey)
	if !isUserUpdated || updateErr != nil {
		utils.ErrorResponse(w, "Failed to regenerate api key", http.StatusBadRequest)
		return
	}

	var responsePayload = map[string]string{
		"api_key": newApiKey,
	}

	utils.SuccessResponse(w, &responsePayload, "Api key regenerated", http.StatusOK)
}
