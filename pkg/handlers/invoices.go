package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
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

	fmt.Println(requestBody)

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
	var invoice, err = c.repo.GetInvoiceOfAUser(invoiceID, user.ID)

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

func (c *invoiceHandler) FinalizeInvoice(w http.ResponseWriter, r *http.Request, user *models.User) {
	var invoiceID = chi.URLParam(r, "invoiceID")
	var prevInvoice, prevInvoiceErr = c.repo.GetInvoiceOfAUser(invoiceID, user.ID)

	if prevInvoiceErr != nil {
		if prevInvoiceErr == sql.ErrNoRows {
			helpers.ErrorResponse(w, "Invoice does not exist", http.StatusNotFound)
		} else {
			helpers.ErrorResponse(w, "Error getting invoice", http.StatusInternalServerError)
		}
		return
	}

	if prevInvoice.Status != models.Draft {
		helpers.ErrorResponse(w, "Invoice is already finalized", http.StatusBadRequest)
		return
	}

	prevInvoice.Status = models.Open
	var invoiceInput = models.UpdateInvoiceInput{
		Description:           prevInvoice.Description,
		Status:                &prevInvoice.Status,
		CustomerName:          &prevInvoice.CustomerName,
		CustomerEmail:         prevInvoice.CustomerEmail,
		AllowMultiplePayments: &prevInvoice.AllowMultiplePayments,
		Currency:              &prevInvoice.Currency,
		Total:                 &prevInvoice.Total,
		DueDate:               prevInvoice.DueDate,
	}

	var updatedInvoice, err = c.repo.UpdateInvoice(invoiceID, invoiceInput)
	if err != nil {
		log.Println(err)
		helpers.ErrorResponse(w, "Invoice not found", http.StatusInternalServerError)
		return
	}

	helpers.SuccessResponse(w, &updatedInvoice, "Invoice finalized", http.StatusOK)
}

func (c *invoiceHandler) UpdateInvoice(w http.ResponseWriter, r *http.Request, user *models.User) {
	var requestBody models.UpdateInvoiceInput

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

	var invoiceID = chi.URLParam(r, "invoiceID")
	var prevInvoice, prevInvoiceErr = c.repo.GetInvoiceOfAUser(invoiceID, user.ID)

	if prevInvoiceErr != nil {
		if prevInvoiceErr == sql.ErrNoRows {
			helpers.ErrorResponse(w, "Invoice does not exist", http.StatusNotFound)
		} else {
			helpers.ErrorResponse(w, "Error updating invoice", http.StatusInternalServerError)
		}
		return
	}

	if prevInvoice.Status != models.Draft {
		helpers.ErrorResponse(w, "Invoice is closed for updates", http.StatusBadRequest)
		return
	}

	mergeInvoices(&requestBody, &prevInvoice)

	var updatedInvoice, err = c.repo.UpdateInvoice(invoiceID, requestBody)
	if err != nil {
		helpers.ErrorResponse(w, "Error updating invoice", http.StatusInternalServerError)
		return
	}

	helpers.SuccessResponse(w, &updatedInvoice, "Invoice finalized", http.StatusOK)
}

func (c *invoiceHandler) DeleteInvoice(w http.ResponseWriter, r *http.Request, user *models.User) {
	var invoiceID = chi.URLParam(r, "invoiceID")
	var invoice, err = c.repo.GetInvoiceOfAUser(invoiceID, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.ErrorResponse(w, "Invoice does not exist", http.StatusNotFound)
		} else {
			helpers.ErrorResponse(w, "Error deleting invoice", http.StatusInternalServerError)
		}
		return
	}

	if invoice.Status != models.Draft {
		helpers.ErrorResponse(w, "Only draft invoices can be deleted", http.StatusBadRequest)
		return
	}

	var deleteErr = c.repo.DeleteInvoice(invoiceID)
	if deleteErr != nil {
		log.Println(deleteErr.Error(), "Line 196")
		helpers.ErrorResponse(w, "Error deleting invoice", http.StatusInternalServerError)
		return
	}

	helpers.SuccessResponse(w, nil, "Invoice deleted", http.StatusOK)
}

func (h *invoiceHandler) VoidInvoice(w http.ResponseWriter, r *http.Request, user *models.User) {
	var invoiceID = chi.URLParam(r, "invoiceID")
	var prevInvoice, prevInvoiceErr = h.repo.GetInvoiceOfAUser(invoiceID, user.ID)

	if prevInvoiceErr != nil {
		if prevInvoiceErr == sql.ErrNoRows {
			helpers.ErrorResponse(w, "Invoice does not exist", http.StatusNotFound)
		} else {
			helpers.ErrorResponse(w, "Error voiding this invoice", http.StatusInternalServerError)
		}
		return
	}

	if prevInvoice.Status != models.Open {
		helpers.ErrorResponse(w, "Only open invoices can be voided", http.StatusBadRequest)
		return
	}

	prevInvoice.Status = models.Void

	var updateInvoiceInput = models.UpdateInvoiceInput{
		Description:           prevInvoice.Description,
		Status:                &prevInvoice.Status,
		CustomerName:          &prevInvoice.CustomerName,
		CustomerEmail:         prevInvoice.CustomerEmail,
		AllowMultiplePayments: &prevInvoice.AllowMultiplePayments,
		Currency:              &prevInvoice.Currency,
		Total:                 &prevInvoice.Total,
		DueDate:               prevInvoice.DueDate,
	}

	var updatedInvoice, err = h.repo.UpdateInvoice(invoiceID, updateInvoiceInput)
	if err != nil {
		helpers.ErrorResponse(w, "Error voiding invoice", http.StatusInternalServerError)
		return
	}

	helpers.SuccessResponse(w, &updatedInvoice, "Invoice is void", http.StatusOK)
}

func mergeInvoices(dest *models.UpdateInvoiceInput, src *models.Invoice) {
	if dest.Description == nil {
		dest.Description = src.Description
	}
	if dest.Status == nil {
		dest.Status = &src.Status
	}
	if dest.CustomerName == nil {
		dest.CustomerName = &src.CustomerName
	}
	if dest.CustomerEmail == nil {
		dest.CustomerEmail = src.CustomerEmail
	}
	if dest.AllowMultiplePayments == nil {
		dest.AllowMultiplePayments = &src.AllowMultiplePayments
	}
	if dest.Currency == nil {
		dest.Currency = &src.Currency
	}
	if dest.Total == nil {
		dest.Total = &src.Total
	}
	if dest.DueDate == nil {
		dest.DueDate = src.DueDate
	}
}
