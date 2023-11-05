package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	database "go-fwallet/database"

	"github.com/gin-gonic/gin"
)

type TransactionType struct {
	TransactionTypeCode string    `json:"transactionTypeCode" bun:",pk"`
	Description         string    `json:"description"`
	UpdatedAt           time.Time `json:"-" bun:"default:current_timestamp"`
}

func GetAllTransactionTypes(c *gin.Context) {
	var transactionTypes []TransactionType
	err := database.DB.NewSelect().Model(&transactionTypes).Scan(c)

	if err != nil {
		log.Printf("Error while getting all transaction types, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong. See logs for more details",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All transaction types",
		"data":    transactionTypes,
	})
}

func CreateTransactionType(c *gin.Context) {
	var transactionType TransactionType
	c.BindJSON(&transactionType)

	_, err := database.DB.NewInsert().Model(&transactionType).Exec(c)

	if err != nil {
		log.Printf("Error while inserting new transaction type into db, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Transaction Type created Successfully",
		"data":    transactionType,
	})
}

func GetSingleTransactionType(c *gin.Context) {
	transactionTypeCode := c.Param("transactionTypeCode")
	transactionType := new(TransactionType)

	err := database.DB.NewSelect().Model(transactionType).Where("transaction_type_code = ?", transactionTypeCode).Scan(c)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "Transaction type not found",
			})
			return
		}
		log.Printf("Error while getting a single transaction type, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single Transaction Type",
		"data":    transactionType,
	})
}

func EditTransactionType(c *gin.Context) {
	transactionTypeCode := c.Param("transactionTypeCode")
	transactionType := &TransactionType{
		TransactionTypeCode: transactionTypeCode,
		UpdatedAt:           time.Now()}
	c.BindJSON(&transactionType)

	res, err := database.DB.NewUpdate().Model(transactionType).OmitZero().WherePK().Exec(c)
	row, _ := res.RowsAffected()

	if err != nil {
		log.Printf("Error while updating a Transaction Type, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	if row == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Transaction not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Transaction Type Edited Successfully",
	})
}

func DeleteTransactionType(c *gin.Context) {
	transactionTypeCode := c.Param("transactionTypeCode")
	transactionType := new(TransactionType)

	res, err := database.DB.NewDelete().Model(transactionType).Where("transaction_type_code = ?", transactionTypeCode).Exec(c)
	if err != nil {
		log.Printf("Error while deleting a Transaction Type, Reason: %v\n", err)
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
			"message": "Transaction Type not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Transaction Type deleted successfully",
	})
}
