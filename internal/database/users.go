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

func (d *Database) GetUsers(c context.Context) ([]models.User, error) {
	var u []models.User
	err := d.Client.NewSelect().Model(&u).Scan(c)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (d *Database) GetUser(idStr string, c context.Context) (*models.User, error) {
	id, err := helpers.StringToInt(idStr)
	if err != nil {
		return nil, middleware.NewHttpError("cannot conver user id to int", err.Error(), http.StatusBadRequest)
	}

	var u models.User
	err = d.Client.NewSelect().Model(&u).Where("id = ?", id).Scan(c)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, middleware.NewHttpError("user not found", fmt.Sprintf("userID:%d", id), http.StatusNotFound)
		}
		return nil, err
	}

	return &u, nil
}

func (d *Database) AddUser(u *models.User, c context.Context) error {
	hashedPass, err := helpers.HashAndSalt([]byte(u.Password))
	if err != nil {
		return err
	}

	u.Password = hashedPass
	_, err = d.Client.NewInsert().Model(u).Exec(c)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) EditUser(idStr string, newUser *models.User, c context.Context) (*models.User, error) {
	oldUsr, err := d.GetUser(idStr, c)
	if err != nil {
		return nil, err
	}

	if newUser.Password != "" {
		hashedPass, err := helpers.HashAndSalt([]byte(newUser.Password))
		if err != nil {
			return nil, err
		}
		newUser.Password = hashedPass
	}

	_, err = d.Client.NewUpdate().Model(newUser).OmitZero().Where("id = ?", oldUsr.ID).Returning("*").Exec(c)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (d *Database) DeleteUser(idStr string, c context.Context) error {
	u, err := d.GetUser(idStr, c)
	if err != nil {
		return err
	}

	_, err = d.Client.NewDelete().Model(u).WherePK().Exec(c)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) LoginCheck(u *models.UserLogin, c context.Context) (*models.User, error) {
	var user models.User
	err := d.Client.NewSelect().Model(&user).Where("email = ?", u.Email).Scan(c)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, middleware.NewHttpError("login or password is incorect", err.Error(), http.StatusUnauthorized)
		}
		return nil, err
	}

	err = helpers.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		return nil, middleware.NewHttpError("login or password is incorect", err.Error(), http.StatusUnauthorized)
	}

	return &user, nil
}
