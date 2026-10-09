package services

import (
	"be-date/config"
	"be-date/models"
	"be-date/workers"
	"context"
	"fmt"
	"sync"
)

func GetRundown(tanggalID int64) (
	models.GetRundownResponse,
	error,
) {

	jobs := make(chan models.RundownGetJob)
	results := make(chan models.RundownGetResult)

	var wg sync.WaitGroup

	const workerCount = 3

	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go workers.GetRundownWorker(
			&wg,
			i+1,
			jobs,
			results,
		)
	}

	// Kirim 3 job
	go func() {

		jobs <- models.RundownGetJob{
			Type:      "tanggal",
			TanggalID: tanggalID,
		}

		jobs <- models.RundownGetJob{
			Type:      "error",
			TanggalID: tanggalID,
		}

		jobs <- models.RundownGetJob{
			Type:      "rundown",
			TanggalID: tanggalID,
		}

		close(jobs)

	}()

	// Tutup results setelah semua worker selesai
	go func() {
		wg.Wait()
		close(results)
	}()

	var response models.GetRundownResponse

	for result := range results {

		if result.Error != nil {
			return response, result.Error
		}

		switch result.Type {

		case "tanggal":
			response.Tanggal = result.Tanggal
			response.Status = result.Status
			response.StatusRundown = result.StatusRundown

		case "error":
			response.TotalPeringatan = result.TotalPeringatan

		case "rundown":
			response.Rundown = result.Rundown
		}
	}

	// Kalau status_rundown sukses,
	// total peringatan harus 0
	if response.StatusRundown == "sukses" {
		response.TotalPeringatan = 0
	}

	if response.StatusRundown == "peringatan" &&
		response.TotalPeringatan == 0 {

		err := updateStatusRundown(
			tanggalID,
			"sukses",
		)

		if err != nil {
			return response, err
		}

		response.StatusRundown = "sukses"
	}

	return response, nil
}

func updateStatusRundown(
	tanggalID int64,
	status string,
) error {

	ctx := context.Background()

	query := `
		UPDATE tanggal
		SET status_rundown = $1
		WHERE id = $2
	`

	_, err := config.DB.Exec(
		ctx,
		query,
		status,
		tanggalID,
	)

	if err != nil {
		return fmt.Errorf(
			"gagal update status_rundown: %w",
			err,
		)
	}

	return nil
}
