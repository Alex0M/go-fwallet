package database

import (
	"context"
	"database/sql"
	"fmt"
	"go-fwallet/internal/helpers"
	"go-fwallet/internal/middleware"
	"go-fwallet/internal/models"
	"net/http"
)

func (d *Database) GetCategories(c context.Context) ([]models.Category, error) {
	var categories []models.Category
	err := d.Client.NewSelect().Model(&categories).Scan(c)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (d *Database) GetCategory(idStr string, c context.Context) (*models.Category, error) {
	id, err := helpers.StringToInt(idStr)
	if err != nil {
		return nil, middleware.NewHttpError("cannot conver category id to int", err.Error(), http.StatusBadRequest)
	}

	var category models.Category
	err = d.Client.NewSelect().Model(&category).Where("id = ?", id).Scan(c)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, middleware.NewHttpError("category not found", fmt.Sprintf("categoryID:%d", id), http.StatusNotFound)
		}
		return nil, err
	}

	return &category, nil
}

func (d *Database) GetCategoryByName(name string, c context.Context) (*models.Category, error) {
	var category models.Category
	err := d.Client.NewSelect().Model(&category).Where("name = ?", name).Scan(c)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, middleware.NewHttpError("category not found", fmt.Sprintf("CategoryNmae:%s", name), http.StatusNotFound)
		}
		return nil, err
	}

	return &category, nil
}

func (d *Database) AddCategory(cat *models.Category, c context.Context) error {
	_, err := d.Client.NewInsert().Model(cat).Exec(c)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) EditCategory(idStr string, newCategory *models.Category, c context.Context) (*models.Category, error) {
	oldCat, err := d.GetCategory(idStr, c)
	if err != nil {
		return nil, err
	}

	_, err = d.Client.NewUpdate().Model(newCategory).OmitZero().Where("id = ?", oldCat.ID).Returning("*").Exec(c)
	if err != nil {
		return nil, err
	}

	return newCategory, nil
}

func (d *Database) DeleteCategory(idStr string, c context.Context) error {
	cat, err := d.GetCategory(idStr, c)
	if err != nil {
		return err
	}

	_, err = d.Client.NewDelete().Model(cat).WherePK().Exec(c)
	if err != nil {
		return err
	}

	return nil
}
