package handlers

import (
	"be-date/models"
	services "be-date/services/tanggal"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetDokumentasi(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "id tidak valid",
			Error:   "id harus berupa angka positif",
		})
		return
	}

	data, err := services.GetDokumentasi(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Status:  false,
				Message: "data tanggal tidak ditemukan",
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  false,
			Message: "gagal mengambil dokumentasi",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status:  true,
		Message: "dokumentasi berhasil diambil",
		Data:    data,
	})
}

func UpdateDokumentasi(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "id tidak valid",
			Error:   "id harus berupa angka positif",
		})
		return
	}

	var input models.UpdateDokumentasiRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "payload tidak valid",
			Error:   err.Error(),
		})
		return
	}

	if input.Link == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "link wajib diisi",
			Error:   "gunakan field link dengan nilai string",
		})
		return
	}

	data, err := services.UpdateDokumentasi(id, *input.Link)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Status:  false,
				Message: "data tanggal tidak ditemukan",
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  false,
			Message: "gagal memperbarui dokumentasi",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status:  true,
		Message: "dokumentasi berhasil diperbarui",
		Data:    data,
	})
}
