package handlers

import (
	"be-date/models"
	services "be-date/services/rundown"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteData(c *gin.Context) {

	table := c.Param("tabel")

	// Validasi table
	if table != "rundown" && table != "error" {

		c.JSON(
			http.StatusBadRequest,
			models.SuccessResponse{
				Status:  false,
				Message: "tabel tidak valid",
				Data:    nil,
			},
		)

		return
	}

	var ids []int64

	if err := c.ShouldBindJSON(&ids); err != nil {

		c.JSON(
			http.StatusBadRequest,
			models.ErrorResponse{
				Status:  false,
				Message: "payload ID tidak valid",
				Error:   err.Error(),
			},
		)

		return
	}

	if len(ids) == 0 {

		c.JSON(
			http.StatusBadRequest,
			models.ErrorResponse{
				Status:  false,
				Message: "ID tidak boleh kosong",
				Error:   "ID tidak boleh kosong",
			},
		)

		return
	}

	result := services.DeleteData(
		table,
		ids,
	)

	c.JSON(
		http.StatusOK,
		models.SuccessResponse{
			Status:  true,
			Message: "proses hapus selesai",
			Data:    result,
		},
	)
}
