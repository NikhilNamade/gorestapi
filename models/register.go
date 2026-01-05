package models

import (
	"database/sql"
	"errors"
	"fmt"

	"example.com/REST-API/db"
)

type Register struct {
	Id      int64
	UserId  int64
	EventId int64
}

func (res Register) Save() error {
	query := `INSERT INTO register (eventId,userId) VALUES(?,?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()
	fmt.Println(res.EventId, res.UserId)
	_, err = GetIDEvent(res.EventId)
	if err != nil {
		return err
	}

	ok := GetRegisterByUserId(res.UserId, res.EventId)
	if ok {
		return errors.New("User already have registration")
	}
	result, err := stmt.Exec(res.EventId, res.UserId)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	res.Id = id
	return nil
}

func Getall() ([]Register, error) {
	query := `SELECT * FROM register`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	var res = []Register{}
	defer rows.Close()
	for rows.Next() {
		var resi Register
		err := rows.Scan(&resi.Id, &resi.EventId, &resi.UserId)
		if err != nil {
			return nil, err
		}

		res = append(res, resi)
	}
	return res, nil
}

func GetRegisterByUserId(userId, eventId int64) bool {
	query := `SELECT userId FROM register WHERE userId = ? AND eventId = ?`
	row := db.DB.QueryRow(query, userId, eventId)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// ✅ No duplicate found — user not registered for this event
			return false
		}
		// ⚠️ Some other unexpected DB error — log it and treat as no duplicate
		fmt.Println("Error checking registration:", err)
		return false
	}

	// ✅ Record exists — user already registered for this event
	return true
}

func GetregisterById(id int64) int64 {
	query := `SELECT userId FROM register WHERE Id = ?`
	row := db.DB.QueryRow(query, id)
	var userId int64
	err := row.Scan(&userId)
	if err != nil {
		return 0
	}
	return userId
}
func Deleteresi(id,userId int64) error {
	query := `DELETE FROM register WHERE userId = ? AND eventId = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	fmt.Println(id,userId)
	_, err = stmt.Exec(userId,id)
	return err
}

func GetAllRegisterByUserId(userId int64)([]Register,error){
	query := `SELECT * FROM register WHERE userId = ?`
	rows,err :=db.DB.Query(query,userId)

	if err != nil{
		return  nil,err
	}

	defer rows.Close()

	var events = []Register{}

	for rows.Next(){
		var event Register
		err := rows.Scan(&event.Id,&event.EventId,&event.UserId)
		if err != nil{
			return nil,err
		}
		events = append(events, event)
	}
	return events,nil
}
func GetDetailEvents(userId int64)([]Event,error){
	query := `SELECT eventId FROM register WHERE userId = ?`
	rows,err :=db.DB.Query(query,userId)

	if err != nil{
		return  nil,err
	}

	defer rows.Close()

	var events = []Event{}

	for rows.Next(){
		var eventId int64
		var event Event
		err := rows.Scan(&eventId)
		if err != nil{
			return nil,err
		}

		event,err = GetIDEvent(eventId)

		if err != nil{
			return nil,err
		}

		events = append(events, event)
	}
	fmt.Println(events);
	return events,nil
}