package handlers

import (
	"be-date/models"
	services "be-date/services/tanggal"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDataHandler(c *gin.Context) {

	dataTanggal, err := services.GetTanggal()

	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "gagal mengambil data Tanggal",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status:  true,
		Message: "data Tanggal berhasil diambil",
		Data:    dataTanggal,
	})
}
