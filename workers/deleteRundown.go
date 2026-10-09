package workers

import (
	"be-date/config"
	"be-date/models"
	"context"
	"fmt"
	"sync"
)

func DeleteWorker(
	wg *sync.WaitGroup,
	workerID int,
	jobs <-chan models.DeleteJob,
	results chan<- models.DeleteJobResult,
) {

	defer wg.Done()

	for job := range jobs {

		result := deleteOneData(job)

		results <- result
	}
}

func deleteOneData(
	job models.DeleteJob,
) models.DeleteJobResult {

	ctx := context.Background()

	result := models.DeleteJobResult{
		ID: job.ID,
	}

	tx, err := config.DB.Begin(ctx)

	if err != nil {
		result.Error = fmt.Errorf(
			"gagal memulai transaction: %w",
			err,
		)

		return result
	}

	defer tx.Rollback(ctx)

	// Ambil nama terlebih dahulu
	var nama string

	selectQuery := fmt.Sprintf(`
		SELECT nama
		FROM %s
		WHERE id = $1
	`, job.Table)

	err = tx.QueryRow(
		ctx,
		selectQuery,
		job.ID,
	).Scan(&nama)

	if err != nil {
		result.Error = fmt.Errorf(
			"data tidak ditemukan",
		)

		return result
	}

	result.Nama = nama

	// Delete
	deleteQuery := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = $1
	`, job.Table)

	deleteResult, err := tx.Exec(
		ctx,
		deleteQuery,
		job.ID,
	)

	if err != nil {
		result.Error = fmt.Errorf(
			"gagal menghapus data: %w",
			err,
		)

		return result
	}

	if deleteResult.RowsAffected() == 0 {
		result.Error = fmt.Errorf(
			"data gagal dihapus",
		)

		return result
	}

	// Commit
	if err := tx.Commit(ctx); err != nil {

		result.Error = fmt.Errorf(
			"gagal commit transaction: %w",
			err,
		)

		return result
	}

	result.Success = true

	return result
}
