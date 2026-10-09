package handlers

import (
	"be-date/models"
	services "be-date/services/rundown"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetRundown(c *gin.Context) {

	tanggalIDString := c.Param("tanggal_id")

	tanggalID, err := strconv.ParseInt(
		tanggalIDString,
		10,
		64,
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.SuccessResponse{
				Status:  false,
				Message: "tanggal_id tidak valid",
				Data:    nil,
			},
		)

		return
	}

	data, err := services.GetRundown(
		tanggalID,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.ErrorResponse{
				Status:  false,
				Message: "gagal mengambil data rundown",
				Error:   err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		models.SuccessResponse{
			Status:  true,
			Message: "rundown berhasil diambil",
			Data:    data,
		},
	)
}
