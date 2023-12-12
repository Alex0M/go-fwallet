package database

import (
	"context"
	"database/sql"
	"go-fwallet/internal/helpers"
	"go-fwallet/internal/models"

	"go.uber.org/zap"
)

func (d *Database) GetCategories(c context.Context) ([]models.Category, error) {
	var categories []models.Category
	err := d.Client.NewSelect().Model(&categories).Scan(c)
	if err != nil {
		d.Logger.Error("error while getting all accounts", zap.Error(err))
		return nil, err
	}

	return categories, nil
}

func (d *Database) GetCategory(idStr string, c context.Context) (*models.Category, error) {
	d.Logger.Info("starting getting a category by id", zap.Any("categoryID", idStr))
	id, err := helpers.StringToInt(idStr)
	if err != nil {
		d.Logger.Error(ErrValidateID.Error(), zap.Error(err))
		return nil, ErrValidateID
	}

	var category models.Category
	err = d.Client.NewSelect().Model(&category).Where("id = ?", id).Scan(c)
	if err != nil {
		if err == sql.ErrNoRows {
			d.Logger.Info("category not found", zap.Any("categoryID", id))
			return nil, nil
		}
		d.Logger.Error("error while getting category", zap.Error(err))
		return nil, err
	}

	d.Logger.Info("get single category", zap.Any("category", &category))
	return &category, nil
}

func (d *Database) GetCategoryByName(name string, c context.Context) (*models.Category, error) {
	d.Logger.Info("starting getting a category by name", zap.Any("categoryName", name))
	var category models.Category
	err := d.Client.NewSelect().Model(&category).Where("name = ?", name).Scan(c)
	if err != nil {
		if err == sql.ErrNoRows {
			d.Logger.Info("category not found", zap.Any("categoryName", name))
			return nil, nil
		}
		d.Logger.Error("error while getting category", zap.Error(err))
		return nil, err
	}

	d.Logger.Info("get single category", zap.Any("category", &category))
	return &category, nil
}

func (d *Database) AddCategory(cat *models.Category, c context.Context) error {
	d.Logger.Info("starting addding new category into DB", zap.Any("category", cat))
	_, err := d.Client.NewInsert().Model(cat).Exec(c)
	if err != nil {
		d.Logger.Error("error while inserting new category into db", zap.Error(err))
		return err
	}
	d.Logger.Info("category was added successfully", zap.Any("category", cat))
	return nil
}

func (d *Database) EditCategory(idStr string, newCategory *models.Category, c context.Context) (*models.Category, error) {
	d.Logger.Info("start updating category", zap.Any("categoryID", idStr))
	oldCat, err := d.GetCategory(idStr, c)
	if err != nil {
		return nil, err
	}

	if oldCat == nil {
		return oldCat, nil
	}

	_, err = d.Client.NewUpdate().Model(newCategory).OmitZero().Where("id = ?", oldCat.ID).Returning("*").Exec(c)
	if err != nil {
		d.Logger.Error("error while updating category", zap.Error(err))
		return nil, err
	}

	d.Logger.Info("category was updated successfully", zap.Any("category", newCategory))
	return newCategory, nil
}

func (d *Database) DeleteCategory(idStr string, c context.Context) (*models.Category, error) {
	d.Logger.Info("start deleting category", zap.Any("categoryID", idStr))
	cat, err := d.GetCategory(idStr, c)
	if err != nil {
		return nil, err
	}

	if cat != nil {
		_, err = d.Client.NewDelete().Model(cat).WherePK().Exec(c)
		if err != nil {
			d.Logger.Error("error while deleting category", zap.Error(err))
			return nil, err
		}

		d.Logger.Info("category was deleted successfully", zap.Any("category", cat))
	}

	return cat, nil
}
