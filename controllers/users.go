package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	database "go-fwallet/database"
	helpers "go-fwallet/helpers"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID        int       `json:"id" bun:",pk,autoincrement"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"-" bun:"default:current_timestamp"`
}

type CreateUserResponse struct {
	ID int `json:"id"`
}

func GetAllUsers(c *gin.Context) {
	var users []User
	err := database.DB.NewSelect().Model(&users).Scan(c)

	if err != nil {
		log.Printf("Error while getting all users, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong. See logs for more details",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All users",
		"data":    users,
	})
}

func CreateUser(c *gin.Context) {
	var user User
	c.BindJSON(&user)

	user.Password = helpers.HashAndSalt([]byte(user.Password))
	/*
		userInsert := User{
			Username: user.Username,
			Password: helpers.HashAndSalt([]byte(user.Password)),
			Email:    user.Email,
		}
	*/
	_, err := database.DB.NewInsert().Model(&user).Exec(c)
	if err != nil {
		log.Printf("Error while inserting new user into db, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "User created Successfully",
		"data": CreateUserResponse{
			ID: user.ID,
		},
	})
}

func GetSingleUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("userId"))
	user := new(User)
	err := database.DB.NewSelect().Model(user).Where("id = ?", userId).Scan(c)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "User not found",
			})
			return
		}
		log.Printf("Error while getting a single user, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single User",
		"data":    user,
	})
}

func EditUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("userId"))
	user := &User{
		ID:        userId,
		UpdatedAt: time.Now(),
	}
	c.BindJSON(&user)

	user.Password = helpers.HashAndSalt([]byte(user.Password))

	res, err := database.DB.NewUpdate().Model(user).WherePK().Exec(c)
	row, _ := res.RowsAffected()

	if err != nil {
		log.Printf("Error while updating a user, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	if row == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "User Edited Successfully",
	})
}

func DeleteUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("userId"))
	user := new(User)

	res, err := database.DB.NewDelete().Model(user).Where("id = ?", userId).Exec(c)
	if err != nil {
		log.Printf("Error while deleting a user, Reason: %v\n", err)
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
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "User deleted successfully",
	})
}
