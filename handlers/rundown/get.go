package handlers

import (
	"be-date/models"
	services "be-date/services/rundown"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetData(c *gin.Context) {
	dataType := c.Param("data")

	tanggalID, err := strconv.ParseInt(
		c.Param("tanggal_id"),
		10,
		64,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "tanggal_id tidak valid",
			Error:   err.Error(),
		})
		return
	}

	switch dataType {
	case "rundown":
		data, err := services.GetRundown(tanggalID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Status:  false,
				Message: "gagal mengambil data rundown",
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.SuccessResponse{
			Status:  true,
			Message: "rundown berhasil diambil",
			Data:    data,
		})

	case "error":
		data, err := services.GetError(tanggalID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Status:  false,
				Message: "gagal mengambil data error",
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.SuccessResponse{
			Status:  true,
			Message: "data error berhasil diambil",
			Data:    data,
		})

	default:
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "jenis data tidak didukung",
			Error:   "Gunakan data rundown atau error",
		})
	}
}
