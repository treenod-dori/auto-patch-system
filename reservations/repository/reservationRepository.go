package repository

import (
	"auto-patch-system/reservations/entity"
	"auto-patch-system/utils"
)

type ReservationRepository interface {
	SaveReservation(reservation entity.Reservation) error
	GetReservations() ([]entity.Reservation, error)
	GetReservationByDate(todayDate string) (error, []entity.Reservation)
	UpdateReservationStatus(todayDate string, status int) error
	GetReservationByDateAndFileName(date string, fileName string) (entity.Reservation, error)
	DeleteReservation(file entity.Reservation) error
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

func (r reservationRepository) GetReservations() ([]entity.Reservation, error) {
	db := utils.GetSqliteDB()
	var reservationList []entity.Reservation

	selectErr := db.Table("reservations").Select("fileName", "reservationDate", "reservationTime", "success").Find(&reservationList).Error
	return reservationList, selectErr
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

func (r reservationRepository) GetReservationByDateAndFileName(date string, fileName string) (entity.Reservation, error) {
	db := utils.GetSqliteDB()

	var result entity.Reservation
	selectErr := db.Table("reservations").Where("reservationDate = ? AND fileName = ?", date, fileName).Find(&result).Error
	return result, selectErr
}

func (r reservationRepository) UpdateReservationStatus(todayDate string, status int) error {
	db := utils.GetSqliteDB()

	updateErr := db.Table("reservations").Where("reservationDate = ?", todayDate).Update("success", status).Error
	if updateErr != nil {
		return updateErr
	}
	return nil
}

func (r reservationRepository) DeleteReservation(data entity.Reservation) error {
	err := utils.GetSqliteDB().Table("reservations").Delete(&data).Error
	return err
}
