package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"time"
	balance_app "users-balance-monitoring"
)

type getAllTransactionsResponse struct {
	Data       []balance_app.TransactionList `json:"data"`
	Page       int                           `json:"page"`
	TotalPages int                           `json:"total_pages"`
}

// @Summary GetAllTransactions
// @Tags transactions
// @description  -----The user can get the transaction history:-----
// @description -by the exact amount
// @description -by a date interval
// @description -by the amount interval
//
//	@description -----Sorting is implemented:-----
//	@description -in ascending order
//	@description -in descending order
//	@description -----Implemented pagination:-----
//	@description -It is necessary to specify the page number.
//
// @description -The maximum number of transactions per page.
// @ID history
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sort_by query string false "Sort by (date, amount)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Param sort_field query string false "Sort field (amount, amount_interval, date_interval) amount is for exact value of the amount"
// @Param exact_value query string false "Exact value of the amount"
// @Param start_value query string false "Start value for interval"
// @Param end_value query string false "End value for interval"
// @Param limit query int false "Limit of transactions per page"
// @Param page query int false "Page number"
// @Success 200 {object} getAllTransactionsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/balance/history [get]
func (h *Handler) getAllTransactions(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	sortBy := c.DefaultQuery("sort_by", "")
	sortOrder := c.DefaultQuery("sort_order", "")
	sortField := c.DefaultQuery("sort_field", "")
	exactValue := c.DefaultQuery("exact_value", "")
	startValue := c.DefaultQuery("start_value", "")
	endValue := c.DefaultQuery("end_value", "")
	limitOfTransactionsOnPage := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := strconv.Atoi(limitOfTransactionsOnPage)
	if err != nil || limit < 1 {
		newErrorResponse(c, http.StatusBadRequest, "Invalid limit parameter")
		return
	}

	var requestParameter string
	var requestOrder string

	switch sortField {
	case "amount":
		requestParameter = " AND amount = " + exactValue
	case "":

	case "amount_interval":
		requestParameter = fmt.Sprintf(" AND amount >= '%s' AND amount <= '%s'", startValue, endValue)
	case "date_interval":
		startDateValue, err := time.Parse("02-01-2006", startValue)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "Invalid start date format, please use dd-mm-yyyy")
			return
		}
		startDay := startDateValue.Format("2006-01-02")
		endDateValue, err := time.Parse("02-01-2006", endValue)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "Invalid end date format, please use dd-mm-yyyy")
			return
		}
		endDay := endDateValue.Format("2006-01-02")
		requestParameter = fmt.Sprintf(" AND date >= '%s' AND date < '%s'", startDay, endDay)
	default:
		newErrorResponse(c, http.StatusBadRequest, "Requested date or amount wasnt found")
		return
	}
	switch sortBy {
	case "date":
		requestOrder = " ORDER BY date"
	case "amount":
		requestOrder = " ORDER BY amount"
	case "":

	default:
		newErrorResponse(c, http.StatusBadRequest, "Invalid sort_by parameter")
		return
	}

	switch sortOrder {
	case "desc":
		requestOrder += " DESC "
	case "asc":
		requestOrder += " ASC "
	case "":

	default:
		newErrorResponse(c, http.StatusBadRequest, "Invalid sort_order parameter")
		return
	}

	transactions, countOfTransactions, err := h.service.GetAllTransactions(userId, requestParameter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if countOfTransactions == 0 {
		newErrorResponse(c, http.StatusBadRequest, "Transaction history was not found by your request")
		return
	}

	countOfPages := int(math.Ceil(float64(countOfTransactions) / float64(limit)))

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 || page > countOfPages {
		newErrorResponse(c, http.StatusBadRequest, "Invalid page parameter")
		return
	}

	transactionOffset := (page - 1) * limit

	transactions, err = h.service.PaginationTransactions(userId, requestParameter, requestOrder, limit, transactionOffset)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllTransactionsResponse{
		Data:       transactions,
		Page:       page,
		TotalPages: countOfPages,
	})
}
