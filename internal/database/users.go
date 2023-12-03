package database

import (
	"context"
	"database/sql"
	"go-fwallet/internal/helpers"
	"go-fwallet/internal/models"

	"go.uber.org/zap"
)

func (d *Database) GetUsers(c context.Context) ([]models.User, error) {
	var u []models.User
	err := d.Client.NewSelect().Model(&u).Scan(c)
	if err != nil {
		d.Logger.Error("error while getting all accounts", zap.Error(err))
		return nil, err
	}

	return u, nil
}

func (d *Database) GetUser(idStr string, c context.Context) (*models.User, error) {
	d.Logger.Info("starting getting a user by id", zap.Any("userID", idStr))
	id, err := helpers.StringToInt(idStr)
	if err != nil {
		d.Logger.Error(ErrValidateID.Error(), zap.Error(err))
		return nil, ErrValidateID
	}

	var u models.User
	err = d.Client.NewSelect().Model(&u).Where("id = ?", id).Scan(c)
	if err != nil {
		if err == sql.ErrNoRows {
			d.Logger.Info("users not found", zap.Any("userID", id))
			return nil, nil
		}
		d.Logger.Error("error while getting user", zap.Error(err))
		return nil, err
	}

	d.Logger.Info("get single user", zap.Any("user", &u))
	return &u, nil
}

func (d *Database) AddUser(u *models.User, c context.Context) error {
	d.Logger.Info("starting addding new user into DB", zap.Any("user", u))
	hashedPass, err := helpers.HashAndSalt([]byte(u.Password))
	if err != nil {
		d.Logger.Error("error while hashing the user password", zap.Error(err))
		return err
	}

	u.Password = hashedPass
	_, err = d.Client.NewInsert().Model(u).Exec(c)
	if err != nil {
		d.Logger.Error("error while inserting new user into db", zap.Error(err))
		return err
	}
	d.Logger.Info("user was added successfully", zap.Any("user", u))
	return nil
}

func (d *Database) EditUser(idStr string, newUser *models.User, c context.Context) (*models.User, error) {
	d.Logger.Info("start updating user", zap.Any("userID", idStr))
	oldUsr, err := d.GetUser(idStr, c)
	if err != nil {
		return nil, err
	}

	if oldUsr == nil {
		return oldUsr, nil
	}

	if newUser.Password != "" {
		hashedPass, err := helpers.HashAndSalt([]byte(newUser.Password))
		if err != nil {
			d.Logger.Error("error while hashing the user password", zap.Error(err))
			return nil, err
		}
		newUser.Password = hashedPass
	}

	_, err = d.Client.NewUpdate().Model(newUser).OmitZero().Where("id = ?", oldUsr.ID).Returning("*").Exec(c)
	if err != nil {
		d.Logger.Error("error while updating user", zap.Error(err))
		return nil, err
	}

	d.Logger.Info("user was updated successfully", zap.Any("account", newUser))
	return newUser, nil
}

func (d *Database) DeleteUser(idStr string, c context.Context) (*models.User, error) {
	d.Logger.Info("start deleting user", zap.Any("userID", idStr))
	u, err := d.GetUser(idStr, c)
	if err != nil {
		return nil, err
	}

	if u != nil {
		_, err = d.Client.NewDelete().Model(u).WherePK().Exec(c)
		if err != nil {
			d.Logger.Error("error while deleting user", zap.Error(err))
			return nil, err
		}

		d.Logger.Info("user was deleted successfully", zap.Any("account", u))
	}

	return u, nil
}
