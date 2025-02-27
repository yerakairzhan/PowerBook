// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
	"time"
)

type ReadingLog struct {
	ID          int32        `json:"id"`
	Userid      string       `json:"userid"`
	Date        time.Time    `json:"date"`
	MinutesRead int32        `json:"minutes_read"`
	CreatedAt   sql.NullTime `json:"created_at"`
}

type User struct {
	ID         int32          `json:"id"`
	Userid     string         `json:"userid"`
	Username   string         `json:"username"`
	Registered sql.NullBool   `json:"registered"`
	Language   sql.NullString `json:"language"`
	Timer      time.Time      `json:"timer"`
	State      sql.NullString `json:"state"`
	CreatedAt  sql.NullTime   `json:"created_at"`
}
