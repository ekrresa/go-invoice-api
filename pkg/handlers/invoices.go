package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

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

	if requestBody.DueDate != nil && requestBody.DueDate.Before(time.Now()) {
		helpers.ErrorResponse(w, "Due date cannot be in the past", http.StatusBadRequest)
		return
	}

	if len(requestBody.Items) > 0 {
		var totalAmount uint
		for _, item := range requestBody.Items {
			totalAmount += item.Quantity * item.UnitPrice
		}

		if totalAmount != requestBody.Total {
			helpers.ErrorResponse(w, "Total amount does not match sum of items", http.StatusBadRequest)
			return
		}
	}

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

func (h *invoiceHandler) PayAnInvoice(w http.ResponseWriter, r *http.Request, user *models.User) {
	var requestBody models.PayInvoiceInput

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

	var invoiceID = chi.URLParam(r, "invoiceID")

	var invoice, getInvoiceErr = h.repo.GetInvoiceOfAUser(invoiceID, user.ID)
	if getInvoiceErr != nil {
		if getInvoiceErr == sql.ErrNoRows {
			helpers.ErrorResponse(w, "Invoice does not exist", http.StatusNotFound)
		} else {
			helpers.ErrorResponse(w, "An error occurred. Please try again", http.StatusInternalServerError)
		}
		return
	}

	if invoice.Status != models.Open {
		helpers.ErrorResponse(w, "This invoice cannot receive payments", http.StatusBadRequest)
		return
	}

	if invoice.Currency != requestBody.Currency {
		helpers.ErrorResponse(w, "Currency does not match with the invoice currency", http.StatusBadRequest)
		return
	}

	if !invoice.AllowMultiplePayments && requestBody.Amount != invoice.Total {
		helpers.ErrorResponse(w, "Amount does not match with the invoice total", http.StatusBadRequest)
		return
	}

	var invoiceAmountReceived = invoice.AmountPaid + requestBody.Amount

	if invoice.AllowMultiplePayments && invoiceAmountReceived > invoice.Total {
		helpers.ErrorResponse(w, "Amount cannot be greater than the invoice total", http.StatusBadRequest)
		return
	}

	var payment, err = h.repo.PayAnInvoice(invoice, requestBody, invoiceAmountReceived == invoice.Total)
	if err != nil {
		helpers.ErrorResponse(w, "Error paying invoice", http.StatusInternalServerError)
		return
	}

	helpers.SuccessResponse(w, &payment, "Payment was successful", http.StatusOK)
}

func (h *invoiceHandler) GetPaymentsOfInvoice(w http.ResponseWriter, r *http.Request, user *models.User) {
	var invoiceID = chi.URLParam(r, "invoiceID")

	var payments, err = h.repo.ListPaymentsOfAnInvoice(invoiceID)

	if err != nil {
		helpers.ErrorResponse(w, "Error getting payments", http.StatusInternalServerError)
		return
	}

	helpers.SuccessResponse(w, &payments, "Payments retrieved successfully", http.StatusOK)
}
