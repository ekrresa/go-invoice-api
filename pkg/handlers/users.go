package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ekrresa/invoice-api/pkg/models"
	"github.com/ekrresa/invoice-api/pkg/repository"
	"github.com/ekrresa/invoice-api/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
)

type UserController struct {
	repo repository.UserRepository
}

func NewUserController(repo repository.UserRepository) *UserController {
	return &UserController{
		repo: repo,
	}
}

type user struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

func (c *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var responseBody user
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	decodeErr := decoder.Decode(&responseBody)

	if decodeErr != nil {
		if errors.As(decodeErr, &unmarshalErr) {
			utils.ErrorResponse(w, "Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			utils.ErrorResponse(w, "Unable to parse body", http.StatusInternalServerError)
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
