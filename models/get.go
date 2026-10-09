package models

type GetTanggal struct {
	ID            int    `json:"id"`
	Tanggal       string `json:"tanggal"`
	Status        string `json:"status"`
	StatusRundown string `json:"status_rundown"`
}

type GetRundown struct {
	ID         int      `json:"id"`
	Nama       string   `json:"nama"`
	JamMulai   string   `json:"jam_mulai"`
	JamSelesai string   `json:"jam_selesai"`
	LinkGmaps  *string  `json:"link_gmaps"`
	Rating     *float64 `json:"rating"`
}

type GetRundownResponse struct {
	Tanggal         string       `json:"tanggal"`
	Status          string       `json:"status"`
	StatusRundown   string       `json:"status_rundown"`
	TotalPeringatan int          `json:"total_peringatan"`
	Rundown         []GetRundown `json:"rundown"`
}

type RundownGetJob struct {
	Type      string
	TanggalID int64
}

// type RundownGetResult struct {
// 	Type          string
// 	Tanggal       string
// 	Status        string
// 	StatusRundown string

// 	TotalPeringatan int
// 	Rundown         []GetRundown
// 	Error           error
// }

type GetError struct {
	ID          *int64   `json:"id"`
	Nama        *string  `json:"nama"`
	JamMulai    *string  `json:"jam_mulai"`
	JamSelesai  *string  `json:"jam_selesai"`
	ErrorDetail *string  `json:"error_detail"`
	LinkGmaps   *string  `json:"link_gmaps"`
	Rating      *float64 `json:"rating"`
}

type GetErrorResponse struct {
	Tanggal         string     `json:"tanggal"`
	Status          string     `json:"status"`
	StatusRundown   string     `json:"status_rundown"`
	TotalPeringatan int        `json:"total_peringatan"`
	Error           []GetError `json:"error"`
}

type RundownGetResult struct {
	Type            string
	Tanggal         string
	Status          string
	StatusRundown   string
	TotalPeringatan int
	Rundown         []GetRundown
	ErrorData       []GetError
	Error           error
}
