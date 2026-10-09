package services

import (
	"be-date/models"
	"be-date/workers"
	"sync"
)

const DeleteWorkerCount = 3

func DeleteData(
	table string,
	ids []int64,
) models.DeleteResponse {

	response := models.DeleteResponse{
		Table:    table,
		Total:    len(ids),
		Berhasil: []models.DeleteResult{},
		Gagal:    []models.DeleteResult{},
	}

	jobs := make(chan models.DeleteJob)
	results := make(chan models.DeleteJobResult)

	var wg sync.WaitGroup

	wg.Add(DeleteWorkerCount)

	// Jalankan 3 worker
	for i := 0; i < DeleteWorkerCount; i++ {
		go workers.DeleteWorker(
			&wg,
			i+1,
			jobs,
			results,
		)
	}

	// Kirim job
	go func() {

		for _, id := range ids {

			jobs <- models.DeleteJob{
				ID:    id,
				Table: table,
			}
		}

		close(jobs)

	}()

	// Tutup results setelah semua worker selesai
	go func() {

		wg.Wait()

		close(results)

	}()

	// Ambil hasil
	for result := range results {

		if result.Success {

			response.Berhasil = append(
				response.Berhasil,
				models.DeleteResult{
					ID:      result.ID,
					Nama:    result.Nama,
					Message: "berhasil dihapus",
				},
			)

			continue
		}

		response.Gagal = append(
			response.Gagal,
			models.DeleteResult{
				ID:      result.ID,
				Nama:    result.Nama,
				Message: result.Error.Error(),
			},
		)
	}

	return response
}
