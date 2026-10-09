package workers

import (
	"be-date/config"
	"be-date/models"
	"be-date/utils"
	"context"
	"fmt"
	"sync"
	"time"
)

const WorkerCount = 3

func WorkerCreate(
	wg *sync.WaitGroup,
	workerID int,
	jobs <-chan models.RundownJob,
	results chan<- models.RundownJobResult,
) {
	defer wg.Done()

	for job := range jobs {

		err := ProcessRundownJob(job)

		result := models.RundownJobResult{
			JobID:     job.ID,
			TanggalID: job.TanggalID,
			Data:      job.Data,
			Success:   err == nil,
		}

		if err != nil {
			result.Error = err
			result.ErrorDetail = err.Error()
		}

		results <- result
	}
}

func ProcessRundownJob(job models.RundownJob) error {

	ctx := context.Background()

	// ==========================================
	// 1. Validasi jam
	// ==========================================

	loc, err := time.LoadLocation("Asia/Jakarta")

	if err != nil {
		return fmt.Errorf(
			"gagal memuat timezone: %w",
			err,
		)
	}

	jamMulai, err := utils.ParseJam(
		job.Data.JamMulai,
		loc,
	)

	if err != nil {
		return fmt.Errorf(
			"jam_mulai tidak valid: %w",
			err,
		)
	}

	jamSelesai, err := utils.ParseJam(
		job.Data.JamSelesai,
		loc,
	)

	if err != nil {
		return fmt.Errorf(
			"jam_selesai tidak valid: %w",
			err,
		)
	}

	if !jamSelesai.After(jamMulai) {
		return fmt.Errorf(
			"jam_selesai harus lebih besar dari jam_mulai",
		)
	}

	// ==========================================
	// 3. Mulai transaction
	// ==========================================

	tx, err := config.DB.Begin(ctx)

	if err != nil {
		return fmt.Errorf(
			"gagal memulai transaction: %w",
			err,
		)
	}

	defer tx.Rollback(ctx)

	// ==========================================
	// 4. Insert rundown
	// ==========================================

	query := `
		INSERT INTO rundown (
			nama,
			jam_mulai,
			jam_selesai,
			tanggal_id,
			link_gmaps,
			rating
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)
	`

	_, err = tx.Exec(
		ctx,
		query,
		job.Data.Nama,

		// PENTING:
		// kirim string, bukan time.Time
		job.Data.JamMulai,

		// PENTING:
		// kirim string, bukan time.Time
		job.Data.JamSelesai,

		job.TanggalID,
		utils.NullIfEmpty(job.Data.LinkGmaps),
		job.Data.Rating,
	)

	if err != nil {
		return fmt.Errorf(
			"gagal insert rundown: %w",
			err,
		)
	}

	// ==========================================
	// 5. Commit
	// ==========================================

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf(
			"gagal commit transaction: %w",
			err,
		)
	}

	return nil
}
