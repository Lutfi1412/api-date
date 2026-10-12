package services

import (
	"be-date/config"
	"be-date/models"
	"context"
	"fmt"
)

func GetDokumentasi(
	id int64,
) (models.DokumentasiResponse, error) {
	var data models.DokumentasiResponse

	query := `
		SELECT id, dokumentasi
		FROM public.tanggal
		WHERE id = $1
	`

	err := config.DB.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&data.ID,
		&data.Dokumentasi,
	)

	if err != nil {
		return data, fmt.Errorf(
			"gagal mengambil dokumentasi: %w",
			err,
		)
	}

	return data, nil
}

func UpdateDokumentasi(
	id int64,
	link string,
) (models.DokumentasiResponse, error) {
	var data models.DokumentasiResponse

	ctx := context.Background()

	tx, err := config.DB.Begin(ctx)
	if err != nil {
		return data, fmt.Errorf(
			"gagal memulai transaksi: %w",
			err,
		)
	}

	// Rollback otomatis jika transaksi belum berhasil di-commit.
	defer tx.Rollback(ctx)

	query := `
		UPDATE public.tanggal
		SET dokumentasi = NULLIF($1, '')
		WHERE id = $2
		RETURNING id, dokumentasi
	`

	err = tx.QueryRow(
		ctx,
		query,
		link,
		id,
	).Scan(
		&data.ID,
		&data.Dokumentasi,
	)

	if err != nil {
		return data, fmt.Errorf(
			"gagal memperbarui dokumentasi: %w",
			err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return data, fmt.Errorf(
			"gagal commit transaksi dokumentasi: %w",
			err,
		)
	}

	return data, nil
}
