package models

import (
	"fmt"
	"time"

	"example.com/REST-API/db"
)

type Event struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name" form:"name" binding:"required"`
	Description  string    `json:"description" form:"description" binding:"required"`
	Location     string    `json:"location" form:"location" binding:"required"`
	Datetime     time.Time `json:"date_time" form:"date_time" time_format:"2006-01-02T15:04:05Z07:00"`
	UserId       int64     `json:"userId" form:"userId"`
	Profile      string    `json:"profile_image" form:"profile_image" binding:"required"`
	Category     string    `json:"category" form:"category" binding:"required"`
	Fees         int       `json:"fees" form:"fees"`
}

var events = []Event{}

func (e *Event) Save() error {
	query := `INSERT INTO event (name,description,location,date_time,userId,profileImage,category,fees)
	VALUES(?,?,?,?,?,?,?,?)`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()
	fmt.Println(e.Datetime)
	formattedTime := e.Datetime.Format("2006-01-02 15:04:05")
	fmt.Println(formattedTime)
	result, err := stmt.Exec(e.Name, e.Description, e.Location, formattedTime, e.UserId, e.Profile, e.Category, e.Fees)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	e.ID = id
	return err
}

func GetAllEvents(userId int64) ([]Event, error) {
	query := `SELECT * FROM event WHERE userId != ?`
	rows, err := db.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	var events []Event
	defer rows.Close()
	for rows.Next() {
		var event Event
		var dtestr string
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &dtestr, &event.UserId, &event.Profile, &event.Category, &event.Fees)

		if err != nil {
			return nil, err
		}
		event.Datetime, _ = time.Parse(time.RFC3339, dtestr)
		events = append(events, event)
	}
	return events, nil
}

func GetIDEvent(id int64) (Event, error) {
	query := `SELECT * FROM event WHERE id = ?`
	row := db.DB.QueryRow(query, id)
	var event Event
	var dtestr string
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &dtestr, &event.UserId, &event.Profile, &event.Category, &event.Fees)

	if err != nil {
		return Event{}, err
	}
	event.Datetime, _ = time.Parse(time.RFC3339, dtestr)
	return event, nil
}

func (event Event) UpdateEvent() error {
	query := `UPDATE event SET name = ?,description = ?,location = ?,date_time = ?, profileImage= ?, category = ?, fees = ? WHERE id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Print("Error1")
		return err
	}
	defer stmt.Close()

	formattedTime := event.Datetime.Format("2006-01-02 15:04:05")
	_, err = stmt.Exec(event.Name, event.Description, event.Location, formattedTime, event.Profile, event.Category, event.Fees, event.ID)
	if err != nil {
		fmt.Print("Error2", err)
		return err
	}
	return nil
}

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

func GetEventsByUser(userId int64) ([]map[string]any, error) {
	//query := `SELECT id,name,description,location,date_time,profileImage FROM event WHERE userId = ?`
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
			COUNT(r.eventId) as CountOfRegisters
			FROM event e
			LEFT JOIN register r ON e.id = r.eventId
			WHERE e.userId = ?
			GROUP BY e.id, e.name, e.description, e.location, e.date_time, e.profileImage
	`
	rows, err := db.DB.Query(query, userId)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var events []map[string]any

	for rows.Next() {
		var dtestr string
		var event Event
		var countOfRegistration int
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &dtestr, &event.Profile, &event.Category, &event.Fees, &countOfRegistration)
		if err != nil {
			return nil, err
		}
		event.Datetime, _ = time.Parse(time.RFC3339, dtestr)
		actualEvent := map[string]any{
			"id":                  event.ID,
			"name":                event.Name,
			"description":         event.Description,
			"location":            event.Location,
			"date_time":           event.Datetime,
			"profile_image":       event.Profile,
			"category":            event.Category,
			"fees":                event.Fees,
			"countOfRegistration": countOfRegistration,
		}
		events = append(events, actualEvent)
	}
	return events, nil
}


