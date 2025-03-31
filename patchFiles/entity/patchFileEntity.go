package entity

type PatchFile struct {
	Title           string `json:"title" gorm:"column:title;primary_key"`
	ReservationDate string `json:"reservationDate" gorm:"column:reservationDate;primary_key"`
	PatchData       []byte `json:"patchData" gorm:"column:patchData"`
}
