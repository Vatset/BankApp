package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	balance_app "users-balance-monitoring"
)

// @Summary		Deposit
// @Tags			admin
// @Description	To deposit the user's account by the admin.
// @ID				deposit
// @Accept			json
// @Produce		json
// @Param			input	body		balance_app.TransactionList	true	"deposit"
// @Success		200		{string}	string						"token"
// @Failure		400,404	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api_admin/balance/deposit [post]
func (h *Handler) balanceDeposit(c *gin.Context) {
	var input balance_app.TransactionList

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.UserChecker(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	newBalance, err := h.service.Deposit(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"new balance": newBalance,
	})
}

// @Summary GetBalance
// @Tags user
// @Description The function of getting the user's balance. The default currency is 'RUB'. Optionally, the user can change the currency to any other fiat currency.
// @ID get-balance
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param currency query string false "Currency code (e.g., USD, EUR, etc.) to get the balance in. Default is RUB."
// @Success		200		{string}	string						"token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/balance [get]
func (h *Handler) getBalance(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	currency := c.DefaultQuery("currency", "RUB")

	balance, err := h.service.GetBalance(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if currency != "RUB" {
		convertedBalance, err := h.convertCurrency(balance, currency)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, "Failed to convert currency")
			return
		}
		c.JSON(http.StatusOK, map[string]interface{}{
			"balance":  convertedBalance,
			"currency": currency,
		})
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"balance": balance,
		})
	}

}

func (h *Handler) convertCurrency(balance float64, currency string) (float64, error) {
	client := http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.exchangerate.host/convert?from=RUB&to=%s&amount=%f", currency, balance), nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	if balanceConverted, ok := result["result"].(float64); ok {
		return balanceConverted, nil
	} else {
		return 0, err
	}
}

// @Summary		Withdraw
// @Tags			admin
// @Description	To withdraw the user's balance by the admin.
// @ID				withdraw
// @Accept			json
// @Produce		json
// @Param			input	body		balance_app.TransactionList	true	"withdraw"
// @Success		200		{string}	string						"token"
// @Failure		400,404	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api_admin/balance/withdraw [post]
func (h *Handler) balanceWithdraw(c *gin.Context) {
	var input balance_app.TransactionList

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.UserChecker(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	newBalance, err := h.service.Withdraw(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"new balance": newBalance,
	})
}

//	@Summary		BalanceTransfer
//
// @Security ApiKeyAuth
//
//	@Tags			user
//	@description	The function of transferring money from one user to another.
//	@description If the sender has a sufficient balance in their account and the recipient has an open account in the application, the transaction will be successful
//	@ID				transfer
//	@Accept			json
//	@Produce		json
//	@Param			input	body		balance_app.TransactionList	true	"balance transfer"
//	@Success		200		{string}	string						"token"
//	@Failure		400,404	{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Failure		default	{object}	errorResponse
//	@Router			/api/balance/transfer [post]
func (h *Handler) balanceTransfer(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input balance_app.TransactionList

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.UserChecker(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Transfer(userId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
