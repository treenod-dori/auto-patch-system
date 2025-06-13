package main

import (
	"auto-patch-system/patchFiles"
	"auto-patch-system/reservations"
	"auto-patch-system/utils"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type ErrorResponse struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"` // 상세 정보를 옵션으로 포함
	Code    int         `json:"code"`
}

func main() {
	sqliteDBConfig, err := utils.NewSQLiteConfig()
	if err != nil {
		log.Println(errors.New(err.Error()))
		return
	}

	err = utils.InitSQLite(sqliteDBConfig)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	r := gin.Default()
	// CORS 설정 커스터마이즈
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Custom-Message", "Authorization"},
		ExposeHeaders:    []string{"X-Custom-Message", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	patchFileController := patchFiles.InitPatchFilesController()
	reservationController := reservations.InitReservationController()

	patchFilesGroup := r.Group("/patchFiles")
	{
		// 크롤링 결과에서 Send 버튼을 눌렀을 때 호출하는 API
		patchFilesGroup.POST("/upload", patchFileController.SaveCrawlingQuery)

		// 다운로드 버튼을 눌렀을 때 호출하는 API. 요청 온 날짜에 해당하는 파일들을 가져와서 하나의 파일로 병합한다.
		patchFilesGroup.POST("/download", patchFileController.DownloadFile)

		// 패치 파일을 테스트 실행할 때 호출하는 API
		patchFilesGroup.POST("/test", patchFileController.TestAllPatchList)

		// 예약된 패치 파일 목록을 가져올 때 호출하는 API
		patchFilesGroup.GET("", patchFileController.GetAllPatchFiles)

		// 예약된 패치 파일을 삭제할 때 호출하는 API
		patchFilesGroup.POST("/delete", patchFileController.DeletePatchFile)
	}

	reservationsGroup := r.Group("/reservations")
	{
		//예약시간에 해당하는 파일이 있다면 실행하는 API. cron에서 호출하도록 구성해보기
		reservationsGroup.GET("/exec", reservationController.ExecReservedPatchFile)

		// 알림만 전송이 필요한 경우 호출하는 API
		reservationsGroup.GET("/only-notification", reservationController.ReserveNotification)

		// 예약 탭에서 현재까지 예약된 패치 목록을 보여주는 API
		reservationsGroup.GET("", reservationController.GetAllReservations)

		// Reserve 버튼을 눌렀을 때 호출하는 API
		reservationsGroup.POST("", reservationController.ReservePatchList)
	}

	r.Run("localhost:8080")
}
