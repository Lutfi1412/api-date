package handlers

import (
	"be-date/models"
	services "be-date/services/rundown"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateRundown(c *gin.Context) {

	// ==========================================
	// 1. Ambil tanggal_id dari URL
	// ==========================================

	tanggalIDString := c.Param("tanggal_id")

	tanggalID, err := strconv.ParseInt(
		tanggalIDString,
		10,
		64,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "tanggal_id tidak valid",
			Error:   "tanggal_id harus berupa angka",
		})
		return
	}

	// ==========================================
	// 2. Ambil array JSON
	// ==========================================

	var input []models.CreateRundown

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "inputan rundown tidak valid",
			Error:   err.Error(),
		})
		return
	}

	// ==========================================
	// 3. Pastikan array tidak kosong
	// ==========================================

	if len(input) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "data rundown tidak boleh kosong",
			Error:   "minimal terdapat 1 rundown",
		})
		return
	}

	// ==========================================
	// 4. Jalankan service
	// ==========================================

	err = services.CreateRundown(
		tanggalID,
		input,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "gagal membuat rundown",
			Error:   err.Error(),
		})
		return
	}

	// ==========================================
	// 5. Response
	// ==========================================

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status:  true,
		Message: "rundown berhasil diproses",
	})
}
