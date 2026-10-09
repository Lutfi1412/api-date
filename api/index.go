package handler

import (
	"be-date/config"
	Rundown "be-date/handlers/rundown"
	Tanggal "be-date/handlers/tanggal"
	"log"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	router   *gin.Engine
	initOnce sync.Once
	initErr  error
)

func setupRouter() error {
	conn, err := config.Connect()
	if err != nil {
		return err
	}

	// Jangan tutup koneksi di sini karena akan digunakan
	// oleh handler selama instance Function masih aktif.
	_ = conn

	gin.SetMode(gin.ReleaseMode)
	router = gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "API be-date Running",
		})
	})

	router.POST("/create-tanggal", Tanggal.CreateTanggal)
	router.DELETE("/delete-tanggal/:id", Tanggal.DeleteTanggal)
	router.GET("/get-tanggal", Tanggal.GetDataHandler)

	router.POST("/create-rundown/:tanggal_id", Rundown.CreateRundown)
	router.GET("/get/:data/:tanggal_id", Rundown.GetData)
	router.PUT("/update-rundown/:id", Rundown.UpdateRundown)
	router.DELETE("/delete-data/:tabel", Rundown.DeleteData)
	router.POST("/try-again/:id/:tanggal_id", Rundown.TryAgain)

	return nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	initOnce.Do(func() {
		initErr = setupRouter()
	})

	if initErr != nil {
		log.Printf("Initialization error: %v", initErr)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	router.ServeHTTP(w, r)
}
