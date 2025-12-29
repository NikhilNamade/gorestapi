package models

import (
	"errors"
	"time"

	"example.com/REST-API/db"
)

type Story struct {
	ID        int
	File      string     `form:"file" binding:"required"`
	UserID    int        `form:"userId"`
	ViewStory *time.Time `form:"viewStory" json:"viewStory"`
}

func (story *Story) Save() error {
	query := `INSERT INTO story (file,userID) VALUES (?,?)`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(story.File, story.UserID)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	story.ID = int(id)
	return nil
}

func GetAllStory() ([]Story, error) {
	query := `SELECT * FROM story`
	var stories = []Story{}
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var story Story

		err := rows.Scan(&story.ID, &story.File, &story.UserID, &story.ViewStory)
		if err != nil {
			return nil, err
		}
		stories = append(stories, story)
	}
	return stories, nil
}

func GetuserStory(id int) ([]Story, error) {
	if id == 0 {
		return nil, errors.New("User is not Authenticated")
	}
	var stories = []Story{}
	query := `SELECT * FROM story WHERE userId = ?`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	result, err := stmt.Query(id)

	if err != nil {
		return nil, err
	}

	for result.Next() {
		var story Story
		err := result.Scan(&story.ID, &story.File, &story.UserID, &story.ViewStory)
		if err != nil {
			return nil, err
		}
		stories = append(stories, story)
	}
	return stories, nil
}

func UpdateStoryView(id int, date time.Time) error {
	if id == 0 {
		return errors.New("Id not found")
	}

	query := `
		UPDATE story SET viewStory = ? WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(date, id)
	if err != nil {
		return err
	}
	return nil
}
