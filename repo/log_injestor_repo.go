package repo

import (
	"Log-Ingestor/contracts"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

type InjestLogRepo struct {
	db *sql.DB
}

func NewInjestLog(db *sql.DB) LogInjestorRepo {
	return &InjestLogRepo{
		db: db,
	}
}

type LogInjestorRepo interface {
	InjestLogs(logs []*contracts.LogEntry) http.Response
}

func (i InjestLogRepo) InjestLogs(logs []*contracts.LogEntry) http.Response {

	tx, err := i.db.Begin()
	if err != nil {
		log.Fatal(err)
		return http.Response{
			Status:     "InternalServerError",
			StatusCode: 500,
		}
	}

	query := fmt.Sprintf("INSERT INTO %s VALUES", "log")

	values := ""
	for i := range logs {
		if i > 0 {
			values += ","
		}
		values += " (?)"
	}

	query += values

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Fatal(err)
		return http.Response{
			Status:     "InternalServerError",
			StatusCode: 500,
		}
	}
	defer stmt.Close()

	var args []interface{}
	for _, logEntry := range logs {
		args = append(args, logEntry)
	}

	_, err = stmt.Exec(args...)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return http.Response{
			Status:     "InternalServerError",
			StatusCode: 500,
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return http.Response{
			Status:     "InternalServerError",
			StatusCode: 500,
		}
	}

	fmt.Println("Logs inserted successfully.")
	return http.Response{
		Status:     "Created",
		StatusCode: 201,
	}
}
