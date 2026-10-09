package services

import (
	"be-date/config"
	"context"
	"fmt"
	"os"
)

func DeleteTanggal(id int) (string, error) {
	query := `DELETE FROM tanggal WHERE id = $1;`

	conn := config.DB

	tx, err := conn.Begin(context.Background())

	_, err = tx.Exec(context.Background(), query, id)

	err = tx.Commit(context.Background())

	if err != nil {
		fmt.Fprintf(os.Stderr, "gagal hapus data: %v\n", err)
		tx.Rollback(context.Background())
	}

	return "tanggal berhasil dihapus", nil
}
