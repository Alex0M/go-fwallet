package routes

import (
	"net/http"

	controllers "go-fwallet/controllers"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {

	/*router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})) */

	router.GET("/", welcome)

	router.GET("/transactiontypes", controllers.GetAllTransactionTypes)
	router.POST("/transactiontypes", controllers.CreateTransactionType)
	router.GET("/transactiontypes/:transactionTypeCode", controllers.GetSingleTransactionType)
	router.PUT("/transactiontypes/:transactionTypeCode", controllers.EditTransactionType)
	router.DELETE("/transactiontypes/:transactionTypeCode", controllers.DeleteTransactionType)

	router.GET("/accountstatements", controllers.GetAccountsStatement)
	router.POST("/accountstatements", controllers.CreateAccountsStatement)

	router.NoRoute(notFound)
}

func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Welcome to FWallet API",
	})
}

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  404,
		"message": "Route Not Found",
	})
}
