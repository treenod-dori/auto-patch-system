package entity

type Reservation struct {
	ReservationDate string `json:"reservationDate" gorm:"column:reservationDate;primary_key"`
	FileName        string `json:"fileName" gorm:"column:fileName;primary_key"`
	PatchData       []byte `json:"patchData" gorm:"column:patchData"`
	UploadDate      string `json:"uploadDate" gorm:"column:uploadDate"`
	Success         int    `json:"success" gorm:"column:success"`
}
