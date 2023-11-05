package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	database "go-fwallet/database"

	"github.com/gin-gonic/gin"
)

type Account struct {
	ID          int       `json:"id" bun:",pk,autoincrement"`
	Name        string    `json:"name"`
	AccountType string    `json:"accountType"`
	Currency    string    `json:"currency"`
	UserID      int       `json:"userID"`
	UpdatedAt   time.Time `json:"-" bun:"default:current_timestamp"`
}

func GetAllAccounts(c *gin.Context) {
	var accounts []Account
	err := database.DB.NewSelect().Model(&accounts).Scan(c)

	if err != nil {
		log.Printf("Error while getting all accounts, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong. See logs for more details",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All accounts",
		"data":    accounts,
	})
}

func CreateAccount(c *gin.Context) {
	var account Account
	c.BindJSON(&account)

	_, err := database.DB.NewInsert().Model(&account).Exec(c)

	if err != nil {
		log.Printf("Error while inserting new account into db, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Account created Successfully",
		"data":    account,
	})
}

func GetSingleAccount(c *gin.Context) {
	accountID, _ := strconv.Atoi(c.Param("accountID"))
	account := new(Account)
	err := database.DB.NewSelect().Model(account).Where("id = ?", accountID).Scan(c)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "Account not found",
			})
			return
		}
		log.Printf("Error while getting a single account, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single Account",
		"data":    account,
	})
}

func EditAccount(c *gin.Context) {
	accountID, _ := strconv.Atoi(c.Param("accountID"))
	account := &Account{
		ID:        accountID,
		UpdatedAt: time.Now(),
	}
	c.BindJSON(&account)

	res, err := database.DB.NewUpdate().Model(account).OmitZero().WherePK().Exec(c)
	row, _ := res.RowsAffected()

	if err != nil {
		log.Printf("Error while updating a Account, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	if row == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Account not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Account Edited Successfully",
	})
}

func DeleteAccount(c *gin.Context) {
	accountID, _ := strconv.Atoi(c.Param("accountID"))
	account := new(Account)

	res, err := database.DB.NewDelete().Model(account).Where("id = ?", accountID).Exec(c)
	if err != nil {
		log.Printf("Error while deleting an Account, Reason: %v\n", err)
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
			"message": "Account not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Account deleted successfully",
	})
}
