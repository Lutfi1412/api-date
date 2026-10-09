package services

import (
	"be-date/config"
	"be-date/models"
	"be-date/workers"
	"context"
	"fmt"
	"sync"
)

func CreateRundown(
	tanggalID int64,
	data []models.CreateRundown,
) error {

	ctx := context.Background()

	// ==========================================
	// 1. Pastikan tanggal_id tersedia
	// ==========================================

	// var existsID int64

	// checkQuery := `
	// 	SELECT id
	// 	FROM tanggal
	// 	WHERE id = $1
	// `

	// err := config.DB.QueryRow(
	// 	ctx,
	// 	checkQuery,
	// 	tanggalID,
	// ).Scan(&existsID)

	// if err != nil {
	// 	return fmt.Errorf(
	// 		"tanggal_id %d tidak ditemukan",
	// 		tanggalID,
	// 	)
	// }

	// ==========================================
	// 2. Update status menjadi proses
	// ==========================================

	updateQuery := `
		UPDATE tanggal
		SET status_rundown = 'proses'
		WHERE id = $1
	`

	_, err := config.DB.Exec(
		ctx,
		updateQuery,
		tanggalID,
	)

	if err != nil {
		return fmt.Errorf(
			"gagal mengubah status_rundown menjadi proses: %w",
			err,
		)
	}

	// ==========================================
	// 3. Buat channel
	// ==========================================

	jobs := make(chan models.RundownJob)
	results := make(chan models.RundownJobResult)

	var wg sync.WaitGroup

	// ==========================================
	// 4. Jalankan 3 worker
	// ==========================================

	wg.Add(workers.WorkerCount)

	for i := 0; i < workers.WorkerCount; i++ {
		go workers.WorkerCreate(
			&wg,
			i+1,
			jobs,
			results,
		)
	}

	// ==========================================
	// 5. Kirim semua job
	// ==========================================

	go func() {
		for i, item := range data {

			jobs <- models.RundownJob{
				ID:        i + 1,
				TanggalID: tanggalID,
				Data:      item,
			}
		}

		close(jobs)
	}()

	// ==========================================
	// 6. Tunggu worker selesai
	// ==========================================

	go func() {
		wg.Wait()
		close(results)
	}()

	// ==========================================
	// 7. Ambil semua hasil
	// ==========================================

	var hasError bool

	for result := range results {

		if result.Success {
			continue
		}

		hasError = true

		// Kalau job gagal,
		// simpan ke tabel error
		err := saveError(result)

		if err != nil {
			fmt.Printf(
				"gagal menyimpan error job %d: %v\n",
				result.JobID,
				err,
			)
		}
	}

	// ==========================================
	// 8. Update status akhir
	// ==========================================

	finalStatus := "sukses"

	if hasError {
		finalStatus = "peringatan"
	}

	err = updateFinalStatus(
		tanggalID,
		finalStatus,
	)

	if err != nil {
		return fmt.Errorf(
			"gagal mengubah status akhir rundown: %w",
			err,
		)
	}

	return nil
}

func saveError(result models.RundownJobResult) error {

	ctx := context.Background()

	query := `
		INSERT INTO error (
			nama,
			jam_mulai,
			jam_selesai,
			error_detail,
			tanggal_id,
			link_gmaps,
			rating
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := config.DB.Exec(
		ctx,
		query,
		result.Data.Nama,
		result.Data.JamMulai,
		result.Data.JamSelesai,
		result.ErrorDetail,
		result.TanggalID,
		result.Data.LinkGmaps,
		result.Data.Rating,
	)

	return err
}

func updateFinalStatus(
	tanggalID int64,
	status string,
) error {

	ctx := context.Background()

	tx, err := config.DB.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	query := `
		UPDATE tanggal
		SET status_rundown = $1
		WHERE id = $2
	`

	_, err = tx.Exec(
		ctx,
		query,
		status,
		tanggalID,
	)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
