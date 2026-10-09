package models

type CreateTanggal struct {
	Tanggal string `json:"tanggal" binding:"required"`
}

type CreateRundown struct {
	Nama       string   `json:"nama" binding:"required"`
	JamMulai   string   `json:"jam_mulai" binding:"required"`
	JamSelesai string   `json:"jam_selesai" binding:"required"`
	LinkGmaps  string   `json:"link_gmap"`
	Rating     *float64 `json:"rating"`
}

type RundownJob struct {
	ID        int
	TanggalID int64
	Data      CreateRundown
}

type RundownJobResult struct {
	JobID       int
	TanggalID   int64
	Data        CreateRundown
	Success     bool
	Error       error
	ErrorDetail string
}
