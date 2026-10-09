package models

type UpdateRundown struct {
	Nama       string   `json:"nama"binding:"required"`
	JamMulai   string   `json:"jam_mulai" binding:"required"`
	JamSelesai string   `json:"jam_selesai" binding:"required"`
	LinkGmaps  string   `json:"link_gmap"`
	Rating     *float64 `json:"rating"`
}
