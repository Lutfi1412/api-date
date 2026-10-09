package models

type DeleteResult struct {
	ID      int64  `json:"id"`
	Nama    string `json:"nama"`
	Message string `json:"message"`
}

type DeleteResponse struct {
	Table    string         `json:"table"`
	Total    int            `json:"total"`
	Berhasil []DeleteResult `json:"berhasil"`
	Gagal    []DeleteResult `json:"gagal"`
}

type DeleteJob struct {
	ID    int64
	Table string
}

type DeleteJobResult struct {
	ID      int64
	Nama    string
	Success bool
	Error   error
}
