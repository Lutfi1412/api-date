package services

import (
	"be-date/config"
	"be-date/models"
	"be-date/utils"
	"context"
	"fmt"
	"time"
)

func UpdateRundown(
	id int64,
	data models.UpdateRundown,
) error {

	loc, err := time.LoadLocation("Asia/Jakarta")

	ctx := context.Background()

	// Validasi jam_mulai
	jamMulai, err := utils.ParseJam(
		data.JamMulai,
		loc,
	)

	if err != nil {
		return fmt.Errorf(
			"jam_mulai tidak valid: %w",
			err,
		)
	}

	// Validasi jam_selesai
	jamSelesai, err := utils.ParseJam(
		data.JamSelesai,
		loc,
	)

	if err != nil {
		return fmt.Errorf(
			"jam_selesai tidak valid: %w",
			err,
		)
	}

	// Validasi urutan jam
	if !jamSelesai.After(jamMulai) {
		return fmt.Errorf(
			"jam_selesai harus lebih besar dari jam_mulai",
		)
	}

	// Mulai transaction
	tx, err := config.DB.Begin(ctx)

	if err != nil {
		return fmt.Errorf(
			"gagal memulai transaction: %w",
			err,
		)
	}

	defer tx.Rollback(ctx)

	query := `
		UPDATE rundown
		SET
			nama = $1,
			jam_mulai = $2,
			jam_selesai = $3,
			link_gmaps = $4,
			rating = $5
		WHERE id = $6
	`

	result, err := tx.Exec(
		ctx,
		query,
		data.Nama,
		data.JamMulai,
		data.JamSelesai,
		utils.NullIfEmpty(data.LinkGmaps),
		data.Rating,
		id,
	)

	if err != nil {
		return fmt.Errorf(
			"gagal update rundown: %w",
			err,
		)
	}

	// Pastikan ID ditemukan
	if result.RowsAffected() == 0 {
		return fmt.Errorf(
			"rundown dengan id %d tidak ditemukan",
			id,
		)
	}

	// Commit
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf(
			"gagal commit transaction: %w",
			err,
		)
	}

	return nil
}
