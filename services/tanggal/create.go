package services

import (
	"be-date/config"
	"context"
	"fmt"
	"os"
)

func CreateTanggal(tanggal string, status string) (string, error) {

	input := `INSERT INTO tanggal (tanggal, status) VALUES ($1, $2)`

	conn := config.DB

	tx, err := conn.Begin(context.Background())

	_, err = conn.Exec(context.Background(), input, tanggal, status)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert data: %v\n", err)
		tx.Rollback(context.Background())
	}

	err = tx.Commit(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to commit transaction: %v\n", err)
		tx.Rollback(context.Background())
	}

	return "tanggal berhasil disimpan", nil
}
