package entity

type Reservation struct {
	ReservationDate string `json:"reservationDate" gorm:"column:reservationDate;primary_key"`
	FileName        string `json:"fileName" gorm:"column:fileName;primary_key"`
	PatchData       []byte `json:"patchData" gorm:"column:patchData"`
	ReservationTime string `json:"reservationTime" gorm:"column:reservationTime"`
	Success         int    `json:"success" gorm:"column:success"`
}
