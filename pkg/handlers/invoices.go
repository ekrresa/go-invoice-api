package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/ekrresa/invoice-api/pkg/helpers"
	"github.com/ekrresa/invoice-api/pkg/models"
	"github.com/ekrresa/invoice-api/pkg/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type invoiceHandler struct {
	repo repository.Repository
}

func NewInvoiceHandler(repo repository.Repository) *invoiceHandler {
	return &invoiceHandler{
		repo: repo,
	}
}

func (c *invoiceHandler) CreateInvoice(w http.ResponseWriter, r *http.Request, user *models.User) {
	var requestBody models.CreateInvoiceInput

	decodeErr := helpers.DecodeJSONBody(w, r.Body, &requestBody)
	if decodeErr != nil {
		var requestError *helpers.RequestError

		if errors.As(decodeErr, &requestError) {
			helpers.ErrorResponse(w, requestError.Message, requestError.StatusCode)
		} else {
			helpers.ErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	validate := validator.New()
	validateErr := validate.Struct(&requestBody)
	if validateErr != nil {
		helpers.ErrorResponse(w, validateErr.Error(), http.StatusBadRequest)
		return
	}

	//TODO: Validate due date that it is a valid timestamp and it is in the future.
	//TODO: Validate total amount if invoice has items.

	createInvoiceErr := c.repo.CreateInvoice(user.ID, &requestBody)
	if createInvoiceErr != nil {
		helpers.ErrorResponse(w, createInvoiceErr.Error(), http.StatusInternalServerError)
		return
	}

	helpers.SuccessResponse(w, nil, "Invoice created", http.StatusOK)
}

func (c *invoiceHandler) ListInvoicesOfUser(w http.ResponseWriter, r *http.Request, user *models.User) {
	var invoices, err = c.repo.ListInvoicesOfUser(user.ID)

	if err != nil {
		helpers.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.SuccessResponse(w, &invoices, "Invoices retrieved", http.StatusOK)
}

func (c *invoiceHandler) GetInvoice(w http.ResponseWriter, r *http.Request, user *models.User) {
	var invoiceID = chi.URLParam(r, "invoiceID")
	var invoice, err = c.repo.GetInvoice(invoiceID, user.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			helpers.ErrorResponse(w, "Invoice not found", http.StatusNotFound)
		} else {
			helpers.ErrorResponse(w, "Error getting invoice", http.StatusInternalServerError)
		}
		return
	}

	helpers.SuccessResponse(w, &invoice, "Invoice retrieved", http.StatusOK)
}
