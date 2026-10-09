package handlers

import (
	"be-date/models"

	services "be-date/services/tanggal"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeleteTanggal(c *gin.Context) {

	idstr := c.Param("id")

	if idstr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "id tidak terbaca",
		})
		return
	}

	id, _ := strconv.Atoi(idstr)

	_, err := services.DeleteTanggal(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  false,
			Message: "gagal menghapus pesan",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status:  true,
		Message: "pesan berhasil di hapus",
	})
}
