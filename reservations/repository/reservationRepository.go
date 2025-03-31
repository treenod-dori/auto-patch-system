package repository

import (
	"auto-patch-system/reservations/entity"
	"auto-patch-system/utils"
)

type ReservationRepository interface {
	SaveReservation(reservation entity.Reservation) error
	GetReservations(list *[]entity.Reservation) error
	GetReservationByDate(todayDate string) (error, []entity.Reservation)
	UpdateReservationStatus(todayDate string, status int) error
}

// 실제 구현체
type reservationRepository struct{}

func NewReservationRepository() ReservationRepository {
	return reservationRepository{}
}

func (r reservationRepository) SaveReservation(reservation entity.Reservation) error {
	db := utils.GetSqliteDB()

	createErr := db.Table("reservations").Create(&reservation).Error
	if createErr != nil {
		return createErr
	}
	return nil
}

func (r reservationRepository) GetReservations(list *[]entity.Reservation) error {
	db := utils.GetSqliteDB()

	selectErr := db.Table("reservations").Select("fileName", "reservationDate", "reservationTime", "success").Find(&list).Error
	if selectErr != nil {
		return selectErr
	}
	return nil
}

func (r reservationRepository) GetReservationByDate(todayDate string) (error, []entity.Reservation) {
	db := utils.GetSqliteDB()

	var result []entity.Reservation
	selectErr := db.Table("reservations").Where("reservationDate = ?", todayDate).Find(&result).Error
	if selectErr != nil {
		return selectErr, nil
	}
	return nil, result
}

func (r reservationRepository) UpdateReservationStatus(todayDate string, status int) error {
	db := utils.GetSqliteDB()

	updateErr := db.Table("reservations").Where("reservationDate = ?", todayDate).Update("success", status).Error
	if updateErr != nil {
		return updateErr
	}
	return nil
}
