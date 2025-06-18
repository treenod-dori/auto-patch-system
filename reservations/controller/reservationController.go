package controller

import (
	"auto-patch-system/notification"
	reservation "auto-patch-system/reservations/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ReservationController struct {
	reservationService  reservation.ReservationService
	notificationService notification.Service
}

func NewReservationController(notificationService notification.Service, reservationService reservation.ReservationService) ReservationController {
	return ReservationController{notificationService: notificationService, reservationService: reservationService}
}

func (controller *ReservationController) ReserveNotification(context *gin.Context) {
	queryParams := context.Request.URL.Query()
	patchDate := queryParams.Get("date")

	createErr := controller.reservationService.SaveOnlyNotification(patchDate)
	if createErr != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data to database"})
		return
	}

	context.JSON(http.StatusOK, "File saved successfully")
}

func (controller *ReservationController) GetAllReservations(context *gin.Context) {
	err, reservationList := controller.reservationService.GetReservations()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data to database"})
		return
	}

	context.JSON(http.StatusOK, reservationList)
}

func (controller *ReservationController) ReservePatchList(context *gin.Context) {
	// 날짜 가져오기
	queryParams := context.Request.URL.Query()
	patchDate := queryParams.Get("date")

	// 요청 파싱
	err := context.Request.ParseMultipartForm(1 << 20) // 최대 1MB
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save data to database"})
		return
	}

	files := context.Request.MultipartForm.File["blobs"] // 단일 Key로 읽기
	titles := context.Request.MultipartForm.Value["titles"]

	uploadDate := time.Now().Format("2006-01-02")
	if patchDate < uploadDate {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Cannot reserve a date in the past"})
		return
	}

	saveErr := controller.reservationService.SaveReservation(uploadDate, patchDate, files, titles)
	if saveErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save data to database"})
		return
	}

	// 성공 응답
	context.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "File saved successfully",
	})
}

func (controller *ReservationController) RunReservedPatchJob() error {
	todayDate := time.Now().Format("2006-01-02")
	return controller.reservationService.GetReservationByDate(todayDate)
}

func (controller *ReservationController) DeleteReservation(context *gin.Context) {
	fileName := context.PostForm("fileName")
	reservationDate := context.PostForm("reservationDate")

	if fileName == "" || reservationDate == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Missing form fields"})
		return
	}

	// 패치 파일 삭제
	deleteErr := controller.reservationService.DeleteReservation(fileName, reservationDate)
	if deleteErr != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": deleteErr.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Reservation deleted successfully",
	})
}
