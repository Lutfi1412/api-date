package handlers

import (
	"be-date/models"
	services "be-date/services/rundown"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func TryAgain(c *gin.Context) {
	var input models.CreateRundown

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "inputan tidak valid",
			Error:   err.Error(),
		})
		return
	}

	tanggalID := c.Param("tanggal_id")

	TanggalID, err := strconv.ParseInt(
		tanggalID,
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

	Id := c.Param("id")

	ID, err := strconv.ParseInt(
		Id,
		10,
		64,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "id tidak valid",
			Error:   err.Error(),
		})
		return
	}

	err = services.TryAgainRundown(input.Nama, input.JamMulai, input.JamSelesai, TanggalID, input.LinkGmaps, input.Rating, ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  false,
			Message: "gagal membuat rundown",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status:  true,
		Message: "berhasil membuat rundown",
		Data:    nil,
	})
}
