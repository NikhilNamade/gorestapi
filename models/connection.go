package models

import (
	"database/sql"
	"errors"
	"fmt"

	"example.com/REST-API/db"
)

type Connection struct {
	ID       int
	FollowBy int `json:"followBy" binding:"required"`
	FollowTo int `json:"followTo" binding:"required"`
}

func (c *Connection) Save() error {
	var exists int
	err := db.DB.QueryRow("SELECT id FROM user WHERE id = ?", c.FollowTo).Scan(&exists)
	if err != nil {
		return errors.New("followTo user not found")
	}

	_, err = db.DB.Exec("INSERT INTO Connection (followBy, followTo) VALUES (?, ?)", c.FollowBy, c.FollowTo)
	if err != nil {
		return err
	}
	return nil
}

func GetAllConnectionByUser(id int) ([]Connection, error) {
	query := `
	SELECT * FROM Connection 
	WHERE followBy = ? OR followTo = ?
	`

	var connections = []Connection{}
	rows, err := db.DB.Query(query, id, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var connect Connection
		err := rows.Scan(&connect.ID, &connect.FollowBy, &connect.FollowTo)
		if err != nil {
			return nil, err
		}
		connections = append(connections, connect)
	}
	return connections, nil
}

func DisConnect(toId int, byId int) error {
	var exists int

	// Check if connection exists
	err := db.DB.QueryRow(
		"SELECT id FROM Connection WHERE followBy = ? AND followTo = ?",
		byId, toId,
	).Scan(&exists)

	if err == sql.ErrNoRows {
		// No connection found
		return fmt.Errorf("connection does not exist")
	}

	if err != nil {
		// Some other DB error
		return err
	}

	// If exists, delete it
	_, err = db.DB.Exec(
		"DELETE FROM Connection WHERE followBy = ? AND followTo = ?",
		byId, toId,
	)

	if err != nil {
		return err
	}

	return nil
}
