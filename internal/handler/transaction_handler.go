package handler

import (
	"errors"
	"net/http"

	"github.com/Udean777/uang-bijak-go/internal/models"
	"github.com/Udean777/uang-bijak-go/internal/service"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	trxService service.TransactionService
}

func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{trxService: svc}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	userID, err := getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if models.TransactionType(req.Type) != models.TransactionExpense && models.TransactionType(req.Type) != models.TransactionIncome {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction type"})
		return
	}

	trx, err := h.trxService.CreateTransaction(c.Request.Context(), req, userID)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid wallet or category ID"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	c.JSON(http.StatusCreated, trx)
}

func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
	userID, err := getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	transactions, err := h.trxService.GetUserTransactions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
