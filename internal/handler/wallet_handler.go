package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Udean777/uang-bijak-go/internal/models"
	"github.com/Udean777/uang-bijak-go/internal/service"
	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletService service.WalletService
}

func NewWalletHandler(svc service.WalletService) *WalletHandler {
	return &WalletHandler{walletService: svc}
}

func (h *WalletHandler) CreateWallet(c *gin.Context) {
	userID, err := getAuthenticatedUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CreateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.walletService.CreateWallet(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wallet with this name already exists"})
		return
	}

	c.JSON(http.StatusCreated, wallet)
}

func (h *WalletHandler) GetUserWallets(c *gin.Context) {
	userID, err := getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	wallets, err := h.walletService.GetUserWallets(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve wallets"})
		return
	}

	c.JSON(http.StatusOK, wallets)
}

func (h *WalletHandler) UpdateWallet(c *gin.Context) {
	userID, err := getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	walletIDParam := c.Param("id")
	walletID, err := strconv.ParseInt(walletIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	var req models.UpdateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.walletService.UpdateWallet(c.Request.Context(), walletID, req, userID)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this wallet"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wallet with this name already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet updated successfully"})
}

func (h *WalletHandler) DeleteWallet(c *gin.Context) {
	userID, err := getAuthenticatedUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	walletIDParam := c.Param("id")
	walletID, err := strconv.ParseInt(walletIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	err = h.walletService.DeleteWallet(c.Request.Context(), walletID, userID)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this wallet"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete wallet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet deleted successfully"})
}
