package handlers

import (
	"be-date/models"
	services "be-date/services/tanggal"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTanggal(c *gin.Context) {

	var input models.CreateTanggal

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "inputan tidak valid",
			Error:   err.Error(),
		})
		return
	}

	// Validasi format tanggal
	tanggal, err := time.Parse("2006-01-02", input.Tanggal)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "format tanggal tidak valid",
			Error:   "format tanggal harus YYYY-MM-DD",
		})
		return
	}

	// Simpan sebagai YYYY-MM-DD
	tanggalString := tanggal.Format("2006-01-02")

	// Timezone Jakarta
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  false,
			Message: "gagal menentukan timezone",
			Error:   err.Error(),
		})
		return
	}

	now := time.Now().In(loc)

	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		loc,
	)

	// Parse tanggal input sebagai tanggal Jakarta
	tanggalDate, err := time.ParseInLocation(
		"2006-01-02",
		tanggalString,
		loc,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "tanggal tidak valid",
			Error:   err.Error(),
		})
		return
	}

	var status string

	switch {
	case tanggalDate.Before(today):
		status = "selesai"

	case tanggalDate.Equal(today):
		status = "proses"

	default:
		status = "akan datang"
	}

	_, err = services.CreateTanggal(
		tanggalString,
		status,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "gagal menyimpan tanggal",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status:  true,
		Message: "tanggal berhasil disimpan",
	})
}
