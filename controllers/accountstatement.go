package controllers

import (
	database "go-fwallet/database"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AccountStatementRequestPayload struct {
	GteDate string `json:"gteDate"`
	LteDate string `json:"lteDate"`
}

type AccountStatement struct {
	AccountID      int       `json:"accountID"`
	ClosingDate    time.Time `json:"closingDate"`
	ClosingBalance int       `json:"closingBalance"`
	TotalCredit    int       `json:"totalCredit"`
	TotalDebit     int       `json:"totalDebit"`
	UpdatedAt      time.Time `json:"-" bun:"default:current_timestamp"`
}

func GetAccountsStatement(c *gin.Context) {
	var accountStatements []AccountStatement
	cDate, cDateRes := c.GetQuery("closing_date")

	q := database.DB.NewSelect().Model(&accountStatements)
	if cDateRes {
		q.Where("closing_date = ?", cDate)
	}

	err := q.Order("closing_date").Scan(c)
	if err != nil {
		log.Printf("Error reading accounts statement from DB: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Accounts statement",
		"data":    accountStatements,
	})

}

func CreateAccountsStatement(c *gin.Context) {
	var requestPayload AccountStatementRequestPayload
	accountStatements := make([]AccountStatement, 0)
	c.BindJSON(&requestPayload)

	closingDate, err := time.Parse("2006-01-02", requestPayload.LteDate)
	if err != nil {
		log.Printf("Error parsing LteDate: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Bad request",
		})
		return
	}

	err = database.DB.NewRaw(
		"SELECT account_id, sum(amount) FILTER (where transaction_type_code = 'W') as total_credit, sum(amount) FILTER (where transaction_type_code = 'D') as total_debit FROM public.transactions WHERE transaction_date >= ? and transaction_date <= ? group by account_id",
		requestPayload.GteDate, requestPayload.LteDate,
	).Scan(c, &accountStatements)

	if err != nil {
		log.Printf("Error reading transactions for account statements from DB: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	for i := range accountStatements {
		accountStatements[i].ClosingBalance = accountStatements[i].TotalDebit - accountStatements[i].TotalCredit
		accountStatements[i].ClosingDate = closingDate.AddDate(0, 0, 1)
	}

	_, err = database.DB.NewInsert().Model(&accountStatements).Exec(c)
	if err != nil {
		log.Printf("Error inserting account statements into DB: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Accounts statement created successfully",
		"data":    accountStatements,
	})
}
