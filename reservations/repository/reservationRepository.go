package repository

import (
	"auto-patch-system/reservations/entity"
	"auto-patch-system/utils"
	"gorm.io/gorm"
)

type ReservationRepository interface {
	SaveReservation(reservation entity.Reservation) error
	GetReservations() ([]entity.Reservation, error)
	GetReservationByDate(todayDate string) ([]entity.Reservation, error)
	UpdateReservationStatus(todayDate string, status int) error
	GetReservationByDateAndFileName(date string, fileName string) (entity.Reservation, error)
	DeleteReservation(file entity.Reservation) error
}

// 실제 구현체
type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository() ReservationRepository {
	return reservationRepository{
		db: utils.GetMySqlGormDB(),
	}
}

func (r reservationRepository) SaveReservation(reservation entity.Reservation) error {
	createErr := r.db.Table("reservations").Create(&reservation).Error
	return createErr
}

func (r reservationRepository) GetReservations() ([]entity.Reservation, error) {
	var reservationList []entity.Reservation

	selectErr := r.db.Table("reservations").
		Select("fileName", "reservationDate", "uploadDate", "success").
		Order("uploadDate DESC").
		Find(&reservationList).
		Error
	return reservationList, selectErr
}

func (r reservationRepository) GetReservationByDate(todayDate string) ([]entity.Reservation, error) {
	var result []entity.Reservation
	selectErr := r.db.Table("reservations").Where("reservationDate = ?", todayDate).Find(&result).Error
	return result, selectErr
}

func (r reservationRepository) GetReservationByDateAndFileName(date string, fileName string) (entity.Reservation, error) {
	var result entity.Reservation
	selectErr := r.db.Table("reservations").Where("reservationDate = ? AND fileName = ?", date, fileName).Find(&result).Error
	return result, selectErr
}

func (r reservationRepository) UpdateReservationStatus(todayDate string, status int) error {
	updateErr := r.db.Table("reservations").Where("reservationDate = ?", todayDate).Update("success", status).Error
	return updateErr
}

func (r reservationRepository) DeleteReservation(data entity.Reservation) error {
	err := r.db.Table("reservations").Delete(&data).Error
	return err
}
