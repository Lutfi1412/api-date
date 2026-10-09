package handlers

import (
	"be-date/models"
	services "be-date/services/rundown"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UpdateRundown(c *gin.Context) {

	idString := c.Param("id")

	id, err := strconv.ParseInt(
		idString,
		10,
		64,
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.SuccessResponse{
				Status:  false,
				Message: "id rundown tidak valid",
				Data:    nil,
			},
		)

		return
	}

	var input models.UpdateRundown

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.SuccessResponse{
				Status:  false,
				Message: "inputan rundown tidak valid",
				Data:    nil,
			},
		)

		return
	}

	err = services.UpdateRundown(
		id,
		input,
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.ErrorResponse{
				Status:  false,
				Message: "gagal update rundown",
				Error:   err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		models.SuccessResponse{
			Status:  true,
			Message: "rundown berhasil diupdate",
			Data:    nil,
		},
	)
}
