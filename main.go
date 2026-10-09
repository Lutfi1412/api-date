package main

import (
	"be-date/config"
	Rundown "be-date/handlers/rundown"
	Tanggal "be-date/handlers/tanggal"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// urlExample := "postgres://username:password@localhost:5432/database_name"

	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env error")
	}

	conn, err := config.Connect()
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}

	defer conn.Close()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	r.POST("/create-tanggal", Tanggal.CreateTanggal)
	r.DELETE("/delete-tanggal/:id", Tanggal.DeleteTanggal)
	r.GET("/get-tanggal", Tanggal.GetDataHandler)

	//rundown

	r.POST(
		"/create-rundown/:tanggal_id",
		Rundown.CreateRundown,
	)
	r.GET(
		"/get-rundown/:tanggal_id",
		Rundown.GetRundown,
	)

	r.PUT(
		"/update-rundown/:id",
		Rundown.UpdateRundown,
	)

	r.DELETE(
		"/delete-data/:tabel",
		Rundown.DeleteData,
	)

	r.POST("/try-again/:id/:tanggal_id", Rundown.TryAgain)

	r.Run(":8080")

}
