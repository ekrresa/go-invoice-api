package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ekrresa/invoice-api/pkg/models"
	"github.com/ekrresa/invoice-api/pkg/repository"
	"github.com/ekrresa/invoice-api/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
)

type invoiceHandler struct {
	repo repository.Repository
}

func NewInvoiceHandler(repo repository.Repository) *invoiceHandler {
	return &invoiceHandler{
		repo: repo,
	}
}

func (c *invoiceHandler) CreateInvoice(w http.ResponseWriter, r *http.Request, user *models.UserWithoutPassword) {
	var requestBody models.CreateInvoicePayload

	decodeErr := utils.DecodeJSONBody(w, r.Body, &requestBody)
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
	validateErr := validate.Struct(&requestBody)
	if validateErr != nil {
		utils.ErrorResponse(w, validateErr.Error(), http.StatusBadRequest)
		return
	}

	// dueDate, dueDateErr := time.Parse("2006-01-02", requestBody.DueDate)
	// if dueDateErr != nil {
	// 	utils.ErrorResponse(w, "Invalid date format for due date", http.StatusBadRequest)
	// 	return

	// }

	newInvoice := models.Invoice{
		ID:            strings.ToLower(string(ulid.Make().String())),
		UserID:        user.ID,
		Description:   requestBody.Description,
		Status:        models.InvoiceStatus(requestBody.Status),
		CustomerName:  requestBody.CustomerName,
		CustomerEmail: requestBody.CustomerEmail,
		Underpay:      requestBody.Underpay,
		Currency:      requestBody.Currency,
		Total:         requestBody.Total,
		DueDate:       time.Now().Add(time.Duration(time.Now().Day() * 12)),
	}

	var newInvoiceItems = make([]models.InvoiceItem, len(requestBody.Items))
	if requestBody.Items != nil {
		for index, item := range requestBody.Items {
			newItems := models.InvoiceItem{
				Name:      item.Name,
				Quantity:  item.Quantity,
				Price:     item.Price,
				InvoiceID: newInvoice.ID,
			}

			newInvoiceItems[index] = newItems
		}
	}

	createInvoiceErr := c.repo.CreateInvoice(&newInvoice, &newInvoiceItems)
	if createInvoiceErr != nil {
		utils.ErrorResponse(w, createInvoiceErr.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, &newInvoice, "Invoice created", http.StatusOK)
}
