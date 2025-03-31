package controller

import (
	"auto-patch-system/reservations/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ReservationController struct {
	reservationService service.ReservationService
}

func NewReservationController(reservationService service.ReservationService) ReservationController {
	return ReservationController{reservationService: reservationService}
}

func (c ReservationController) ReserveNotification(context *gin.Context) {
	queryParams := context.Request.URL.Query()
	patchDate := queryParams.Get("date")
	reservedTime := "10:00:00"

	createErr := c.reservationService.SaveOnlyNotification(patchDate, reservedTime)
	if createErr != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data to database"})
		return
	}

	context.JSON(http.StatusOK, "File saved successfully")
}

func (c ReservationController) GetAllReservations(context *gin.Context) {
	err, reservationList := c.reservationService.GetReservations()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data to database"})
		return
	}

	context.JSON(http.StatusOK, reservationList)
}

func (c ReservationController) ReservePatchList(context *gin.Context) {
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

	saveErr := c.reservationService.SaveReservation(patchDate, files, titles)
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

func (c ReservationController) ExecReservedPatchFile(context *gin.Context) {
	todayDate := time.Now().Format("2006-01-02")

	getErr := c.reservationService.GetReservationByDate(todayDate)
	if getErr != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data from database"})
		return
	}

	return
}
