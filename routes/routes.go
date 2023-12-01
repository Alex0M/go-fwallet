package routes

import (
	"net/http"

	controllers "go-fwallet/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.GET("/", welcome)

	router.GET("/users", controllers.GetAllUsers)
	router.POST("/users", controllers.CreateUser)
	router.GET("/users/:userId", controllers.GetSingleUser)
	router.PUT("/users/:userId", controllers.EditUser)
	router.DELETE("/users/:userId", controllers.DeleteUser)

	router.GET("/transactiontypes", controllers.GetAllTransactionTypes)
	router.POST("/transactiontypes", controllers.CreateTransactionType)
	router.GET("/transactiontypes/:transactionTypeCode", controllers.GetSingleTransactionType)
	router.PUT("/transactiontypes/:transactionTypeCode", controllers.EditTransactionType)
	router.DELETE("/transactiontypes/:transactionTypeCode", controllers.DeleteTransactionType)

	router.GET("/categories", controllers.GetAllCategories)
	router.POST("/categories", controllers.CreateCategory)
	router.GET("/categories/:categoryID", controllers.GetSingleCategory)
	router.GET("/categories/name/:categoryName", controllers.GetSingleCategoryByName)
	router.PUT("/categories/:categoryID", controllers.EditCategory)
	router.DELETE("/categories/:categoryID", controllers.DeleteCategory)

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
