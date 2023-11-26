package controllers

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

type UserBody struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

func (c *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userBody UserBody
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	decodeErr := decoder.Decode(&userBody)

	if decodeErr != nil {
		if errors.As(decodeErr, &unmarshalErr) {
			utils.ErrorResponse(w, "Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			utils.ErrorResponse(w, "Unable to parse body", http.StatusBadRequest)
		}
		return
	}

	validate := validator.New()
	validateErr := validate.Struct(userBody)

	if validateErr != nil {
		utils.ErrorResponse(w, validateErr.Error(), http.StatusBadRequest)
		return
	}

	newUser := &models.User{}
	newUser.ID = strings.ToLower(ulid.Make().String())
	newUser.Name = userBody.Name
	newUser.Email = userBody.Email
	newUser.Password = userBody.Password

	createErr := c.repo.CreateUser(newUser)

	if createErr != nil {
		utils.ErrorResponse(w, createErr.Error(), http.StatusBadRequest)
		return
	}

	utils.SuccessResponse(w, &newUser, "User created", http.StatusCreated)
}
