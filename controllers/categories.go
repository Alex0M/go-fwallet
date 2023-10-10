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

type Category struct {
	ID        int       `json:"id" bun:",pk,autoincrement"`
	Name      string    `json:"name"`
	ParentID  int       `json:"parentID" bun:",nullzero"`
	UpdatedAt time.Time `json:"-" bun:"default:current_timestamp"`
}

func GetAllCategories(c *gin.Context) {
	var categories []Category
	err := database.DB.NewSelect().Model(&categories).Scan(c)

	if err != nil {
		log.Printf("Error while getting all categories, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong. See logs for more details",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All categories",
		"data":    categories,
	})
}

func CreateCategory(c *gin.Context) {
	var category Category
	c.BindJSON(&category)

	_, err := database.DB.NewInsert().Model(&category).Exec(c)

	if err != nil {
		log.Printf("Error while inserting new category into db, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Category created Successfully",
		"data":    category,
	})
}

func GetSingleCategory(c *gin.Context) {
	categoryID, _ := strconv.Atoi(c.Param("categoryID"))
	category := new(Category)
	err := database.DB.NewSelect().Model(category).Where("id = ?", categoryID).Scan(c)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "Category not found",
			})
			return
		}
		log.Printf("Error while getting a single category, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single Category",
		"data":    category,
	})
}

func GetSingleCategoryByName(c *gin.Context) {
	categoryName := c.Param("categoryName")
	category := new(Category)
	err := database.DB.NewSelect().Model(category).Where("name = ?", categoryName).Scan(c)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "Category not found",
			})
			return
		}
		log.Printf("Error while getting a single category, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single Category",
		"data":    category,
	})
}

func EditCategory(c *gin.Context) {
	categoryID, _ := strconv.Atoi(c.Param("categoryID"))
	category := &Category{
		ID:        categoryID,
		UpdatedAt: time.Now(),
	}
	c.BindJSON(&category)

	res, err := database.DB.NewUpdate().Model(category).WherePK().Exec(c)
	row, _ := res.RowsAffected()

	if err != nil {
		log.Printf("Error while updating Category, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	if row == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Category not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Category Edited Successfully",
	})
}

func DeleteCategory(c *gin.Context) {
	categoryID, _ := strconv.Atoi(c.Param("categoryID"))
	category := new(Category)

	res, err := database.DB.NewDelete().Model(category).Where("id = ?", categoryID).Exec(c)
	if err != nil {
		log.Printf("Error while deleting Category, Reason: %v\n", err)
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
			"message": "Category not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Category deleted successfully",
	})
}
