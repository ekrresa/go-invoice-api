package handlers

import (
	"errors"
	"net/http"

	"github.com/ekrresa/invoice-api/pkg/helpers"
	"github.com/ekrresa/invoice-api/pkg/models"
	"github.com/ekrresa/invoice-api/pkg/repository"
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

	var decodeErr = helpers.DecodeJSONBody(w, r.Body, &requestBody)

	if decodeErr != nil {
		var requestError *helpers.RequestError

		if errors.As(decodeErr, &requestError) {
			helpers.ErrorResponse(w, requestError.Message, requestError.StatusCode)
		} else {
			helpers.ErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	var validate = validator.New()
	var validateErr = validate.Struct(&requestBody)
	if validateErr != nil {
		helpers.ErrorResponse(w, validateErr.Error(), http.StatusBadRequest)
		return
	}

	if _, err := c.repo.GetUserByEmail(requestBody.Email); err == nil {
		helpers.ErrorResponse(w, "User already exists", http.StatusBadRequest)
		return
	}

	var hashedPassword, passwordError = helpers.HashPassword(requestBody.Password)
	if passwordError != nil {
		helpers.ErrorResponse(w, "Password failed validation", http.StatusInternalServerError)
		return
	}

	requestBody.Password = hashedPassword

	var apiKey, createErr = c.repo.CreateUser(requestBody)
	if createErr != nil {
		helpers.ErrorResponse(w, createErr.Error(), http.StatusBadRequest)
		return
	}

	var responsePayload = map[string]string{
		"api_key": apiKey,
	}

	helpers.SuccessResponse(w, &responsePayload, "User created", http.StatusCreated)
}

func (c *UserHandler) RegenerateApiKey(w http.ResponseWriter, r *http.Request) {
	var requestBody models.RegenerateAPIKeyInput

	var decodeErr = helpers.DecodeJSONBody(w, r.Body, &requestBody)
	if decodeErr != nil {
		var requestError *helpers.RequestError

		if errors.As(decodeErr, &requestError) {
			helpers.ErrorResponse(w, requestError.Message, requestError.StatusCode)
		} else {
			helpers.ErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	var validate = validator.New()
	var validateErr = validate.Struct(&requestBody)
	if validateErr != nil {
		helpers.ErrorResponse(w, validateErr.Error(), http.StatusBadRequest)
		return
	}

	var user, userByEmailError = c.repo.GetUserByEmail(requestBody.Email)
	if userByEmailError != nil {
		helpers.ErrorResponse(w, "Invalid email/password", http.StatusBadRequest)
		return
	}

	var passwordIsCorrect = helpers.CheckPasswordHash(requestBody.Password, user.Password)
	if !passwordIsCorrect {
		helpers.ErrorResponse(w, "Invalid email/password", http.StatusBadRequest)
		return
	}

	var newApiKey = ulid.Make().String()
	var hashedApiKey = helpers.HashApiKey(newApiKey)

	var isUserUpdated, updateErr = c.repo.UpdateUserAPIKey(user.ID, hashedApiKey)
	if !isUserUpdated || updateErr != nil {
		helpers.ErrorResponse(w, "Failed to regenerate api key", http.StatusBadRequest)
		return
	}

	var responsePayload = map[string]string{
		"api_key": newApiKey,
	}

	helpers.SuccessResponse(w, &responsePayload, "Api key regenerated", http.StatusOK)
}
