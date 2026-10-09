package services

import (
	"be-date/config"
	"be-date/models"
	"context"
	"fmt"
	"time"
)

func GetTanggal() ([]models.GetTanggal, error) {
	ctx := context.Background()

	query := `
		SELECT id, tanggal, status, status_rundown
		FROM tanggal
		ORDER BY tanggal DESC LIMIT 10
	`

	rows, err := config.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Gunakan timezone Asia/Jakarta
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, err
	}

	// Ambil tanggal hari ini berdasarkan WIB
	now := time.Now().In(loc)

	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0,
		0,
		0,
		0,
		loc,
	)

	var dataTanggal []models.GetTanggal

	for rows.Next() {

		var (
			tanggal   models.GetTanggal
			tanggalDB time.Time
		)

		if err := rows.Scan(
			&tanggal.ID,
			&tanggalDB,
			&tanggal.Status,
			&tanggal.StatusRundown,
		); err != nil {
			return nil, err
		}

		// Ambil tanggal saja dan jadikan timezone Jakarta
		tanggalDate := time.Date(
			tanggalDB.Year(),
			tanggalDB.Month(),
			tanggalDB.Day(),
			0,
			0,
			0,
			0,
			loc,
		)

		var statusBaru string

		switch {
		case tanggalDate.Before(today):
			statusBaru = "selesai"

		case tanggalDate.Equal(today):
			statusBaru = "proses"

		default:
			statusBaru = "akan datang"
		}

		if tanggal.Status != statusBaru {

			updateQuery := `
				UPDATE tanggal
				SET status = $1
				WHERE id = $2
			`

			_, err := config.DB.Exec(
				ctx,
				updateQuery,
				statusBaru,
				tanggal.ID,
			)

			if err != nil {
				return nil, err
			}
		}

		// Masukkan status terbaru ke response
		tanggal.Status = statusBaru

		// =========================
		// CEK STATUS RUNDOWN
		// =========================

		if tanggal.StatusRundown == "peringatan" {

			var totalError int

			errorQuery := `
				SELECT COUNT(*)
				FROM error
				WHERE tanggal_id = $1
			`

			err := config.DB.QueryRow(
				ctx,
				errorQuery,
				tanggal.ID,
			).Scan(&totalError)

			if err != nil {
				return nil, err
			}

			// Kalau sudah tidak ada error,
			// ubah status rundown menjadi sukses
			if totalError == 0 {

				updateRundownQuery := `
					UPDATE tanggal
					SET status_rundown = 'sukses'
					WHERE id = $1
				`

				_, err := config.DB.Exec(
					ctx,
					updateRundownQuery,
					tanggal.ID,
				)

				if err != nil {
					return nil, err
				}

				// Update response juga
				tanggal.StatusRundown = "sukses"
			}
		}

		// Format tanggal menjadi YYYY-MM-DD
		// Format tanggal menjadi 20 Oktober 2026
		namaBulan := []string{
			"Januari", "Februari", "Maret", "April",
			"Mei", "Juni", "Juli", "Agustus",
			"September", "Oktober", "November", "Desember",
		}

		tanggal.Tanggal = fmt.Sprintf(
			"%d %s %d",
			tanggalDate.Day(),
			namaBulan[int(tanggalDate.Month())-1],
			tanggalDate.Year(),
		)

		dataTanggal = append(dataTanggal, tanggal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dataTanggal, nil
}
