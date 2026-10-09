package workers

import (
	"be-date/config"
	"be-date/models"
	"context"
	"fmt"
	"sync"
)

func GetRundownWorker(
	wg *sync.WaitGroup,
	workerID int,
	jobs <-chan models.RundownGetJob,
	results chan<- models.RundownGetResult,
) {

	defer wg.Done()

	for job := range jobs {

		var result models.RundownGetResult

		switch job.Type {

		case "tanggal":
			result = getTanggal(job)

		case "error":
			result = getTotalError(job)

		case "rundown":
			result = getRundownData(job)

		case "error_data":
			result = getErrorData(job)

		default:
			result.Error = fmt.Errorf(
				"type job tidak dikenal: %s",
				job.Type,
			)
		}

		results <- result
	}
}

func getErrorData(
	job models.RundownGetJob,
) models.RundownGetResult {
	ctx := context.Background()

	result := models.RundownGetResult{
		Type:      "error_data",
		ErrorData: []models.GetError{},
	}

	query := `
		SELECT
			id,
			nama,
			to_char(jam_mulai, 'HH24:MI:SS'),
			to_char(jam_selesai, 'HH24:MI:SS'),
			error_detail,
			link_gmaps,
			rating
		FROM public.error
		WHERE tanggal_id = $1
		ORDER BY jam_mulai ASC
	`

	rows, err := config.DB.Query(
		ctx,
		query,
		job.TanggalID,
	)
	if err != nil {
		result.Error = fmt.Errorf(
			"gagal mengambil data error: %w",
			err,
		)
		return result
	}

	defer rows.Close()

	for rows.Next() {
		var data models.GetError

		err := rows.Scan(
			&data.ID,
			&data.Nama,
			&data.JamMulai,
			&data.JamSelesai,
			&data.ErrorDetail,
			&data.LinkGmaps,
			&data.Rating,
		)
		if err != nil {
			result.Error = fmt.Errorf(
				"gagal membaca data error: %w",
				err,
			)
			return result
		}

		result.ErrorData = append(
			result.ErrorData,
			data,
		)
	}

	if err := rows.Err(); err != nil {
		result.Error = fmt.Errorf(
			"error saat membaca rows error: %w",
			err,
		)
		return result
	}

	return result
}
func getTanggal(
	job models.RundownGetJob,
) models.RundownGetResult {

	ctx := context.Background()

	result := models.RundownGetResult{
		Type: "tanggal",
	}

	query := `
		SELECT
			tanggal,
			status,
			status_rundown
		FROM tanggal
		WHERE id = $1
	`

	err := config.DB.QueryRow(
		ctx,
		query,
		job.TanggalID,
	).Scan(
		&result.Tanggal,
		&result.Status,
		&result.StatusRundown,
	)

	if err != nil {
		result.Error = fmt.Errorf(
			"gagal mengambil data tanggal: %w",
			err,
		)

		return result
	}

	return result
}

func getTotalError(
	job models.RundownGetJob,
) models.RundownGetResult {

	ctx := context.Background()

	result := models.RundownGetResult{
		Type: "error",
	}

	query := `
		SELECT COUNT(*)
		FROM error
		WHERE tanggal_id = $1
	`

	err := config.DB.QueryRow(
		ctx,
		query,
		job.TanggalID,
	).Scan(
		&result.TotalPeringatan,
	)

	if err != nil {
		result.Error = fmt.Errorf(
			"gagal menghitung total peringatan: %w",
			err,
		)

		return result
	}

	return result
}

func getRundownData(
	job models.RundownGetJob,
) models.RundownGetResult {

	ctx := context.Background()

	result := models.RundownGetResult{
		Type:    "rundown",
		Rundown: []models.GetRundown{},
	}

	query := `
		SELECT
			id,
			nama,
			jam_mulai,
			jam_selesai,
			link_gmaps,
			rating
		FROM rundown
		WHERE tanggal_id = $1
		ORDER BY jam_mulai ASC
	`

	rows, err := config.DB.Query(
		ctx,
		query,
		job.TanggalID,
	)

	if err != nil {
		result.Error = fmt.Errorf(
			"gagal mengambil data rundown: %w",
			err,
		)

		return result
	}

	defer rows.Close()

	for rows.Next() {

		var data models.GetRundown

		err := rows.Scan(
			&data.ID,
			&data.Nama,
			&data.JamMulai,
			&data.JamSelesai,
			&data.LinkGmaps,
			&data.Rating,
		)

		if err != nil {
			result.Error = fmt.Errorf(
				"gagal membaca data rundown: %w",
				err,
			)

			return result
		}

		result.Rundown = append(
			result.Rundown,
			data,
		)
	}

	if err := rows.Err(); err != nil {
		result.Error = fmt.Errorf(
			"error saat membaca rows rundown: %w",
			err,
		)

		return result
	}

	return result
}
