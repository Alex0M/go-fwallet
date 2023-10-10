package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	database "go-fwallet/database"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type Transaction struct {
	ID                  int       `json:"id" bun:",pk,autoincrement"`
	AccountID           int       `json:"accountID"`
	CategoryID          int       `json:"categoryID"`
	TransactionTypeCode string    `json:"transactionTypeCode"`
	Amount              int       `json:"amount"`
	Description         string    `json:"description"`
	TransactionDate     time.Time `json:"transactionDate"`
	UpdatedAt           time.Time `json:"-" bun:"default:current_timestamp"`
}

type GetTransaction struct {
	bun.BaseModel `bun:"transactions,alias:t"`

	ID                  int       `json:"id"`
	TransactionDate     time.Time `json:"transactionDate"`
	Category            string    `json:"category"`
	SubCategory         string    `json:"subCategory"`
	Amount              int       `json:"amount"`
	Account             string    `json:"account"`
	Description         string    `json:"description"`
	TransactionTypeCode string    `json:"transactionTypeCode"`
}

type SumTransactionByCategory struct {
	bun.BaseModel `bun:"transactions,alias:t"`

	Category string `json:"category"`
	Amount   int    `json:"amount"`
}

type TransactionCategoryFromBSD struct {
	GteDate             string `json:"gteDate" binding:"required"`
	LteDate             string `json:"lteDate" binding:"required"`
	TransactionTypeCode string `json:"transactionTypeCode" binding:"required"`
}

func GetAllTransaction(c *gin.Context) {
	categoryId, _ := strconv.Atoi(c.Param("categoryID"))
	gteDate, gteRes := c.GetQuery("date_gte")
	lteDate, lteRes := c.GetQuery("date_lte")
	tTypeCode, tTypeCodeRes := c.GetQuery("transaction_type_code")
	var transactions []GetTransaction

	q := database.DB.NewSelect().Model(&transactions).ColumnExpr("t.id, t.transaction_date, t.amount, t.description, t.transaction_type_code").
		ColumnExpr("a.name AS account").
		ColumnExpr("c.name AS sub_category").
		ColumnExpr("cc.name AS category").
		Join("LEFT JOIN categories AS c ON t.category_id = c.id").
		Join("LEFT JOIN categories AS cc ON c.parent_id = cc.id").
		Join("LEFT JOIN accounts AS a ON t.account_id = a.id")

	if categoryId != 0 {
		q.Where("cc.id = ?", categoryId)
	}

	if tTypeCodeRes {
		q.Where("t.transaction_type_code = ?", tTypeCode)
	}

	if gteRes && lteRes {
		q.Where("t.transaction_date >= ? and t.transaction_date <= ?", gteDate, lteDate)
	} else if gteRes {
		q.Where("t.transaction_date >= ?", gteDate)
	} else if lteRes {
		q.Where("t.transaction_date <= ?", lteDate)
	}

	err := q.Order("t.transaction_date").Scan(c)

	if err != nil {
		log.Printf("Error while getting all transactions, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong. See logs for more details",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All transactions",
		"data":    transactions,
	})
}

func GetSumOfTransactionsByCategory(c *gin.Context) {
	gteDate, gteRes := c.GetQuery("date_gte")
	lteDate, lteRes := c.GetQuery("date_lte")
	tTypeCode, tTypeCodeRes := c.GetQuery("transaction_type_code")
	var sumTxsByCategory []SumTransactionByCategory

	q := database.DB.NewSelect().Model(&sumTxsByCategory).ColumnExpr("sum(amount) as amount").
		ColumnExpr("c.name as category").
		Join("LEFT JOIN categories AS cc ON t.category_id = cc.id").
		Join("LEFT JOIN categories AS c ON cc.parent_id = c.id")

	if tTypeCodeRes {
		q.Where("t.transaction_type_code = ?", tTypeCode)
	}

	if gteRes && lteRes {
		q.Where("t.transaction_date >= ? and t.transaction_date <= ?", gteDate, lteDate)
	} else if gteRes {
		q.Where("t.transaction_date >= ?", gteDate)
	} else if lteRes {
		q.Where("t.transaction_date <= ?", lteDate)
	}

	err := q.Group("category").Scan(c)

	if err != nil {
		log.Printf("Error while getting all transactions, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong. See logs for more details",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Sum of transactions by category",
		"data":    sumTxsByCategory,
	})
}

func CreateTransaction(c *gin.Context) {
	var transaction Transaction
	c.BindJSON(&transaction)

	_, trnsError := database.DB.NewInsert().Model(&transaction).Exec(c)
	if trnsError != nil {
		log.Printf("Error while inserting new transaction into db, Reason: %v\n", trnsError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Transaction created Successfully",
		"data":    transaction,
	})
}

/*
ToDo: Update EditTransaction API

	func EditTransaction(c *gin.Context) {
		transactionId, _ := strconv.Atoi(c.Param("transactionID"))
		log.Println(transactionId)
		transaction := &Transaction{ID: transactionId}
		c.BindJSON(&transaction)
		log.Println(transaction)

		_, err := dbConnect.Model(transaction).WherePK().Update()
		if err != nil {
			log.Printf("Error while updating a transaction, Reason: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  500,
				"message": "Something went wrong",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "Transaction Edited Successfully",
		})
		return
	}
*/
func DeleteTransaction(c *gin.Context) {
	transactionId, _ := strconv.Atoi(c.Param("transactionID"))
	transaction := new(Transaction)

	res, err := database.DB.NewDelete().Model(transaction).Where("id = ?", transactionId).Exec(c)
	if err != nil {
		log.Printf("Error while deleting a transaction, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	row, _ := res.RowsAffected()
	if row == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Transaction not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Transaction deleted successfully",
	})
}

func UpdateTransactionCategory(c *gin.Context) {
	var tcbsd TransactionCategoryFromBSD
	err := c.BindJSON(&tcbsd)
	if err != nil {
		log.Printf("Error trying to BindJSON, Reason: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Bad Request",
		})
		return
	}

	//err := database.DB.NewSelect().Model(&transactions).Where("transaction_date >= ? and transaction_date <= ? and category_id = 1", gteDate, lteDate).Scan(c)
	res, err := database.DB.NewUpdate().
		Table("transactions", "bank_statement_descriptions").
		SetColumn("category_id", "bank_statement_descriptions.category_id").
		Set("updated_at = ?", time.Now()).
		Where("transactions.description = bank_statement_descriptions.bank_statement_description and transactions.transaction_date >= ? and transactions.transaction_date <= ? and transactions.transaction_type_code = ? and transactions.category_id = 1", tcbsd.GteDate, tcbsd.LteDate, tcbsd.TransactionTypeCode).
		Exec(c)

	if err != nil {
		log.Printf("Error while getting transactions, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	r, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error while getting transactions, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Query returned successfully. Row updated: %d", r),
	})
}
