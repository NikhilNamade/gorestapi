package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"example.com/REST-API/db"
)

type Address struct {
	Location string  `json:"location" form:"location" binding:"required"`
	Lat      float64 `json:"lat" form:"lat" binding:"required"`
	Lgn      float64 `json:"lgn" form:"lgn" binding:"required"`
}

type Event struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" form:"name" binding:"required"`
	Description string    `json:"description" form:"description" binding:"required"`
	Location    Address   `json:"location" binding:"required"`
	Datetime    time.Time `json:"date_time" form:"date_time"`
	UserId      int64     `json:"userId"`
	Profile     string    `json:"profile_image"`
	Category    string    `json:"category"`
	Fees        int       `json:"fees"`
}

/* =========================
   CREATE EVENT
========================= */

func (e *Event) Save() error {
	query := `
	INSERT INTO event 
	(name, description, location, date_time, userId, profileImage, category, fees)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	locationJSON, err := json.Marshal(e.Location)
	if err != nil {
		return err
	}

	formattedTime := e.Datetime.Format("2006-01-02 15:04:05")

	result, err := stmt.Exec(
		e.Name,
		e.Description,
		locationJSON,
		formattedTime,
		e.UserId,
		e.Profile,
		e.Category,
		e.Fees,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	e.ID = id
	return err
}

/* =========================
   GET ALL EVENTS
========================= */

func GetAllEvents(userId int64) ([]Event, error) {
	query := `SELECT * FROM event WHERE userId != ?`

	rows, err := db.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		var dtestr sql.NullString
		var locationJSON []byte

		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&locationJSON,
			&dtestr,
			&event.UserId,
			&event.Profile,
			&event.Category,
			&event.Fees,
		)
		if err != nil {
			return nil, err
		}

		_ = json.Unmarshal(locationJSON, &event.Location)
		if dtestr.Valid {
			// try DB format first
			if t, err := time.Parse("2006-01-02 15:04:05", dtestr.String); err == nil {
				event.Datetime = t
			} else if t, err := time.Parse(time.RFC3339, dtestr.String); err == nil {
				event.Datetime = t
			}
		}

		events = append(events, event)
	}

	return events, nil
}

func GetIDEvent(id int64) (Event, error) {
	query := `SELECT * FROM event WHERE id = ?`
	row := db.DB.QueryRow(query, id)

	var event Event
	var dtestr sql.NullString
	var locationJSON []byte

	err := row.Scan(
		&event.ID,
		&event.Name,
		&event.Description,
		&locationJSON,
		&dtestr,
		&event.UserId,
		&event.Profile,
		&event.Category,
		&event.Fees,
	)
	if err != nil {
		return Event{}, err
	}

	_ = json.Unmarshal(locationJSON, &event.Location)
	if dtestr.Valid {
			// try DB format first
			if t, err := time.Parse("2006-01-02 15:04:05", dtestr.String); err == nil {
				event.Datetime = t
			} else if t, err := time.Parse(time.RFC3339, dtestr.String); err == nil {
				event.Datetime = t
			}
		}

	return event, nil
}

func (event Event) UpdateEvent() error {
	query := `
	UPDATE event 
	SET name = ?, description = ?, location = ?, date_time = ?, 
	    profileImage = ?, category = ?, fees = ?
	WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	locationJSON, err := json.Marshal(event.Location)
	if err != nil {
		return err
	}

	formattedTime := event.Datetime.Format("2006-01-02 15:04:05")

	_, err = stmt.Exec(
		event.Name,
		event.Description,
		locationJSON,
		formattedTime,
		event.Profile,
		event.Category,
		event.Fees,
		event.ID,
	)
	return err
}

/* =========================
   DELETE EVENT
========================= */

func (event Event) DeleteEvent() error {
	query := `DELETE FROM event WHERE id = ?`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(event.ID)
	return err
}

/* =========================
   GET EVENTS BY USER
========================= */

func GetEventsByUser(userId int64) ([]map[string]any, error) {
	query := `
	SELECT 
		e.id,
		e.name,
		e.description,
		e.location,
		e.date_time,
		e.profileImage,
		e.category,
		e.fees,
		COUNT(r.eventId) as countOfRegisters
	FROM event e
	LEFT JOIN register r ON e.id = r.eventId
	WHERE e.userId = ?
	GROUP BY e.id
	`

	rows, err := db.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []map[string]any

	for rows.Next() {
		var event Event
		var dtestr sql.NullString
		var locationJSON []byte
		var count int

		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&locationJSON,
			&dtestr,
			&event.Profile,
			&event.Category,
			&event.Fees,
			&count,
		)
		if err != nil {
			return nil, err
		}

		_ = json.Unmarshal(locationJSON, &event.Location)
		if dtestr.Valid {
			// try DB format first
			if t, err := time.Parse("2006-01-02 15:04:05", dtestr.String); err == nil {
				event.Datetime = t
			} else if t, err := time.Parse(time.RFC3339, dtestr.String); err == nil {
				event.Datetime = t
			}
		}

		events = append(events, map[string]any{
			"id":                  event.ID,
			"name":                event.Name,
			"description":         event.Description,
			"location":            event.Location,
			"date_time":           event.Datetime,
			"profile_image":       event.Profile,
			"category":            event.Category,
			"fees":                event.Fees,
			"countOfRegistration": count,
		})
	}

	return events, nil
}


