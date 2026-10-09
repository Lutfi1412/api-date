package services

import (
	"be-date/config"
	"be-date/utils"
	"context"
	"fmt"
	"time"
)

func TryAgainRundown(nama string, jamMulai string, jamSelesai string, tanggalID int64, linkGmaps string, rating *float64, id int64) error {

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

	jamMulaiParse, err := utils.ParseJam(
		jamMulai,
		loc,
	)

	if err != nil {
		return fmt.Errorf(
			"jam_mulai tidak valid: %w",
			err,
		)
	}

	jamSelesaiParse, err := utils.ParseJam(
		jamSelesai,
		loc,
	)

	if err != nil {
		return fmt.Errorf(
			"jam_selesai tidak valid: %w",
			err,
		)
	}

	if !jamSelesaiParse.After(jamMulaiParse) {
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
		nama,

		// PENTING:
		// kirim string, bukan time.Time
		jamMulai,

		// PENTING:
		// kirim string, bukan time.Time
		jamSelesai,

		tanggalID,
		utils.NullIfEmpty(linkGmaps),
		rating,
	)

	if err != nil {
		return fmt.Errorf(
			"gagal insert rundown: %w",
			err,
		)
	}

	query2 := `DELETE FROM error WHERE id = $1;`

	_, err = tx.Exec(context.Background(), query2, id)

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
