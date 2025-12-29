package models

import (
	"errors"
	"fmt"

	"example.com/REST-API/db"
	"example.com/REST-API/utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (user User) Save() error {
	query := `INSERT INTO user(name,email,password) VALUES(?,?,?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()
	hashPass, err := utils.HashPass(user.Password)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(user.Name, user.Email, hashPass)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	user.ID = id
	return err
}

func GetAllUsers() ([]User, error) {
	query := `SELECT id, name , email FROM user`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	var users = []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func LoginUser(user User) (string, error) {
	query := `SELECT id,password FROM user WHERE email = ?`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer stmt.Close()
	var userPassword string
	var id int64
	err = stmt.QueryRow(user.Email).Scan(&id, &userPassword)
	fmt.Println(id)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(user.Password))

	if err != nil {
		fmt.Println(err)
		fmt.Println("Invalid password")
		return "", errors.New("invalid password")
	}

	token, err := utils.GenerateToken(id)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return token, err
}

func GetUserByToken(token string) (User, error) {
	userId, err := utils.AuthenticateUser(token)

	if err != nil {
		return User{}, err
	}
	var user User
	query := `SELECT id,name,email FROM user WHERE id = ?`
	row := db.DB.QueryRow(query, userId)
	err = row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func GetUserById(id int64) (User, error) {
	if id == 0 {
		return User{}, errors.New("no Id found")
	}
	query := `
		SELECT id,name,email FROM user WHERE id = ?
	`
	var user User
	row := db.DB.QueryRow(query, id)

	err := row.Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		return User{}, err
	}
	return user, nil
}
